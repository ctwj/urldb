package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/ctwj/urldb/utils"
	"github.com/meilisearch/meilisearch-go"
)

// MeilisearchDocument 搜索文档结构
type MeilisearchDocument struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	SaveURL     string    `json:"save_url"`
	FileSize    string    `json:"file_size"`
	Key         string    `json:"key"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	PanName     string    `json:"pan_name"`
	PanID       *uint     `json:"pan_id"`
	Author      string    `json:"author"`
	Cover       string    `json:"cover"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// 高亮字段
	TitleHighlight       string   `json:"_title_highlight,omitempty"`
	DescriptionHighlight string   `json:"_description_highlight,omitempty"`
	CategoryHighlight    string   `json:"_category_highlight,omitempty"`
	TagsHighlight        []string `json:"_tags_highlight,omitempty"`
}

// MeilisearchService Meilisearch服务
type MeilisearchService struct {
	client    meilisearch.ServiceManager
	index     meilisearch.IndexManager
	indexName string
	enabled   bool
}

// NewMeilisearchService 创建Meilisearch服务
func NewMeilisearchService(host, port, masterKey, indexName string, enabled bool) *MeilisearchService {
	if !enabled {
		return &MeilisearchService{
			enabled: false,
		}
	}

	// 构建服务器URL
	serverURL := fmt.Sprintf("http://%s:%s", host, port)

	// 创建客户端
	var client meilisearch.ServiceManager

	if masterKey != "" {
		client = meilisearch.New(serverURL, meilisearch.WithAPIKey(masterKey))
	} else {
		client = meilisearch.New(serverURL)
	}

	// 获取索引
	index := client.Index(indexName)

	return &MeilisearchService{
		client:    client,
		index:     index,
		indexName: indexName,
		enabled:   enabled,
	}
}

// IsEnabled 检查是否启用
func (m *MeilisearchService) IsEnabled() bool {
	return m.enabled
}

// HealthCheck 健康检查
func (m *MeilisearchService) HealthCheck() error {
	if !m.enabled {
		utils.Debug("Meilisearch未启用，跳过健康检查")
		return fmt.Errorf("Meilisearch未启用")
	}

	// utils.Debug("开始Meilisearch健康检查")

	// 使用官方SDK的健康检查
	_, err := m.client.Health()
	if err != nil {
		// utils.Error("Meilisearch健康检查失败: %v", err)
		return fmt.Errorf("Meilisearch健康检查失败: %v", err)
	}

	utils.Debug("Meilisearch健康检查成功")
	return nil
}

// CreateIndex 创建索引
func (m *MeilisearchService) CreateIndex() error {
	if !m.enabled {
		return nil
	}

	// 创建索引配置
	indexConfig := &meilisearch.IndexConfig{
		Uid:        m.indexName,
		PrimaryKey: "id",
	}

	// 创建索引
	_, err := m.client.CreateIndex(indexConfig)
	if err != nil {
		// 如果索引已存在，返回成功
		utils.Debug("Meilisearch索引创建失败或已存在: %v", err)
		return nil
	}

	utils.Debug("Meilisearch索引创建成功: %s", m.indexName)

	// 配置索引设置
	settings := &meilisearch.Settings{
		// 配置可过滤的属性
		FilterableAttributes: []string{
			"pan_id",
			"pan_name",
			"category",
			"tags",
		},
		// 配置可搜索的属性
		SearchableAttributes: []string{
			"title",
			"description",
			"category",
			"tags",
		},
		// 配置可排序的属性
		SortableAttributes: []string{
			"created_at",
			"updated_at",
			"id",
		},
	}

	// 更新索引设置
	_, err = m.index.UpdateSettings(settings)
	if err != nil {
		utils.Error("更新Meilisearch索引设置失败: %v", err)
		return err
	}

	utils.Debug("Meilisearch索引设置更新成功")
	return nil
}

// UpdateIndexSettings 更新索引设置
func (m *MeilisearchService) UpdateIndexSettings() error {
	if !m.enabled {
		return nil
	}

	// 配置索引设置
	settings := &meilisearch.Settings{
		// 配置可过滤的属性
		FilterableAttributes: []string{
			"pan_id",
			"pan_name",
			"category",
			"tags",
		},
		// 配置可搜索的属性
		SearchableAttributes: []string{
			"title",
			"description",
			"category",
			"tags",
		},
		// 配置可排序的属性
		SortableAttributes: []string{
			"created_at",
			"updated_at",
			"id",
		},
	}

	// 更新索引设置
	_, err := m.index.UpdateSettings(settings)
	if err != nil {
		utils.Error("更新Meilisearch索引设置失败: %v", err)
		return err
	}

	utils.Debug("Meilisearch索引设置更新成功")
	return nil
}

// BatchAddDocuments 批量添加文档
func (m *MeilisearchService) BatchAddDocuments(docs []MeilisearchDocument) error {
	utils.Debug(fmt.Sprintf("开始批量添加文档到Meilisearch - 文档数量: %d", len(docs)))

	if !m.enabled {
		utils.Debug("Meilisearch未启用，跳过批量添加")
		return fmt.Errorf("Meilisearch未启用")
	}

	if len(docs) == 0 {
		utils.Debug("文档列表为空，跳过批量添加")
		return nil
	}

	// 转换为interface{}切片
	var documents []interface{}
	for i, doc := range docs {
		utils.Debug(fmt.Sprintf("转换文档 %d - ID: %d, 标题: %s, 标签数量: %d", i+1, doc.ID, doc.Title, len(doc.Tags)))
		if len(doc.Tags) > 0 {
			utils.Debug(fmt.Sprintf("文档 %d 的标签: %v", i+1, doc.Tags))
		}
		documents = append(documents, doc)
	}

	utils.Debug(fmt.Sprintf("开始调用Meilisearch API添加 %d 个文档", len(documents)))

	// 批量添加文档
	_, err := m.index.AddDocuments(documents, nil)
	if err != nil {
		utils.Error(fmt.Sprintf("Meilisearch批量添加文档失败: %v", err))
		return fmt.Errorf("Meilisearch批量添加文档失败: %v", err)
	}

	utils.Debug(fmt.Sprintf("成功批量添加 %d 个文档到Meilisearch", len(docs)))
	return nil
}

// Search 搜索文档
func (m *MeilisearchService) Search(query string, filters map[string]interface{}, page, pageSize int) ([]MeilisearchDocument, int64, error) {

	if !m.enabled {
		return nil, 0, fmt.Errorf("Meilisearch未启用")
	}

	// 构建搜索请求
	searchRequest := &meilisearch.SearchRequest{
		Query:  query,
		Offset: int64((page - 1) * pageSize),
		Limit:  int64(pageSize),
		// 启用高亮功能
		AttributesToHighlight: []string{"title", "description", "category", "tags"},
		HighlightPreTag:       "<mark>",
		HighlightPostTag:      "</mark>",
	}

	// 添加过滤器
	if len(filters) > 0 {
		var filterStrings []string
		for key, value := range filters {
			switch key {
			case "pan_id":
				// 直接使用pan_id进行过滤
				filterStrings = append(filterStrings, fmt.Sprintf("pan_id = %v", value))
			case "pan_name":
				// 使用pan_name进行过滤
				filterStrings = append(filterStrings, fmt.Sprintf("pan_name = %q", value))
			case "category":
				filterStrings = append(filterStrings, fmt.Sprintf("category = %q", value))
			case "tags":
				filterStrings = append(filterStrings, fmt.Sprintf("tags = %q", value))
			default:
				filterStrings = append(filterStrings, fmt.Sprintf("%s = %q", key, value))
			}
		}
		if len(filterStrings) > 0 {
			searchRequest.Filter = filterStrings
		}
	}

	// 执行搜索
	result, err := m.index.Search(query, searchRequest)
	if err != nil {
		return nil, 0, fmt.Errorf("搜索失败: %v", err)
	}

	// 解析结果
	var documents []MeilisearchDocument

	// 如果没有任何结果，直接返回
	if len(result.Hits) == 0 {
		utils.Debug("没有搜索结果")
		return documents, result.EstimatedTotalHits, nil
	}

	for _, hit := range result.Hits {
		// 将hit转换为MeilisearchDocument
		doc := MeilisearchDocument{}

		// 解析JSON数据 - 使用反射
		hitValue := reflect.ValueOf(hit)

		if hitValue.Kind() == reflect.Map {
			for _, key := range hitValue.MapKeys() {
				keyStr := key.String()
				value := hitValue.MapIndex(key).Interface()

				// 处理_formatted字段（包含所有高亮内容）
				if keyStr == "_formatted" {
					if rawValue, ok := value.(json.RawMessage); ok {
						// 解析_formatted字段中的高亮内容
						var formattedData map[string]interface{}
						if err := json.Unmarshal(rawValue, &formattedData); err == nil {
							// 提取高亮字段
							if titleHighlight, ok := formattedData["title"].(string); ok {
								doc.TitleHighlight = titleHighlight
							}
							if descHighlight, ok := formattedData["description"].(string); ok {
								doc.DescriptionHighlight = descHighlight
							}
							if categoryHighlight, ok := formattedData["category"].(string); ok {
								doc.CategoryHighlight = categoryHighlight
							}
							if tagsHighlight, ok := formattedData["tags"].([]interface{}); ok {
								var tags []string
								for _, tag := range tagsHighlight {
									if tagStr, ok := tag.(string); ok {
										tags = append(tags, tagStr)
									}
								}
								doc.TagsHighlight = tags
							}
						}
					}
				}

				switch keyStr {
				case "id":
					if rawID, ok := value.(json.RawMessage); ok {
						var id float64
						if err := json.Unmarshal(rawID, &id); err == nil {
							doc.ID = uint(id)
						}
					}
				case "title":
					if rawTitle, ok := value.(json.RawMessage); ok {
						var title string
						if err := json.Unmarshal(rawTitle, &title); err == nil {
							doc.Title = title
						}
					}
				case "description":
					if rawDesc, ok := value.(json.RawMessage); ok {
						var description string
						if err := json.Unmarshal(rawDesc, &description); err == nil {
							doc.Description = description
						}
					}
				case "url":
					if rawURL, ok := value.(json.RawMessage); ok {
						var url string
						if err := json.Unmarshal(rawURL, &url); err == nil {
							doc.URL = url
						}
					}
				case "save_url":
					if rawSaveURL, ok := value.(json.RawMessage); ok {
						var saveURL string
						if err := json.Unmarshal(rawSaveURL, &saveURL); err == nil {
							doc.SaveURL = saveURL
						}
					}
				case "file_size":
					if rawFileSize, ok := value.(json.RawMessage); ok {
						var fileSize string
						if err := json.Unmarshal(rawFileSize, &fileSize); err == nil {
							doc.FileSize = fileSize
						}
					}
				case "key":
					if rawKey, ok := value.(json.RawMessage); ok {
						var key string
						if err := json.Unmarshal(rawKey, &key); err == nil {
							doc.Key = key
						}
					}
				case "category":
					if rawCategory, ok := value.(json.RawMessage); ok {
						var category string
						if err := json.Unmarshal(rawCategory, &category); err == nil {
							doc.Category = category
						}
					}
				case "tags":
					if rawTags, ok := value.(json.RawMessage); ok {
						var tags []string
						if err := json.Unmarshal(rawTags, &tags); err == nil {
							doc.Tags = tags
						}
					}
				case "pan_name":
					if rawPanName, ok := value.(json.RawMessage); ok {
						var panName string
						if err := json.Unmarshal(rawPanName, &panName); err == nil {
							doc.PanName = panName
						}
					}
				case "pan_id":
					if rawPanID, ok := value.(json.RawMessage); ok {
						var panID float64
						if err := json.Unmarshal(rawPanID, &panID); err == nil {
							panIDUint := uint(panID)
							doc.PanID = &panIDUint
						}
					}
				case "author":
					if rawAuthor, ok := value.(json.RawMessage); ok {
						var author string
						if err := json.Unmarshal(rawAuthor, &author); err == nil {
							doc.Author = author
						}
					}
				case "cover":
					if rawCover, ok := value.(json.RawMessage); ok {
						var cover string
						if err := json.Unmarshal(rawCover, &cover); err == nil {
							doc.Cover = cover
						}
					}
				case "created_at":
					if rawCreatedAt, ok := value.(json.RawMessage); ok {
						var createdAt string
						if err := json.Unmarshal(rawCreatedAt, &createdAt); err == nil {
							// 尝试多种时间格式
							var t time.Time
							var parseErr error
							formats := []string{
								time.RFC3339,
								"2006-01-02T15:04:05Z",
								"2006-01-02 15:04:05",
								"2006-01-02T15:04:05.000Z",
							}
							for _, format := range formats {
								if t, parseErr = time.Parse(format, createdAt); parseErr == nil {
									doc.CreatedAt = t
									break
								}
							}
						}
					}
				case "updated_at":
					if rawUpdatedAt, ok := value.(json.RawMessage); ok {
						var updatedAt string
						if err := json.Unmarshal(rawUpdatedAt, &updatedAt); err == nil {
							// 尝试多种时间格式
							var t time.Time
							var parseErr error
							formats := []string{
								time.RFC3339,
								"2006-01-02T15:04:05Z",
								"2006-01-02 15:04:05",
								"2006-01-02T15:04:05.000Z",
							}
							for _, format := range formats {
								if t, parseErr = time.Parse(format, updatedAt); parseErr == nil {
									doc.UpdatedAt = t
									break
								}
							}
						}
					}
					// 高亮字段处理 - 已移除，现在使用_formatted字段
				}
			}
		} else {
			utils.Error("hit不是Map类型，无法解析")
		}

		documents = append(documents, doc)
	}

	return documents, result.EstimatedTotalHits, nil
}

// GetAllDocuments 获取所有文档（用于调试）
func (m *MeilisearchService) GetAllDocuments() ([]MeilisearchDocument, error) {
	if !m.enabled {
		return nil, fmt.Errorf("Meilisearch未启用")
	}

	// 构建搜索请求，获取所有文档
	searchRequest := &meilisearch.SearchRequest{
		Query:  "",
		Offset: 0,
		Limit:  1000, // 获取前1000个文档
	}

	// 执行搜索
	result, err := m.index.Search("", searchRequest)
	if err != nil {
		return nil, fmt.Errorf("获取所有文档失败: %v", err)
	}

	utils.Debug("获取所有文档，总数: %d", result.EstimatedTotalHits)
	utils.Debug("获取到的文档数量: %d", len(result.Hits))

	// 解析结果
	var documents []MeilisearchDocument
	utils.Debug("获取到 %d 个文档", len(result.Hits))

	// 只显示前3个文档的字段信息
	for i, hit := range result.Hits {
		if i >= 3 {
			break
		}
		utils.Debug("文档%d的字段:", i+1)
		hitValue := reflect.ValueOf(hit)
		if hitValue.Kind() == reflect.Map {
			for _, key := range hitValue.MapKeys() {
				keyStr := key.String()
				value := hitValue.MapIndex(key).Interface()
				if rawValue, ok := value.(json.RawMessage); ok {
					utils.Debug("  %s: %s", keyStr, string(rawValue))
				} else {
					utils.Debug("  %s: %v", keyStr, value)
				}
			}
		}
	}

	return documents, nil
}

// GetIndexStats 获取索引统计信息
func (m *MeilisearchService) GetIndexStats() (map[string]interface{}, error) {
	if !m.enabled {
		return map[string]interface{}{
			"enabled": false,
			"message": "Meilisearch未启用",
		}, nil
	}

	// 获取索引统计
	stats, err := m.index.GetStats()
	if err != nil {
		return nil, fmt.Errorf("获取索引统计失败: %v", err)
	}

	utils.Debug("Meilisearch统计 - 文档数: %d, 索引中: %v", stats.NumberOfDocuments, stats.IsIndexing)

	// 转换为map
	result := map[string]interface{}{
		"enabled":           true,
		"numberOfDocuments": stats.NumberOfDocuments,
		"isIndexing":        stats.IsIndexing,
		"fieldDistribution": stats.FieldDistribution,
	}
	return result, nil
}

// DeleteDocument 删除单个文档
func (m *MeilisearchService) DeleteDocument(documentID uint) error {
	if !m.enabled {
		return fmt.Errorf("Meilisearch未启用")
	}

	// utils.Debug("开始删除Meilisearch文档 - ID: %d", documentID)

	// 删除单个文档
	documentIDStr := fmt.Sprintf("%d", documentID)
	_, err := m.index.DeleteDocument(documentIDStr)
	if err != nil {
		return fmt.Errorf("删除Meilisearch文档失败: %v", err)
	}

	utils.Debug("成功删除Meilisearch文档 - ID: %d", documentID)
	return nil
}

// ClearIndex 清空索引
func (m *MeilisearchService) ClearIndex() error {
	if !m.enabled {
		return fmt.Errorf("Meilisearch未启用")
	}

	// 清空索引
	_, err := m.index.DeleteAllDocuments()
	if err != nil {
		return fmt.Errorf("清空索引失败: %v", err)
	}

	utils.Debug("Meilisearch索引已清空")
	return nil
}
