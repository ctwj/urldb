package dto

// CopyrightClaimCreateRequest 版权申述创建请求
type CopyrightClaimCreateRequest struct {
	ResourceKey  string `json:"resource_key" validate:"required,max=255"`
	Identity     string `json:"identity" validate:"required,max=50"`
	ProofType    string `json:"proof_type" validate:"required,max=50"`
	Reason       string `json:"reason" validate:"required,max=2000"`
	ContactInfo  string `json:"contact_info" validate:"required,max=255"`
	ClaimantName string `json:"claimant_name" validate:"required,max=100"`
	ProofFiles   string `json:"proof_files" validate:"omitempty,max=2000"`
	UserAgent    string `json:"user_agent" validate:"omitempty,max=1000"`
	IPAddress    string `json:"ip_address" validate:"omitempty,max=45"`
}

// CopyrightClaimUpdateRequest 版权申述更新请求
type CopyrightClaimUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=pending approved rejected"`
	Note   string `json:"note" validate:"omitempty,max=1000"`
}

// CopyrightClaimResponse 版权申述响应
type CopyrightClaimResponse struct {
	ID           uint   `json:"id"`
	ResourceKey  string `json:"resource_key"`
	Identity     string `json:"identity"`
	ProofType    string `json:"proof_type"`
	Reason       string `json:"reason"`
	ContactInfo  string `json:"contact_info"`
	ClaimantName string `json:"claimant_name"`
	ProofFiles   string `json:"proof_files"`
	UserAgent    string `json:"user_agent"`
	IPAddress    string `json:"ip_address"`
	Status       string `json:"status"`
	Note         string `json:"note"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CopyrightClaimListRequest 版权申述列表请求
type CopyrightClaimListRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PageSize int    `query:"page_size" validate:"omitempty,min=1,max=100"`
	Status   string `query:"status" validate:"omitempty,oneof=pending approved rejected"`
}