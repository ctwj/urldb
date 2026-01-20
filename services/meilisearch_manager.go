package services

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

// MeilisearchManager Meilisearch管理器
type MeilisearchManager struct {
	service    *MeilisearchService
	repoMgr    *repo.RepositoryManager
	configRepo repo.SystemConfigRepository
	mutex      sync.RWMutex
	status     MeilisearchStatus
	stopChan   chan struct{}
	isRunning  bool

	// 同步进度控制
	syncMutex    sync.RWMutex
	syncProgress SyncProgress
	isSyncing    bool
	syncStopChan chan struct{}
}

// SyncProgress 同步进度
type SyncProgress struct {
	IsRunning      bool      `json:"is_running"`
	TotalCount     int64     `json:"total_count"`
	ProcessedCount int64     `json:"processed_count"`
	SyncedCount    int64     `json:"synced_count"`
	FailedCount    int64     `json:"failed_count"`
	StartTime      time.Time `json:"start_time"`
	EstimatedTime  string    `json:"estimated_time"`
	CurrentBatch   int       `json:"current_batch"`
	TotalBatches   int       `json:"total_batches"`
	ErrorMessage   string    `json:"error_message"`
}

// MeilisearchStatus Meilisearch状态
type MeilisearchStatus struct {
	Enabled       bool      `json:"enabled"`
	Healthy       bool      `json:"healthy"`
	LastCheck     time.Time `json:"last_check"`
	ErrorCount    int       `json:"error_count"`
	LastError     string    `json:"last_error"`
	DocumentCount int64     `json:"document_count"`
}

// NewMeilisearchManager 创建Meilisearch管理器
func NewMeilisearchManager(repoMgr *repo.RepositoryManager) *MeilisearchManager {
	return &MeilisearchManager{
		repoMgr:      repoMgr,
		stopChan:     make(chan struct{}),
		syncStopChan: make(chan struct{}),
		status: MeilisearchStatus{
			Enabled:   false,
			Healthy:   false,
			LastCheck: time.Now(),
		},
	}
}

// Initialize 初始化Meilisearch服务
func (m *MeilisearchManager) Initialize() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 设置configRepo
	m.configRepo = m.repoMgr.SystemConfigRepository

	// 获取配置
	enabled, err := m.configRepo.GetConfigBool(entity.ConfigKeyMeilisearchEnabled)
	if err != nil {
		utils.Error("获取Meilisearch启用状态失败: %v", err)
		return err
	}

	if !enabled {
		utils.Debug("Meilisearch未启用，清理服务状态")
		m.status.Enabled = false
		m.service = nil
		// 停止监控循环
		if m.stopChan != nil {
			close(m.stopChan)
			m.stopChan = make(chan struct{})
		}
		return nil
	}

	host, err := m.configRepo.GetConfigValue(entity.ConfigKeyMeilisearchHost)
	if err != nil {
		utils.Error("获取Meilisearch主机配置失败: %v", err)
		return err
	}

	port, err := m.configRepo.GetConfigValue(entity.ConfigKeyMeilisearchPort)
	if err != nil {
		utils.Error("获取Meilisearch端口配置失败: %v", err)
		return err
	}

	masterKey, err := m.configRepo.GetConfigValue(entity.ConfigKeyMeilisearchMasterKey)
	if err != nil {
		utils.Error("获取Meilisearch主密钥配置失败: %v", err)
		return err
	}

	indexName, err := m.configRepo.GetConfigValue(entity.ConfigKeyMeilisearchIndexName)
	if err != nil {
		utils.Error("获取Meilisearch索引名配置失败: %v", err)
		return err
	}

	m.service = NewMeilisearchService(host, port, masterKey, indexName, enabled)
	m.status.Enabled = enabled

	// 如果启用，创建索引并更新设置
	if enabled {
		utils.Debug("Meilisearch已启用，创建索引并更新设置")

		// 创建索引
		if err := m.service.CreateIndex(); err != nil {
			utils.Error("创建Meilisearch索引失败: %v", err)
		}

		// 更新索引设置
		if err := m.service.UpdateIndexSettings(); err != nil {
			utils.Error("更新Meilisearch索引设置失败: %v", err)
		}

		// 立即进行一次健康检查
		go func() {
			m.checkHealth()
			// 启动监控
			go m.monitorLoop()
		}()
	} else {
		utils.Debug("Meilisearch未启用")
	}

	utils.Debug("Meilisearch服务初始化完成")
	return nil
}

// IsEnabled 检查是否启用
func (m *MeilisearchManager) IsEnabled() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.status.Enabled
}

// ReloadConfig 重新加载配置
func (m *MeilisearchManager) ReloadConfig() error {
	utils.Debug("重新加载Meilisearch配置")
	return m.Initialize()
}

// GetService 获取Meilisearch服务
func (m *MeilisearchManager) GetService() *MeilisearchService {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.service
}

// GetStatus 获取状态
func (m *MeilisearchManager) GetStatus() (MeilisearchStatus, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	utils.Debug("获取Meilisearch状态 - 启用状态: %v, 健康状态: %v, 服务实例: %v", m.status.Enabled, m.status.Healthy, m.service != nil)

	if m.service != nil && m.service.IsEnabled() {
		utils.Debug("Meilisearch服务已初始化且启用，尝试获取索引统计")

		// 获取索引统计
		stats, err := m.service.GetIndexStats()
		if err != nil {
			utils.Error("获取Meilisearch索引统计失败: %v", err)
			// 即使获取统计失败，也返回当前状态
		} else {
			utils.Debug("Meilisearch索引统计: %+v", stats)

			// 更新文档数量
			if count, ok := stats["numberOfDocuments"].(float64); ok {
				m.status.DocumentCount = int64(count)
				utils.Debug("文档数量 (float64): %d", int64(count))
			} else if count, ok := stats["numberOfDocuments"].(int64); ok {
				m.status.DocumentCount = count
				utils.Debug("文档数量 (int64): %d", count)
			} else if count, ok := stats["numberOfDocuments"].(int); ok {
				m.status.DocumentCount = int64(count)
				utils.Debug("文档数量 (int): %d", int64(count))
			} else {
				utils.Error("无法解析文档数量，类型: %T, 值: %v", stats["numberOfDocuments"], stats["numberOfDocuments"])
			}

			// 不更新启用状态，保持配置中的状态
			// 启用状态应该由配置控制，而不是由服务状态控制
		}
	} else {
		utils.Debug("Meilisearch服务未初始化或未启用 - service: %v, enabled: %v", m.service != nil, m.service != nil && m.service.IsEnabled())
	}

	return m.status, nil
}

// GetStatusWithHealthCheck 获取状态并同时进行健康检查
func (m *MeilisearchManager) GetStatusWithHealthCheck() (MeilisearchStatus, error) {
	// 先进行健康检查
	m.checkHealth()

	// 然后获取状态
	return m.GetStatus()
}

// SyncResourceToMeilisearch 同步资源到Meilisearch
func (m *MeilisearchManager) SyncResourceToMeilisearch(resource *entity.Resource) error {
	utils.Debug(fmt.Sprintf("开始同步资源到Meilisearch - 资源ID: %d, URL: %s", resource.ID, resource.URL))

	if m.service == nil || !m.service.IsEnabled() {
		utils.Debug("Meilisearch服务未初始化或未启用")
		return fmt.Errorf("Meilisearch服务未初始化或未启用")
	}

	// 先进行健康检查
	if err := m.service.HealthCheck(); err != nil {
		utils.Error(fmt.Sprintf("Meilisearch健康检查失败: %v", err))
		return fmt.Errorf("Meilisearch健康检查失败: %v", err)
	}

	// 确保索引存在
	if err := m.service.CreateIndex(); err != nil {
		utils.Error(fmt.Sprintf("创建Meilisearch索引失败: %v", err))
		return fmt.Errorf("创建Meilisearch索引失败: %v", err)
	}

	// 重新加载资源及其关联数据，确保Tags被正确加载
	resourcesWithRelations, err := m.repoMgr.ResourceRepository.FindByIDs([]uint{resource.ID})
	if err != nil {
		utils.Error(fmt.Sprintf("重新加载资源失败: %v", err))
		return fmt.Errorf("重新加载资源失败: %v", err)
	}

	if len(resourcesWithRelations) == 0 {
		utils.Error(fmt.Sprintf("资源未找到: %d", resource.ID))
		return fmt.Errorf("资源未找到: %d", resource.ID)
	}

	resourceWithRelations := resourcesWithRelations[0]
	doc := m.convertResourceToDocument(&resourceWithRelations)

	// 添加调试日志，记录标签数量
	utils.Debug(fmt.Sprintf("资源ID %d 的标签数量: %d", resource.ID, len(resourceWithRelations.Tags)))
	for i, tag := range resourceWithRelations.Tags {
		utils.Debug(fmt.Sprintf("  标签 %d: ID=%d, Name=%s", i+1, tag.ID, tag.Name))
	}

	// 验证转换后的文档
	utils.Debug(fmt.Sprintf("转换后的文档标签数量: %d", len(doc.Tags)))
	if len(doc.Tags) > 0 {
		utils.Debug(fmt.Sprintf("转换后的文档标签内容: %v", doc.Tags))
	}
	err = m.service.BatchAddDocuments([]MeilisearchDocument{doc})
	if err != nil {
		return err
	}

	// 标记为已同步
	return m.repoMgr.ResourceRepository.MarkAsSyncedToMeilisearch([]uint{resource.ID})
}

// SyncAllResources 同步所有资源
func (m *MeilisearchManager) SyncAllResources() (int, error) {
	if m.service == nil || !m.service.IsEnabled() {
		return 0, fmt.Errorf("Meilisearch未启用")
	}

	// 检查是否已经在同步中
	m.syncMutex.Lock()
	if m.isSyncing {
		m.syncMutex.Unlock()
		return 0, fmt.Errorf("同步操作正在进行中")
	}

	// 初始化同步状态
	m.isSyncing = true
	m.syncProgress = SyncProgress{
		IsRunning:      true,
		TotalCount:     0,
		ProcessedCount: 0,
		SyncedCount:    0,
		FailedCount:    0,
		StartTime:      time.Now(),
		CurrentBatch:   0,
		TotalBatches:   0,
		ErrorMessage:   "",
	}
	// 重新创建停止通道
	m.syncStopChan = make(chan struct{})
	m.syncMutex.Unlock()

	// 在goroutine中执行同步，避免阻塞
	go func() {
		defer func() {
			m.syncMutex.Lock()
			m.isSyncing = false
			m.syncProgress.IsRunning = false
			m.syncMutex.Unlock()
		}()

		m.syncAllResourcesInternal()
	}()

	return 0, nil
}

// DebugGetAllDocuments 调试：获取所有文档
func (m *MeilisearchManager) DebugGetAllDocuments() error {
	if m.service == nil || !m.service.IsEnabled() {
		return fmt.Errorf("Meilisearch未启用")
	}

	utils.Debug("开始调试：获取Meilisearch中的所有文档")
	_, err := m.service.GetAllDocuments()
	if err != nil {
		utils.Error("调试获取所有文档失败: %v", err)
		return err
	}

	utils.Debug("调试完成：已获取所有文档")
	return nil
}

// syncAllResourcesInternal 内部同步方法
func (m *MeilisearchManager) syncAllResourcesInternal() {
	// 健康检查
	if err := m.service.HealthCheck(); err != nil {
		m.updateSyncProgress("", "", fmt.Sprintf("Meilisearch不可用: %v", err))
		return
	}

	// 创建索引
	if err := m.service.CreateIndex(); err != nil {
		m.updateSyncProgress("", "", fmt.Sprintf("创建索引失败: %v", err))
		return
	}

	utils.Debug("开始同步所有资源到Meilisearch...")

	// 获取总资源数量
	totalCount, err := m.repoMgr.ResourceRepository.CountUnsyncedToMeilisearch()
	if err != nil {
		m.updateSyncProgress("", "", fmt.Sprintf("获取资源总数失败: %v", err))
		return
	}

	// 分批处理
	batchSize := 100
	totalBatches := int((totalCount + int64(batchSize) - 1) / int64(batchSize))

	// 更新总数量和总批次
	m.syncMutex.Lock()
	m.syncProgress.TotalCount = totalCount
	m.syncProgress.TotalBatches = totalBatches
	m.syncMutex.Unlock()

	offset := 0
	totalSynced := 0
	currentBatch := 0

	// 预加载所有分类和平台数据到缓存
	categoryCache := make(map[uint]string)
	panCache := make(map[uint]string)

	// 获取所有分类
	categories, err := m.repoMgr.CategoryRepository.FindAll()
	if err == nil {
		for _, category := range categories {
			categoryCache[category.ID] = category.Name
		}
	}

	// 获取所有平台
	pans, err := m.repoMgr.PanRepository.FindAll()
	if err == nil {
		for _, pan := range pans {
			panCache[pan.ID] = pan.Name
		}
	}

	for {
		// 检查是否需要停止
		select {
		case <-m.syncStopChan:
			utils.Debug("同步操作被停止")
			return
		default:
		}

		currentBatch++

		// 获取一批资源（在goroutine中执行，避免阻塞）
		resourcesChan := make(chan []entity.Resource, 1)
		errChan := make(chan error, 1)

		go func() {
			// 直接查询未同步的资源，不使用分页
			resources, _, err := m.repoMgr.ResourceRepository.FindUnsyncedToMeilisearch(1, batchSize)
			if err != nil {
				errChan <- err
				return
			}
			resourcesChan <- resources
		}()

		// 等待数据库查询结果或停止信号（添加超时）
		select {
		case resources := <-resourcesChan:
			if len(resources) == 0 {
				utils.Info("资源同步完成，总共同步 %d 个资源", totalSynced)
				return
			}

			// 检查是否需要停止
			select {
			case <-m.syncStopChan:
				utils.Debug("同步操作被停止")
				return
			default:
			}

			// 转换为Meilisearch文档（确保Tags被正确加载）
			var docs []MeilisearchDocument
			for _, resource := range resources {
				utils.Debug(fmt.Sprintf("批量同步开始处理资源 %d，标签数量: %d", resource.ID, len(resource.Tags)))
				// 使用带缓存的转换方法，但传入的资源已经预加载了Tags数据
				doc := m.convertResourceToDocumentWithCache(&resource, categoryCache, panCache)
				docs = append(docs, doc)
				utils.Debug(fmt.Sprintf("批量同步资源 %d 处理完成，最终标签数量: %d", resource.ID, len(doc.Tags)))
			}

			// 检查是否需要停止
			select {
			case <-m.syncStopChan:
				utils.Debug("同步操作被停止")
				return
			default:
			}

			// 批量添加到Meilisearch（在goroutine中执行，避免阻塞）
			meilisearchErrChan := make(chan error, 1)
			go func() {
				err := m.service.BatchAddDocuments(docs)
				meilisearchErrChan <- err
			}()

			// 等待Meilisearch操作结果或停止信号（添加超时）
			select {
			case err := <-meilisearchErrChan:
				if err != nil {
					m.updateSyncProgress("", "", fmt.Sprintf("批量添加文档失败: %v", err))
					return
				}
			case <-time.After(60 * time.Second): // 60秒超时
				m.updateSyncProgress("", "", "Meilisearch操作超时")
				utils.Error("Meilisearch操作超时")
				return
			case <-m.syncStopChan:
				utils.Debug("同步操作被停止")
				return
			}

			// 检查是否需要停止
			select {
			case <-m.syncStopChan:
				utils.Debug("同步操作被停止")
				return
			default:
			}

			// 标记为已同步（在goroutine中执行，避免阻塞）
			var resourceIDs []uint
			for _, resource := range resources {
				resourceIDs = append(resourceIDs, resource.ID)
			}

			markErrChan := make(chan error, 1)
			go func() {
				err := m.repoMgr.ResourceRepository.MarkAsSyncedToMeilisearch(resourceIDs)
				markErrChan <- err
			}()

			// 等待标记操作结果或停止信号（添加超时）
			select {
			case err := <-markErrChan:
				if err != nil {
					utils.Error("标记资源同步状态失败: %v", err)
				}
			case <-time.After(30 * time.Second): // 30秒超时
				utils.Error("标记资源同步状态超时")
			case <-m.syncStopChan:
				utils.Debug("同步操作被停止")
				return
			}

			totalSynced += len(docs)
			offset += len(resources)

			// 更新进度
			m.updateSyncProgress(fmt.Sprintf("%d", totalSynced), fmt.Sprintf("%d", currentBatch), "")

			utils.Debug("已同步 %d 个资源到Meilisearch (批次 %d/%d)", totalSynced, currentBatch, totalBatches)

			// 检查是否已经同步完所有资源
			if len(resources) == 0 {
				utils.Info("资源同步完成，总共同步 %d 个资源", totalSynced)
				return
			}

		case <-time.After(30 * time.Second): // 30秒超时
			m.updateSyncProgress("", "", "数据库查询超时")
			utils.Error("数据库查询超时")
			return

		case err := <-errChan:
			m.updateSyncProgress("", "", fmt.Sprintf("获取资源失败: %v", err))
			return
		case <-m.syncStopChan:
			utils.Info("同步操作被停止")
			return
		}

		// 避免过于频繁的请求
		time.Sleep(100 * time.Millisecond)
	}

	utils.Info("资源同步完成，总共同步 %d 个资源", totalSynced)
}

// updateSyncProgress 更新同步进度
func (m *MeilisearchManager) updateSyncProgress(syncedCount, currentBatch, errorMessage string) {
	m.syncMutex.Lock()
	defer m.syncMutex.Unlock()

	if syncedCount != "" {
		if count, err := strconv.ParseInt(syncedCount, 10, 64); err == nil {
			m.syncProgress.SyncedCount = count
		}
	}

	if currentBatch != "" {
		if batch, err := strconv.Atoi(currentBatch); err == nil {
			m.syncProgress.CurrentBatch = batch
		}
	}

	if errorMessage != "" {
		m.syncProgress.ErrorMessage = errorMessage
		m.syncProgress.IsRunning = false
	}

	// 计算预估时间
	if m.syncProgress.SyncedCount > 0 {
		elapsed := time.Since(m.syncProgress.StartTime)
		rate := float64(m.syncProgress.SyncedCount) / elapsed.Seconds()
		if rate > 0 {
			remaining := float64(m.syncProgress.TotalCount-m.syncProgress.SyncedCount) / rate
			m.syncProgress.EstimatedTime = fmt.Sprintf("%.0f秒", remaining)
		}
	}
}

// GetUnsyncedCount 获取未同步资源数量
func (m *MeilisearchManager) GetUnsyncedCount() (int64, error) {
	// 直接查询未同步的资源数量
	return m.repoMgr.ResourceRepository.CountUnsyncedToMeilisearch()
}

// GetUnsyncedResources 获取未同步的资源
func (m *MeilisearchManager) GetUnsyncedResources(page, pageSize int) ([]entity.Resource, int64, error) {
	// 查询未同步到Meilisearch的资源
	return m.repoMgr.ResourceRepository.FindUnsyncedToMeilisearch(page, pageSize)
}

// GetSyncedResources 获取已同步的资源
func (m *MeilisearchManager) GetSyncedResources(page, pageSize int) ([]entity.Resource, int64, error) {
	// 查询已同步到Meilisearch的资源
	return m.repoMgr.ResourceRepository.FindSyncedToMeilisearch(page, pageSize)
}

// GetAllResources 获取所有资源
func (m *MeilisearchManager) GetAllResources(page, pageSize int) ([]entity.Resource, int64, error) {
	// 查询所有资源
	return m.repoMgr.ResourceRepository.FindAllWithPagination(page, pageSize)
}

// GetSyncProgress 获取同步进度
func (m *MeilisearchManager) GetSyncProgress() SyncProgress {
	m.syncMutex.RLock()
	defer m.syncMutex.RUnlock()
	return m.syncProgress
}

// StopSync 停止同步
func (m *MeilisearchManager) StopSync() {
	m.syncMutex.Lock()
	defer m.syncMutex.Unlock()

	if m.isSyncing {
		// 发送停止信号
		select {
		case <-m.syncStopChan:
			// 通道已经关闭，不需要再次关闭
		default:
			close(m.syncStopChan)
		}

		m.isSyncing = false
		m.syncProgress.IsRunning = false
		m.syncProgress.ErrorMessage = "同步已停止"
		utils.Debug("同步操作已停止")
	}
}

// ClearIndex 清空索引
func (m *MeilisearchManager) ClearIndex() error {
	if m.service == nil || !m.service.IsEnabled() {
		return fmt.Errorf("Meilisearch未启用")
	}

	// 清空Meilisearch索引
	if err := m.service.ClearIndex(); err != nil {
		return err
	}

	// 标记所有资源为未同步
	return m.repoMgr.ResourceRepository.MarkAllAsUnsyncedToMeilisearch()
}

// convertResourceToDocument 转换资源为搜索文档
func (m *MeilisearchManager) convertResourceToDocument(resource *entity.Resource) MeilisearchDocument {
	// 获取关联数据
	var categoryName string
	if resource.CategoryID != nil {
		category, err := m.repoMgr.CategoryRepository.FindByID(*resource.CategoryID)
		if err == nil {
			categoryName = category.Name
		}
	}

	var panName string
	if resource.PanID != nil {
		pan, err := m.repoMgr.PanRepository.FindByID(*resource.PanID)
		if err == nil {
			panName = pan.Name
		}
	}

	// 获取标签 - 从关联的Tags字段获取
	var tagNames []string
	if len(resource.Tags) > 0 {
		utils.Debug(fmt.Sprintf("处理资源 %d 的 %d 个标签", resource.ID, len(resource.Tags)))
		for i, tag := range resource.Tags {
			if tag.Name != "" {
				utils.Debug(fmt.Sprintf("标签 %d: ID=%d, Name='%s'", i+1, tag.ID, tag.Name))
				tagNames = append(tagNames, tag.Name)
			} else {
				utils.Debug(fmt.Sprintf("标签 %d: ID=%d, Name为空，跳过", i+1, tag.ID))
			}
		}
	} else {
		utils.Debug(fmt.Sprintf("资源 %d 没有关联的标签", resource.ID))
	}

	utils.Debug(fmt.Sprintf("资源 %d 最终标签数量: %d", resource.ID, len(tagNames)))

	return MeilisearchDocument{
		ID:          resource.ID,
		Title:       resource.Title,
		Description: resource.Description,
		URL:         resource.URL,
		SaveURL:     resource.SaveURL,
		FileSize:    resource.FileSize,
		Key:         resource.Key,
		Category:    categoryName,
		Tags:        tagNames,
		PanName:     panName,
		PanID:       resource.PanID,
		Author:      resource.Author,
		Cover:       resource.Cover,
		CreatedAt:   resource.CreatedAt,
		UpdatedAt:   resource.UpdatedAt,
	}
}

// convertResourceToDocumentWithCache 转换资源为搜索文档（使用缓存）
func (m *MeilisearchManager) convertResourceToDocumentWithCache(resource *entity.Resource, categoryCache map[uint]string, panCache map[uint]string) MeilisearchDocument {
	// 从缓存获取关联数据
	var categoryName string
	if resource.CategoryID != nil {
		if name, exists := categoryCache[*resource.CategoryID]; exists {
			categoryName = name
		}
	}

	var panName string
	if resource.PanID != nil {
		if name, exists := panCache[*resource.PanID]; exists {
			panName = name
		}
	}

	// 获取标签 - 从关联的Tags字段获取
	var tagNames []string
	if len(resource.Tags) > 0 {
		utils.Debug(fmt.Sprintf("批量同步处理资源 %d 的 %d 个标签", resource.ID, len(resource.Tags)))
		for i, tag := range resource.Tags {
			if tag.Name != "" {
				utils.Debug(fmt.Sprintf("批量同步标签 %d: ID=%d, Name='%s'", i+1, tag.ID, tag.Name))
				tagNames = append(tagNames, tag.Name)
			} else {
				utils.Debug(fmt.Sprintf("批量同步标签 %d: ID=%d, Name为空，跳过", i+1, tag.ID))
			}
		}
	} else {
		utils.Debug(fmt.Sprintf("批量同步资源 %d 没有关联的标签", resource.ID))
	}

	utils.Debug(fmt.Sprintf("批量同步资源 %d 最终标签数量: %d", resource.ID, len(tagNames)))

	return MeilisearchDocument{
		ID:          resource.ID,
		Title:       resource.Title,
		Description: resource.Description,
		URL:         resource.URL,
		SaveURL:     resource.SaveURL,
		FileSize:    resource.FileSize,
		Key:         resource.Key,
		Category:    categoryName,
		Tags:        tagNames,
		PanName:     panName,
		PanID:       resource.PanID,
		Author:      resource.Author,
		Cover:       resource.Cover,
		CreatedAt:   resource.CreatedAt,
		UpdatedAt:   resource.UpdatedAt,
	}
}

// monitorLoop 监控循环
func (m *MeilisearchManager) monitorLoop() {
	if m.isRunning {
		return
	}

	m.isRunning = true
	ticker := time.NewTicker(30 * time.Second) // 每30秒检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkHealth()
		case <-m.stopChan:
			return
		}
	}
}

// checkHealth 检查健康状态
func (m *MeilisearchManager) checkHealth() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.status.LastCheck = time.Now()

	utils.Debug("开始健康检查 - 服务实例: %v, 启用状态: %v", m.service != nil, m.service != nil && m.service.IsEnabled())

	if m.service == nil || !m.service.IsEnabled() {
		utils.Debug("Meilisearch服务未初始化或未启用")
		m.status.Healthy = false
		m.status.LastError = "Meilisearch未启用"
		return
	}

	// utils.Debug("开始检查Meilisearch健康状态")

	if err := m.service.HealthCheck(); err != nil {
		m.status.Healthy = false
		m.status.ErrorCount++
		m.status.LastError = err.Error()
		utils.Error("Meilisearch健康检查失败: %v", err)
	} else {
		m.status.Healthy = true
		m.status.ErrorCount = 0
		m.status.LastError = ""
		utils.Debug("Meilisearch健康检查成功")

		// 健康检查通过后，更新文档数量
		if stats, err := m.service.GetIndexStats(); err == nil {
			if count, ok := stats["numberOfDocuments"].(float64); ok {
				m.status.DocumentCount = int64(count)
			} else if count, ok := stats["numberOfDocuments"].(int64); ok {
				m.status.DocumentCount = count
			} else if count, ok := stats["numberOfDocuments"].(int); ok {
				m.status.DocumentCount = int64(count)
			}
		}
	}
}

// Stop 停止监控
func (m *MeilisearchManager) Stop() {
	if !m.isRunning {
		return
	}

	close(m.stopChan)
	m.isRunning = false
	utils.Debug("Meilisearch监控服务已停止")
}
