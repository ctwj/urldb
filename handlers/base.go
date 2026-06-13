package handlers

import (
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/services"
)

var repoManager *repo.RepositoryManager
var meilisearchManager *services.MeilisearchManager
var linkCheckService services.LinkCheckService

// SetRepositoryManager 设置Repository管理器
func SetRepositoryManager(manager *repo.RepositoryManager) {
	repoManager = manager
}

// SetMeilisearchManager 设置Meilisearch管理器
func SetMeilisearchManager(manager *services.MeilisearchManager) {
	meilisearchManager = manager
}

// SetLinkCheckService 设置链接检测服务
func SetLinkCheckService(svc services.LinkCheckService) {
	linkCheckService = svc
}
