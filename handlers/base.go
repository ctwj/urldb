package handlers

import (
	"res_db/db/repo"
)

var repoManager *repo.RepositoryManager

// SetRepositoryManager 设置Repository管理器
func SetRepositoryManager(rm *repo.RepositoryManager) {
	repoManager = rm
}
