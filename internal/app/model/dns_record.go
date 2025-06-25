package model

import "github.com/google/uuid"

type DNS_Record struct {
	BaseModel
	AssetID uuid.UUID `gorm:"type:uuid;not null;index" json:"asset_id"` // 关联到域名资产
	Host    string    `gorm:"type:varchar(512);not null" json:"host"`
	Type    string    `gorm:"type:varchar(20);not null" json:"type"` // e.g., "A", "AAAA", "CNAME"
	Value   string    `gorm:"type:varchar(512);not null" json:"value"`
	Source  string    `gorm:"type:varchar(100)" json:"source"` // "dnsx"
}
