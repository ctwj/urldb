package services

import (
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"
)

var repoManager *repo.RepositoryManager
var meilisearchManager *MeilisearchManager

// SetRepositoryManager 设置Repository管理器
func SetRepositoryManager(manager *repo.RepositoryManager) {
	repoManager = manager
}

// SetMeilisearchManager 设置Meilisearch管理器
func SetMeilisearchManager(manager *MeilisearchManager) {
	meilisearchManager = manager
}

// UnifiedSearchResources 执行统一搜索（优先使用Meilisearch，否则使用数据库搜索）并处理违禁词
func UnifiedSearchResources(keyword string, limit int, systemConfigRepo repo.SystemConfigRepository, resourceRepo repo.ResourceRepository) ([]entity.Resource, error) {
	var resources []entity.Resource
	var total int64
	var err error

	// 如果启用了Meilisearch，优先使用Meilisearch搜索
	if meilisearchManager != nil && meilisearchManager.IsEnabled() {
		// 构建MeiliSearch过滤器
		filters := make(map[string]interface{})

		// 使用Meilisearch搜索
		service := meilisearchManager.GetService()
		if service != nil {
			docs, docTotal, err := service.Search(keyword, filters, 1, limit)
			if err == nil {
				// 将Meilisearch文档转换为Resource实体
				for _, doc := range docs {
					resource := entity.Resource{
						ID:          doc.ID,
						Title:       doc.Title,
						Description: doc.Description,
						URL:         doc.URL,
						SaveURL:     doc.SaveURL,
						FileSize:    doc.FileSize,
						Key:         doc.Key,
						PanID:       doc.PanID,
						CreatedAt:   doc.CreatedAt,
						UpdatedAt:   doc.UpdatedAt,
					}
					resources = append(resources, resource)
				}
				total = docTotal

				// 获取违禁词配置并处理违禁词
				cleanWords, err := utils.GetForbiddenWordsFromConfig(func() (string, error) {
					return systemConfigRepo.GetConfigValue(entity.ConfigKeyForbiddenWords)
				})
				if err != nil {
					utils.Error("获取违禁词配置失败: %v", err)
					cleanWords = []string{} // 如果获取失败，使用空列表
				}

				// 处理违禁词替换
				if len(cleanWords) > 0 {
					resources = utils.ProcessResourcesForbiddenWords(resources, cleanWords)
				}

				return resources, nil
			} else {
				utils.Error("MeiliSearch搜索失败，回退到数据库搜索: %v", err)
			}
		}
	}

	// 如果MeiliSearch未启用、搜索失败或没有搜索关键词，使用数据库搜索
	resources, total, err = resourceRepo.Search(keyword, nil, 1, limit)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return []entity.Resource{}, nil
	}

	// 获取违禁词配置并处理违禁词
	cleanWords, err := utils.GetForbiddenWordsFromConfig(func() (string, error) {
		return systemConfigRepo.GetConfigValue(entity.ConfigKeyForbiddenWords)
	})
	if err != nil {
		utils.Error("获取违禁词配置失败: %v", err)
		cleanWords = []string{} // 如果获取失败，使用空列表
	}

	// 处理违禁词替换
	if len(cleanWords) > 0 {
		resources = utils.ProcessResourcesForbiddenWords(resources, cleanWords)
	}

	return resources, nil
}
