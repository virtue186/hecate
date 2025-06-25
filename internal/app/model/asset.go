package model

import "github.com/google/uuid"

// AssetType 定义资产的类型
type AssetType string

const (
	AssetTypeSubdomain AssetType = "subdomain"
	AssetTypeIP        AssetType = "ip"
)

type Asset struct {
	BaseModel
	ProjectID uuid.UUID `gorm:"type:uuid;not null;index" json:"project_id"`
	Value     string    `gorm:"type:varchar(512);not null;uniqueIndex:idx_asset_value_type" json:"value"`
	Type      AssetType `gorm:"type:varchar(50);not null;uniqueIndex:idx_asset_value_type" json:"type"`
	Source    string    `gorm:"type:varchar(255)" json:"source"` // 发现来源，如 "subfinder,nmap"
}
