package converter

import (
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"time"
)

// CopyrightClaimToResponse 将版权申述实体转换为响应对象
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