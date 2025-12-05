package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	panutils "github.com/ctwj/urldb/common"
	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// GetCks 获取Cookie列表
func GetCks(c *gin.Context) {
	cks, err := repoManager.CksRepository.FindAll()
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 使用新的逻辑创建 CksResponse
	var responses []dto.CksResponse
	for _, ck := range cks {
		// 获取平台信息
		var pan *dto.PanResponse
		if ck.PanID != 0 {
			panEntity, err := repoManager.PanRepository.FindByID(ck.PanID)
			if err == nil && panEntity != nil {
				pan = &dto.PanResponse{
					ID:     panEntity.ID,
					Name:   panEntity.Name,
					Key:    panEntity.Key,
					Icon:   panEntity.Icon,
					Remark: panEntity.Remark,
				}
			}
		}

		// 统计转存资源数
		count, err := repoManager.ResourceRepository.CountResourcesByCkID(ck.ID)
		if err != nil {
			count = 0 // 统计失败时设为0
		}

		response := dto.CksResponse{
			ID:               ck.ID,
			PanID:            ck.PanID,
			Idx:              ck.Idx,
			Ck:               ck.Ck,
			IsValid:          ck.IsValid,
			Space:            ck.Space,
			LeftSpace:        ck.LeftSpace,
			UsedSpace:        ck.UsedSpace,
			Username:         ck.Username,
			VipStatus:        ck.VipStatus,
			ServiceType:      ck.ServiceType,
			Remark:           ck.Remark,
			TransferredCount: count,
			Pan:              pan,
		}
		responses = append(responses, response)
	}

	SuccessResponse(c, responses)
}

// CreateCks 创建Cookie
func CreateCks(c *gin.Context) {
	var req dto.CreateCksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	// 获取平台信息以确定服务类型
	pan, err := repoManager.PanRepository.FindByID(req.PanID)
	if err != nil {
		ErrorResponse(c, "平台不存在", http.StatusBadRequest)
		return
	}

	// 根据平台名称确定服务类型
	var serviceType panutils.ServiceType
	switch pan.Name {
	case "quark":
		serviceType = panutils.Quark
	case "alipan":
		serviceType = panutils.Alipan
	case "baidu":
		serviceType = panutils.BaiduPan
	case "uc":
		serviceType = panutils.UC
	case "xunlei":
		serviceType = panutils.Xunlei
	default:
		ErrorResponse(c, "不支持的平台类型", http.StatusBadRequest)
		return
	}

	// 创建网盘服务实例
	factory := panutils.GetInstance()
	service, err := factory.CreatePanServiceByType(serviceType, &panutils.PanConfig{})
	if err != nil {
		ErrorResponse(c, "创建网盘服务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var cks *entity.Cks
	// 迅雷网盘，使用账号密码登录
	if serviceType == panutils.Xunlei {
		// 解析账号密码信息
		credentials, err := panutils.ParseCredentialsFromCk(req.Ck)
		if err != nil {
			ErrorResponse(c, "账号密码格式错误: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 验证账号密码
		if credentials.Username == "" || credentials.Password == "" {
			ErrorResponse(c, "请提供完整的账号和密码", http.StatusBadRequest)
			return
		}

		var tokenData *panutils.XunleiTokenData
		var username string

		// 使用账号密码登录
		xunleiService := service.(*panutils.XunleiPanService)
		token, err := xunleiService.LoginWithCredentials(credentials.Username, credentials.Password)
		if err != nil {
			ErrorResponse(c, "账号密码登录失败: "+err.Error(), http.StatusBadRequest)
			return
		}
		tokenData = &token
		username = credentials.Username

		// 构建extra数据
		extra := panutils.XunleiExtraData{
			Token:   tokenData,
			Captcha: &panutils.CaptchaData{},
		}

		// 如果有账号密码信息，保存到extra中
		if credentials.Username != "" && credentials.Password != "" {
			extra.Credentials = credentials
		}

		extraStr, _ := json.Marshal(extra)

		// 声明userInfo变量
		var userInfo *panutils.UserInfo

		// 设置CKSRepository以便获取用户信息
		xunleiService.SetCKSRepository(repoManager.CksRepository, entity.Cks{})

		// 获取用户信息
		userInfo, err = xunleiService.GetUserInfo(nil)
		if err != nil {
			log.Printf("获取迅雷用户信息失败，使用默认值: %v", err)
			// 如果获取失败，使用默认值
			userInfo = &panutils.UserInfo{
				Username:    username,
				VIPStatus:   false,
				ServiceType: "xunlei",
				TotalSpace:  0,
				UsedSpace:   0,
			}
		}

		leftSpaceBytes := userInfo.TotalSpace - userInfo.UsedSpace

		// 创建Cks实体
		cks = &entity.Cks{
			PanID:       req.PanID,
			Idx:         req.Idx,
			Ck:          req.Ck, // 保持原始输入
			IsValid:     userInfo.VIPStatus, // 根据VIP状态设置有效性
			Space:       userInfo.TotalSpace,
			LeftSpace:   leftSpaceBytes,
			UsedSpace:   userInfo.UsedSpace,
			Username:    userInfo.Username,
			VipStatus:   userInfo.VIPStatus,
			ServiceType: userInfo.ServiceType,
			Extra:       string(extraStr),
			Remark:      req.Remark,
		}
	} else {
		// 获取用户信息
		utils.Info("开始获取用户信息，平台: %s, Cookie长度: %d", pan.Name, len(req.Ck))
		userInfo, err := service.GetUserInfo(&req.Ck)
		if err != nil {
			utils.Error("获取用户信息失败: %v", err)
			ErrorResponse(c, "无法获取用户信息，账号创建失败: "+err.Error(), http.StatusBadRequest)
			return
		}
		utils.Info("成功获取用户信息: %+v", userInfo)

		leftSpaceBytes := userInfo.TotalSpace - userInfo.UsedSpace

		// 创建Cks实体
		cks = &entity.Cks{
			PanID:       req.PanID,
			Idx:         req.Idx,
			Ck:          req.Ck,
			IsValid:     userInfo.VIPStatus, // 根据VIP状态设置有效性
			Space:       userInfo.TotalSpace,
			LeftSpace:   leftSpaceBytes,
			UsedSpace:   userInfo.UsedSpace,
			Username:    userInfo.Username,
			VipStatus:   userInfo.VIPStatus,
			ServiceType: userInfo.ServiceType,
			Extra:       userInfo.ExtraData,
			Remark:      req.Remark,
		}
	}

	err = repoManager.CksRepository.Create(cks)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "账号创建成功",
		"cks":     converter.ToCksResponse(cks),
	})
}

// parseCapacityToBytes 将容量字符串转换为字节数
func parseCapacityToBytes(capacity string) int64 {
	if capacity == "未知" || capacity == "" {
		return 0
	}

	// 移除空格并转换为小写
	capacity = strings.TrimSpace(strings.ToLower(capacity))

	var multiplier int64 = 1
	if strings.Contains(capacity, "gb") {
		multiplier = 1024 * 1024 * 1024
		capacity = strings.Replace(capacity, "gb", "", -1)
	} else if strings.Contains(capacity, "mb") {
		multiplier = 1024 * 1024
		capacity = strings.Replace(capacity, "mb", "", -1)
	} else if strings.Contains(capacity, "kb") {
		multiplier = 1024
		capacity = strings.Replace(capacity, "kb", "", -1)
	} else if strings.Contains(capacity, "b") {
		capacity = strings.Replace(capacity, "b", "", -1)
	}

	// 解析数字
	capacity = strings.TrimSpace(capacity)
	if capacity == "" {
		return 0
	}

	// 尝试解析浮点数
	if strings.Contains(capacity, ".") {
		if val, err := strconv.ParseFloat(capacity, 64); err == nil {
			return int64(val * float64(multiplier))
		}
	} else {
		if val, err := strconv.ParseInt(capacity, 10, 64); err == nil {
			return val * multiplier
		}
	}

	return 0
}

// GetCksByID 根据ID获取Cookie详情
func GetCksByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	cks, err := repoManager.CksRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "Cookie不存在", http.StatusNotFound)
		return
	}

	response := converter.ToCksResponse(cks)
	SuccessResponse(c, response)
}

// UpdateCks 更新Cookie
func UpdateCks(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateCksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	cks, err := repoManager.CksRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "Cookie不存在", http.StatusNotFound)
		return
	}

	if req.PanID != 0 {
		cks.PanID = req.PanID
	}
	if req.Idx != 0 {
		cks.Idx = req.Idx
	}
	if req.Ck != "" {
		cks.Ck = req.Ck
	}
	// 对于 bool 类型，我们需要检查请求中是否包含该字段
	// 由于 Go 的 JSON 解析，如果字段存在且为 false，也会被正确解析
	cks.IsValid = req.IsValid
	if req.LeftSpace != 0 {
		cks.LeftSpace = req.LeftSpace
	}
	if req.UsedSpace != 0 {
		cks.UsedSpace = req.UsedSpace
	}
	if req.Username != "" {
		cks.Username = req.Username
	}
	cks.VipStatus = req.VipStatus
	if req.ServiceType != "" {
		cks.ServiceType = req.ServiceType
	}
	if req.Remark != "" {
		cks.Remark = req.Remark
	}

	// 使用专门的方法更新，确保更新所有字段包括零值
	err = repoManager.CksRepository.UpdateWithAllFields(cks)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "Cookie更新成功"})
}

// DeleteCks 删除Cookie
func DeleteCks(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	err = repoManager.CksRepository.Delete(uint(id))
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{"message": "Cookie删除成功"})
}

// GetCksByID 根据ID获取Cookie详情（使用全局repoManager）
func GetCksByIDGlobal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	cks, err := repoManager.CksRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "Cookie不存在", http.StatusNotFound)
		return
	}

	response := converter.ToCksResponse(cks)
	SuccessResponse(c, response)
}

// RefreshCapacity 刷新账号容量信息
func RefreshCapacity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	// 获取账号信息
	cks, err := repoManager.CksRepository.FindByID(uint(id))
	if err != nil {
		ErrorResponse(c, "账号不存在", http.StatusNotFound)
		return
	}

	// 获取平台信息以确定服务类型
	pan, err := repoManager.PanRepository.FindByID(cks.PanID)
	if err != nil {
		ErrorResponse(c, "平台不存在", http.StatusBadRequest)
		return
	}

	// 根据平台名称确定服务类型
	var serviceType panutils.ServiceType
	switch pan.Name {
	case "quark":
		serviceType = panutils.Quark
	case "alipan":
		serviceType = panutils.Alipan
	case "baidu":
		serviceType = panutils.BaiduPan
	case "uc":
		serviceType = panutils.UC
	case "xunlei":
		serviceType = panutils.Xunlei
	default:
		ErrorResponse(c, "不支持的平台类型", http.StatusBadRequest)
		return
	}

	// 创建网盘服务实例
	factory := panutils.GetInstance()
	service, err := factory.CreatePanServiceByType(serviceType, &panutils.PanConfig{})
	if err != nil {
		ErrorResponse(c, "创建网盘服务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var userInfo *panutils.UserInfo
	service.SetCKSRepository(repoManager.CksRepository, *cks) // 迅雷需要初始化 token 后才能获取，

	// 根据服务类型调用不同的GetUserInfo方法
	switch s := service.(type) {
	case *panutils.XunleiPanService:
		// 迅雷网盘使用存储在extra中的token，不需要传递ck参数
		userInfo, err = s.GetUserInfo(nil)
	default:
		// 其他网盘使用ck参数
		userInfo, err = service.GetUserInfo(&cks.Ck)
	}
	if err != nil {
		ErrorResponse(c, "无法获取用户信息，刷新失败: "+err.Error(), http.StatusBadRequest)
		return
	}
	leftSpaceBytes := userInfo.TotalSpace - userInfo.UsedSpace

	// 更新账号信息
	cks.Username = userInfo.Username
	cks.VipStatus = userInfo.VIPStatus
	cks.ServiceType = userInfo.ServiceType
	cks.Space = userInfo.TotalSpace
	cks.LeftSpace = leftSpaceBytes
	cks.UsedSpace = userInfo.UsedSpace
	// cks.IsValid = userInfo.VIPStatus // 根据VIP状态更新有效性

	err = repoManager.CksRepository.UpdateWithAllFields(cks)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message": "容量信息刷新成功",
		"cks":     converter.ToCksResponse(cks),
	})
}

// DeleteRelatedResources 删除关联资源
func DeleteRelatedResources(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	// 调用资源库删除关联资源
	affectedRows, err := repoManager.ResourceRepository.DeleteRelatedResources(uint(id))
	if err != nil {
		ErrorResponse(c, "删除关联资源失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"message":       "关联资源删除成功",
		"affected_rows": affectedRows,
	})
}
