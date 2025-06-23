package dto

// CreateProjectRequest 定义了创建项目的 API 请求体
type CreateProjectRequest struct {
	Name            string   `json:"name" binding:"required"`
	Description     string   `json:"description"`
	Targets         []string `json:"targets" binding:"required,min=1,dive,required"`
	ExcludedTargets []string `json:"excluded_targets"`
	ProjectSource   string   `json:"project_source"`
	Rules           string   `json:"rules"`
}

// ProjectResponse 定义了项目的通用 API 响应体
type ProjectResponse struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	ProjectSource string           `json:"project_source"`
	Rules         string           `json:"rules"`
	Targets       []TargetResponse `json:"targets"`
	CreatedAt     string           `json:"created_at"`
}

// TargetResponse 定义了目标的通用 API 响应体
type TargetResponse struct {
	Value      string `json:"value"`
	IsExcluded bool   `json:"is_excluded"`
}

type PaginatedResponse struct {
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int64       `json:"total"`
	Data     interface{} `json:"data"`
}
