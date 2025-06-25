package model

import "github.com/google/uuid"

type Port struct {
	BaseModel
	AssetID    uuid.UUID `gorm:"type:uuid;not null;index" json:"asset_id"`
	PortNumber int       `gorm:"not null" json:"port_number"`
	Protocol   string    `gorm:"type:varchar(20);not null" json:"protocol"` // e.g., "tcp", "udp"
	Service    string    `gorm:"type:varchar(100)" json:"service"`          // e.g., "http", "ssh"
}
