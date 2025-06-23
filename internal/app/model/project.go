package model

import "github.com/google/uuid"

type Project struct {
	BaseModel
	Name          string   `gorm:"type:varchar(255);not null;unique" json:"name"`
	Description   string   `gorm:"type:text" json:"description"`
	ProjectSource string   `gorm:"type:varchar(100)" json:"project_source"` // e.g., HackerOne, SRC Name
	Rules         string   `gorm:"type:text" json:"rules"`                  // Scope and rules
	Targets       []Target `gorm:"foreignKey:ProjectID" json:"targets"`
}

// Target 代表项目的一个具体目标
type Target struct {
	BaseModel
	ProjectID  uuid.UUID `gorm:"not null" json:"project_id"`
	Value      string    `gorm:"type:varchar(512);not null" json:"value"` // e.g., example.com, 192.168.1.0/24
	IsExcluded bool      `gorm:"default:false" json:"is_excluded"`        // 标记是否为排除目标
}
