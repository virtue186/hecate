package store

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hecate/internal/app/model"
)

type AssetStore interface {
	CreateBatch(assets []*model.Asset) error
	FindByID(id uuid.UUID) (*model.Asset, error)
}

func NewAssetStore(db *gorm.DB) AssetStore {
	return &dbAssetStore{db: db}
}

type dbAssetStore struct {
	db *gorm.DB
}

// CreateBatch 批量创建资产。
// GORM 在使用 PostgreSQL 时，Create a slice 会自动使用 `INSERT ... RETURNING "id"`，
// 这会将新创建记录的 ID 填充回传入的 assets 切片中。
// 我们还使用 OnConflict 来优雅地处理已存在的资产。
func (s *dbAssetStore) CreateBatch(assets []*model.Asset) error {
	if len(assets) == 0 {
		return nil
	}
	// OnConflict(clause.DoNothing) 会在遇到唯一索引冲突时，直接忽略该条记录的插入。
	// 重要的是，GORM 的 Create 方法在 Create a slice 时，会回填ID。
	return s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&assets).Error
}

func (s *dbAssetStore) FindByID(id uuid.UUID) (*model.Asset, error) {
	var asset model.Asset
	err := s.db.First(&asset, "id = ?", id).Error
	return &asset, err
}
