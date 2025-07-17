package converter

import (
	"github.com/ctwj/panResManage/db/dto"
	"github.com/ctwj/panResManage/db/entity"
)

// ToResourceResponse 将Resource实体转换为ResourceResponse
func ToResourceResponse(resource *entity.Resource) dto.ResourceResponse {
	response := dto.ResourceResponse{
		ID:          resource.ID,
		Title:       resource.Title,
		Description: resource.Description,
		URL:         resource.URL,
		PanID:       resource.PanID,
		QuarkURL:    resource.QuarkURL,
		FileSize:    resource.FileSize,
		CategoryID:  resource.CategoryID,
		ViewCount:   resource.ViewCount,
		IsValid:     resource.IsValid,
		IsPublic:    resource.IsPublic,
		CreatedAt:   resource.CreatedAt,
		UpdatedAt:   resource.UpdatedAt,
	}

	// 设置分类名称
	if resource.Category.ID != 0 {
		response.CategoryName = resource.Category.Name
	}

	// 转换标签
	response.Tags = make([]dto.TagResponse, len(resource.Tags))
	for i, tag := range resource.Tags {
		response.Tags[i] = dto.TagResponse{
			ID:            tag.ID,
			Name:          tag.Name,
			Description:   tag.Description,
			ResourceCount: 0, // 在资源上下文中，标签的资源数量不相关
		}
	}

	return response
}

// ToResourceResponseList 将Resource实体列表转换为ResourceResponse列表
func ToResourceResponseList(resources []entity.Resource) []dto.ResourceResponse {
	responses := make([]dto.ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = ToResourceResponse(&resource)
	}
	return responses
}

// ToCategoryResponse 将Category实体转换为CategoryResponse
func ToCategoryResponse(category *entity.Category, resourceCount int64, tagNames []string) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:            category.ID,
		Name:          category.Name,
		Description:   category.Description,
		ResourceCount: resourceCount,
		TagNames:      tagNames,
	}
}

// ToCategoryResponseList 将Category实体列表转换为CategoryResponse列表
func ToCategoryResponseList(categories []entity.Category, resourceCounts map[uint]int64, tagNamesMap map[uint][]string) []dto.CategoryResponse {
	responses := make([]dto.CategoryResponse, len(categories))
	for i, category := range categories {
		resourceCount := resourceCounts[category.ID]
		tagNames := tagNamesMap[category.ID]
		responses[i] = ToCategoryResponse(&category, resourceCount, tagNames)
	}
	return responses
}

// ToTagResponse 将Tag实体转换为TagResponse
func ToTagResponse(tag *entity.Tag, resourceCount int64) dto.TagResponse {
	response := dto.TagResponse{
		ID:            tag.ID,
		Name:          tag.Name,
		Description:   tag.Description,
		CategoryID:    tag.CategoryID,
		ResourceCount: resourceCount,
	}

	// 设置分类名称
	if tag.CategoryID != nil && tag.Category.ID != 0 {
		response.CategoryName = tag.Category.Name
	} else if tag.CategoryID != nil {
		// 如果CategoryID存在但Category没有预加载，设置为"未知分类"
		response.CategoryName = "未知分类"
	}

	return response
}

// ToTagResponseList 将Tag实体列表转换为TagResponse列表
func ToTagResponseList(tags []entity.Tag, resourceCounts map[uint]int64) []dto.TagResponse {
	responses := make([]dto.TagResponse, len(tags))
	for i, tag := range tags {
		resourceCount := resourceCounts[tag.ID]
		responses[i] = ToTagResponse(&tag, resourceCount)
	}
	return responses
}

// ToPanResponse 将Pan实体转换为PanResponse
func ToPanResponse(pan *entity.Pan) dto.PanResponse {
	return dto.PanResponse{
		ID:     pan.ID,
		Name:   pan.Name,
		Key:    pan.Key,
		Icon:   pan.Icon,
		Remark: pan.Remark,
	}
}

// ToPanResponseList 将Pan实体列表转换为PanResponse列表
func ToPanResponseList(pans []entity.Pan) []dto.PanResponse {
	responses := make([]dto.PanResponse, len(pans))
	for i, pan := range pans {
		responses[i] = ToPanResponse(&pan)
	}
	return responses
}

// ToCksResponse 将Cks实体转换为CksResponse
func ToCksResponse(cks *entity.Cks) dto.CksResponse {
	response := dto.CksResponse{
		ID:          cks.ID,
		PanID:       cks.PanID,
		Idx:         cks.Idx,
		Ck:          cks.Ck,
		IsValid:     cks.IsValid,
		Space:       cks.Space,
		LeftSpace:   cks.LeftSpace,
		UsedSpace:   cks.UsedSpace,
		Username:    cks.Username,
		VipStatus:   cks.VipStatus,
		ServiceType: cks.ServiceType,
		Remark:      cks.Remark,
	}

	// 设置平台信息
	if cks.Pan.ID != 0 {
		response.Pan = &dto.PanResponse{
			ID:     cks.Pan.ID,
			Name:   cks.Pan.Name,
			Key:    cks.Pan.Key,
			Icon:   cks.Pan.Icon,
			Remark: cks.Pan.Remark,
		}
	}

	return response
}

// ToCksResponseList 将Cks实体列表转换为CksResponse列表
func ToCksResponseList(cksList []entity.Cks) []dto.CksResponse {
	responses := make([]dto.CksResponse, len(cksList))
	for i, cks := range cksList {
		responses[i] = ToCksResponse(&cks)
	}
	return responses
}

// ToReadyResourceResponse 将ReadyResource实体转换为ReadyResourceResponse
func ToReadyResourceResponse(resource *entity.ReadyResource) dto.ReadyResourceResponse {
	return dto.ReadyResourceResponse{
		ID:         resource.ID,
		Title:      resource.Title,
		URL:        resource.URL,
		Category:   resource.Category,
		Tags:       resource.Tags,
		Img:        resource.Img,
		Source:     resource.Source,
		Extra:      resource.Extra,
		CreateTime: resource.CreateTime,
		IP:         resource.IP,
	}
}

// ToReadyResourceResponseList 将ReadyResource实体列表转换为ReadyResourceResponse列表
func ToReadyResourceResponseList(resources []entity.ReadyResource) []dto.ReadyResourceResponse {
	responses := make([]dto.ReadyResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = ToReadyResourceResponse(&resource)
	}
	return responses
}

// RequestToReadyResource 将ReadyResourceRequest转换为ReadyResource实体
func RequestToReadyResource(req *dto.ReadyResourceRequest) *entity.ReadyResource {
	if req == nil {
		return nil
	}

	return &entity.ReadyResource{
		Title:       &req.Title,
		Description: req.Description,
		URL:         req.Url,
		Category:    req.Category,
		Tags:        req.Tags,
		Img:         req.Img,
		Source:      req.Source,
		Extra:       req.Extra,
	}
}
