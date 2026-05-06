package handlers

import (
	"net/http"
	"strconv"

	pan "github.com/ctwj/urldb/common"
	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"

	"github.com/gin-gonic/gin"
)

var ruleManager *pan.RuleManager

// SetRuleManager 设置规则管理器
func SetRuleManager(rm *pan.RuleManager) {
	ruleManager = rm
}

// GetPanRules 获取所有网盘规则
func GetPanRules(c *gin.Context) {
	rules, err := repoManager.PanRuleRepository.FindAll()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := converter.ToPanRuleResponseList(rules)
	ListResponse(c, responses, int64(len(responses)))
}

// GetPanRule 根据ID获取网盘规则
func GetPanRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	rule, err := repoManager.PanRuleRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "规则不存在", http.StatusNotFound)
		return
	}

	response := converter.ToPanRuleResponse(rule)
	SuccessResponse(c, response)
}

// CreatePanRule 创建网盘规则
func CreatePanRule(c *gin.Context) {
	var req dto.CreatePanRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	rule := &entity.PanRule{
		PanID:       req.PanID,
		Name:        req.Name,
		Domains:     req.Domains,
		URLPatterns: req.URLPatterns,
		Priority:    req.Priority,
		Enabled:     req.Enabled,
		Remark:      req.Remark,
	}

	err := repoManager.PanRuleRepository.Create(rule)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if ruleManager != nil {
		ruleManager.Refresh()
	}

	SuccessResponse(c, gin.H{
		"id":      rule.ID,
		"message": "规则创建成功",
	})
}

// UpdatePanRule 更新网盘规则
func UpdatePanRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdatePanRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	rule, err := repoManager.PanRuleRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "规则不存在", http.StatusNotFound)
		return
	}

	if req.Name != "" {
		rule.Name = req.Name
	}
	if req.Domains != "" {
		rule.Domains = req.Domains
	}
	if req.URLPatterns != "" {
		rule.URLPatterns = req.URLPatterns
	}
	if req.Priority != 0 {
		rule.Priority = req.Priority
	}
	rule.Enabled = req.Enabled
	if req.Remark != "" {
		rule.Remark = req.Remark
	}

	err = repoManager.PanRuleRepository.Update(rule)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if ruleManager != nil {
		ruleManager.Refresh()
	}

	SuccessResponse(c, gin.H{"message": "规则更新成功"})
}

// DeletePanRule 删除网盘规则
func DeletePanRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.PanRuleRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if ruleManager != nil {
		ruleManager.Refresh()
	}

	SuccessResponse(c, gin.H{"message": "规则删除成功"})
}

// RefreshPanRules 手动刷新规则缓存
func RefreshPanRules(c *gin.Context) {
	if ruleManager == nil {
		ErrorResponse(c, "规则管理器未初始化", http.StatusInternalServerError)
		return
	}

	err := ruleManager.Refresh()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "规则缓存刷新成功"})
}