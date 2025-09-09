package converter

import (
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
)

// HotDramaToResponse 将热播剧实体转换为响应DTO
func HotDramaToResponse(drama *entity.HotDrama) *dto.HotDramaResponse {
	if drama == nil {
		return nil
	}

	return &dto.HotDramaResponse{
		ID:           drama.ID,
		CreatedAt:    drama.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    drama.UpdatedAt.Format("2006-01-02 15:04:05"),
		Title:        drama.Title,
		CardSubtitle: drama.CardSubtitle,
		EpisodesInfo: drama.EpisodesInfo,
		IsNew:        drama.IsNew,
		Rating:       drama.Rating,
		RatingCount:  drama.RatingCount,
		Year:         drama.Year,
		Region:       drama.Region,
		Genres:       drama.Genres,
		Directors:    drama.Directors,
		Actors:       drama.Actors,
		PosterURL:    drama.PosterURL,
		Category:     drama.Category,
		SubType:      drama.SubType,
		Rank:         drama.Rank,
		Source:       drama.Source,
		DoubanID:     drama.DoubanID,
		DoubanURI:    drama.DoubanURI,
	}
}

// RequestToHotDrama 将请求DTO转换为热播剧实体
func RequestToHotDrama(req *dto.HotDramaRequest) *entity.HotDrama {
	if req == nil {
		return nil
	}

	return &entity.HotDrama{
		Title:     req.Title,
		Rating:    req.Rating,
		Year:      req.Year,
		Directors: req.Directors,
		Actors:    req.Actors,
		Category:  req.Category,
		SubType:   req.SubType,
		Rank:      req.Rank,
		Source:    req.Source,
		DoubanID:  req.DoubanID,
	}
}

// HotDramaListToResponse 将热播剧实体列表转换为响应DTO
func HotDramaListToResponse(dramas []entity.HotDrama) *dto.HotDramaListResponse {
	items := make([]dto.HotDramaResponse, len(dramas))
	for i, drama := range dramas {
		response := HotDramaToResponse(&drama)
		if response != nil {
			items[i] = *response
		}
	}

	return &dto.HotDramaListResponse{
		Total: len(items),
		Items: items,
	}
}

// DoubanItemToHotDrama 将豆瓣项目转换为热播剧实体
func DoubanItemToHotDrama(item interface{}, category, subType string) *entity.HotDrama {
	// 这里需要根据实际的豆瓣数据结构进行转换
	// 暂时返回一个示例结构
	return &entity.HotDrama{
		Title:     "示例剧名",
		Rating:    0.0,
		Year:      "",
		Directors: "",
		Actors:    "",
		Category:  category,
		SubType:   subType,
		Source:    "douban",
		DoubanID:  "",
	}
}
