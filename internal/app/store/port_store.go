package store

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hecate/internal/app/model"
)

type PortStore interface {
	CreateBatch(ports []*model.Port) error
}

func NewPortStore(db *gorm.DB) PortStore {
	return &dbPortStore{db: db}
}

type dbPortStore struct {
	db *gorm.DB
}

// CreateBatch 批量创建端口信息，如果已存在则忽略
func (s *dbPortStore) CreateBatch(ports []*model.Port) error {
	if len(ports) == 0 {
		return nil
	}
	// 遇到唯一索引冲突时，忽略插入
	return s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&ports).Error
}
