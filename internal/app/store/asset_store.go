package store

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hecate/internal/app/model"
)

type AssetStore interface {
	CreateBatch(assets []*model.Asset) error
}

func NewAssetStore(db *gorm.DB) AssetStore {
	return &dbAssetStore{db: db}
}

type dbAssetStore struct {
	db *gorm.DB
}

// CreateBatch 批量创建资产，如果已存在则忽略
func (s *dbAssetStore) CreateBatch(assets []*model.Asset) error {
	if len(assets) == 0 {
		return nil
	}
	// 使用 OnConflict(clause.DoNothing) 来避免因为重复发现而导致的插入错误
	// 它会优雅地忽略掉那些已经存在的资产
	return s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&assets).Error
}
