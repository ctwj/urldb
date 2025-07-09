package handlers

import (
	"net/http"

	"res_db/models"

	"github.com/gin-gonic/gin"
)

// GetCategories 获取分类列表
func GetCategories(c *gin.Context) {
	rows, err := models.DB.Query(`
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, categories)
}

// CreateCategory 创建分类
func CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err := models.DB.QueryRow(query, req.Name, req.Description).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "分类创建成功"})
}

// UpdateCategory 更新分类
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		UPDATE categories 
		SET name = COALESCE($1, name),
		    description = COALESCE($2, description),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := models.DB.Exec(query, req.Name, req.Description, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分类更新成功"})
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有资源使用此分类
	var count int
	err := models.DB.QueryRow("SELECT COUNT(*) FROM resources WHERE category_id = $1", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下还有资源，无法删除"})
		return
	}

	result, err := models.DB.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分类删除成功"})
}
