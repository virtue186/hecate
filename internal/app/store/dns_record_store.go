package store

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hecate/internal/app/model"
)

type DnsRecordStore interface {
	CreateBatch(records []*model.DNS_Record) error
}

func NewDnsRecordStore(db *gorm.DB) DnsRecordStore {
	return &dbDnsRecordStore{db: db}
}

type dbDnsRecordStore struct {
	db *gorm.DB
}

// CreateBatch 批量创建DNS记录
func (s *dbDnsRecordStore) CreateBatch(records []*model.DNS_Record) error {
	if len(records) == 0 {
		return nil
	}
	return s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&records).Error
}
