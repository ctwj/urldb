package handlers

import (
	"net/http"
	"strconv"

	"res_db/models"

	"github.com/gin-gonic/gin"
)

// GetResources 获取资源列表
func GetResources(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	categoryID := c.Query("category_id")

	offset := (page - 1) * limit

	query := `
		SELECT r.id, r.title, r.description, r.url, r.file_path, r.file_size, 
		       r.file_type, r.category_id, c.name as category_name, r.tags, 
		       r.download_count, r.view_count, r.is_public, r.created_at, r.updated_at
		FROM resources r
		LEFT JOIN categories c ON r.category_id = c.id
		WHERE r.is_public = true
	`
	args := []interface{}{}

	if categoryID != "" {
		query += " AND r.category_id = $1"
		args = append(args, categoryID)
	}

	query += " ORDER BY r.created_at DESC LIMIT $2 OFFSET $3"
	args = append(args, limit, offset)

	rows, err := models.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var resources []models.Resource
	for rows.Next() {
		var r models.Resource
		err := rows.Scan(
			&r.ID, &r.Title, &r.Description, &r.URL, &r.FilePath, &r.FileSize,
			&r.FileType, &r.CategoryID, &r.CategoryName, &r.Tags,
			&r.DownloadCount, &r.ViewCount, &r.IsPublic, &r.CreatedAt, &r.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resources = append(resources, r)
	}

	c.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"page":      page,
		"limit":     limit,
	})
}

// GetResourceByID 根据ID获取资源
func GetResourceByID(c *gin.Context) {
	id := c.Param("id")

	var r models.Resource
	query := `
		SELECT r.id, r.title, r.description, r.url, r.file_path, r.file_size, 
		       r.file_type, r.category_id, c.name as category_name, r.tags, 
		       r.download_count, r.view_count, r.is_public, r.created_at, r.updated_at
		FROM resources r
		LEFT JOIN categories c ON r.category_id = c.id
		WHERE r.id = $1 AND r.is_public = true
	`

	err := models.DB.QueryRow(query, id).Scan(
		&r.ID, &r.Title, &r.Description, &r.URL, &r.FilePath, &r.FileSize,
		&r.FileType, &r.CategoryID, &r.CategoryName, &r.Tags,
		&r.DownloadCount, &r.ViewCount, &r.IsPublic, &r.CreatedAt, &r.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 增加浏览次数
	models.DB.Exec("UPDATE resources SET view_count = view_count + 1 WHERE id = $1", id)

	c.JSON(http.StatusOK, r)
}

// CreateResource 创建资源
func CreateResource(c *gin.Context) {
	var req models.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO resources (title, description, url, file_path, file_size, file_type, category_id, tags, is_public)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var id int
	err := models.DB.QueryRow(
		query,
		req.Title, req.Description, req.URL, req.FilePath, req.FileSize,
		req.FileType, req.CategoryID, req.Tags, req.IsPublic,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "资源创建成功"})
}

// UpdateResource 更新资源
func UpdateResource(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateResourceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		UPDATE resources 
		SET title = COALESCE($1, title),
		    description = COALESCE($2, description),
		    url = COALESCE($3, url),
		    file_path = COALESCE($4, file_path),
		    file_size = COALESCE($5, file_size),
		    file_type = COALESCE($6, file_type),
		    category_id = $7,
		    tags = COALESCE($8, tags),
		    is_public = COALESCE($9, is_public),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
	`

	result, err := models.DB.Exec(
		query,
		req.Title, req.Description, req.URL, req.FilePath, req.FileSize,
		req.FileType, req.CategoryID, req.Tags, req.IsPublic, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "资源更新成功"})
}

// DeleteResource 删除资源
func DeleteResource(c *gin.Context) {
	id := c.Param("id")

	result, err := models.DB.Exec("DELETE FROM resources WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "资源删除成功"})
}

// SearchResources 搜索资源
func SearchResources(c *gin.Context) {
	query := c.Query("q")
	categoryID := c.Query("category_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	sqlQuery := `
		SELECT r.id, r.title, r.description, r.url, r.file_path, r.file_size, 
		       r.file_type, r.category_id, c.name as category_name, r.tags, 
		       r.download_count, r.view_count, r.is_public, r.created_at, r.updated_at
		FROM resources r
		LEFT JOIN categories c ON r.category_id = c.id
		WHERE r.is_public = true
	`

	args := []interface{}{}
	argCount := 0

	if query != "" {
		argCount++
		sqlQuery += " AND (r.title ILIKE $" + strconv.Itoa(argCount) +
			" OR r.description ILIKE $" + strconv.Itoa(argCount) +
			" OR EXISTS (SELECT 1 FROM unnest(r.tags) tag WHERE tag ILIKE $" + strconv.Itoa(argCount) + "))"
		args = append(args, "%"+query+"%")
	}

	if categoryID != "" {
		argCount++
		sqlQuery += " AND r.category_id = $" + strconv.Itoa(argCount)
		args = append(args, categoryID)
	}

	sqlQuery += " ORDER BY r.created_at DESC LIMIT $" + strconv.Itoa(argCount+1) +
		" OFFSET $" + strconv.Itoa(argCount+2)
	args = append(args, limit, offset)

	rows, err := models.DB.Query(sqlQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var resources []models.Resource
	for rows.Next() {
		var r models.Resource
		err := rows.Scan(
			&r.ID, &r.Title, &r.Description, &r.URL, &r.FilePath, &r.FileSize,
			&r.FileType, &r.CategoryID, &r.CategoryName, &r.Tags,
			&r.DownloadCount, &r.ViewCount, &r.IsPublic, &r.CreatedAt, &r.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resources = append(resources, r)
	}

	c.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"page":      page,
		"limit":     limit,
		"query":     query,
	})
}

// GetStats 获取统计信息
func GetStats(c *gin.Context) {
	var stats models.Stats

	// 总资源数
	err := models.DB.QueryRow("SELECT COUNT(*) FROM resources WHERE is_public = true").Scan(&stats.TotalResources)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 总分类数
	err = models.DB.QueryRow("SELECT COUNT(*) FROM categories").Scan(&stats.TotalCategories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 总下载数
	err = models.DB.QueryRow("SELECT COALESCE(SUM(download_count), 0) FROM resources WHERE is_public = true").Scan(&stats.TotalDownloads)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 总浏览数
	err = models.DB.QueryRow("SELECT COALESCE(SUM(view_count), 0) FROM resources WHERE is_public = true").Scan(&stats.TotalViews)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
