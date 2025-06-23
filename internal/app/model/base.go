package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 使用软删除
}

// BeforeCreate 会在创建模型前，自动生成 UUID
func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	// 如果 ID 为空，则生成新的 UUID
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return
}
