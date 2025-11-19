package converter

import (
	"time"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
)

// CopyrightClaimToResponseWithResources 将版权申述实体和关联资源转换为响应对象
func CopyrightClaimToResponseWithResources(claim *entity.CopyrightClaim, resources []*entity.Resource) *dto.CopyrightClaimResponse {
	if claim == nil {
		return nil
	}

	// 转换关联的资源信息
	var resourceInfos []dto.ResourceInfo
	for _, resource := range resources {
		categoryName := ""
		if resource.Category.ID != 0 {
			categoryName = resource.Category.Name
		}

		panName := ""
		if resource.Pan.ID != 0 {
			panName = resource.Pan.Name
		}

		resourceInfo := dto.ResourceInfo{
			ID:          resource.ID,
			Title:       resource.Title,
			Description: resource.Description,
			URL:         resource.URL,
			SaveURL:     resource.SaveURL,
			FileSize:    resource.FileSize,
			Category:    categoryName,
			PanName:     panName,
			ViewCount:   resource.ViewCount,
			IsValid:     resource.IsValid,
			CreatedAt:   resource.CreatedAt.Format(time.RFC3339),
		}
		resourceInfos = append(resourceInfos, resourceInfo)
	}

	return &dto.CopyrightClaimResponse{
		ID:           claim.ID,
		ResourceKey:  claim.ResourceKey,
		Identity:     claim.Identity,
		ProofType:    claim.ProofType,
		Reason:       claim.Reason,
		ContactInfo:  claim.ContactInfo,
		ClaimantName: claim.ClaimantName,
		ProofFiles:   claim.ProofFiles,
		UserAgent:    claim.UserAgent,
		IPAddress:    claim.IPAddress,
		Status:       claim.Status,
		Note:         claim.Note,
		CreatedAt:    claim.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    claim.UpdatedAt.Format(time.RFC3339),
		Resources:    resourceInfos,
	}
}

// CopyrightClaimToResponse 将版权申述实体转换为响应对象（不包含资源详情）
func CopyrightClaimToResponse(claim *entity.CopyrightClaim) *dto.CopyrightClaimResponse {
	if claim == nil {
		return nil
	}

	return &dto.CopyrightClaimResponse{
		ID:           claim.ID,
		ResourceKey:  claim.ResourceKey,
		Identity:     claim.Identity,
		ProofType:    claim.ProofType,
		Reason:       claim.Reason,
		ContactInfo:  claim.ContactInfo,
		ClaimantName: claim.ClaimantName,
		ProofFiles:   claim.ProofFiles,
		UserAgent:    claim.UserAgent,
		IPAddress:    claim.IPAddress,
		Status:       claim.Status,
		Note:         claim.Note,
		CreatedAt:    claim.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    claim.UpdatedAt.Format(time.RFC3339),
		Resources:   []dto.ResourceInfo{}, // 空的资源列表
	}
}

// CopyrightClaimsToResponse 将版权申述实体列表转换为响应对象列表
func CopyrightClaimsToResponse(claims []*entity.CopyrightClaim) []*dto.CopyrightClaimResponse {
	var responses []*dto.CopyrightClaimResponse
	for _, claim := range claims {
		responses = append(responses, CopyrightClaimToResponse(claim))
	}
	return responses
}