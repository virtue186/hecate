package store

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hecate/internal/app/model"
)

// ProjectStore defines the interface for project data operations.
type ProjectStore interface {
	Create(project *model.Project) error
	FindByID(id uuid.UUID) (*model.Project, error)              // <--- 确保这一行存在！
	FindAll(page, pageSize int) ([]model.Project, int64, error) // <--- 确保这一行存在！
}

// NewProjectStore creates a new ProjectStore.
func NewProjectStore(db *gorm.DB) ProjectStore {
	return &dbProjectStore{db: db}
}

type dbProjectStore struct {
	db *gorm.DB
}

// Create creates a new project and its associated targets in a transaction.
func (s *dbProjectStore) Create(project *model.Project) error {
	// 使用事务确保项目和其目标被原子性地创建
	// 如果任何一步失败，所有操作都将回滚
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建项目主体
		if err := tx.Create(project).Error; err != nil {
			return err
		}
		// Targets 已经在 project 结构体中，GORM 会自动处理关联创建
		// 如果 Targets 数组不为空，GORM 会批量插入它们
		return nil
	})
}

func (s *dbProjectStore) FindByID(id uuid.UUID) (*model.Project, error) {
	var project model.Project
	// Preload("Targets") 会在查询 Project 的同时，自动带上其关联的所有 Target
	err := s.db.Preload("Targets").First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err // gorm.ErrRecordNotFound 会被直接返回
	}
	return &project, nil
}

func (s *dbProjectStore) FindAll(page, pageSize int) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 先计算总数，用于分页
	// 我们只关心 project 表本身的总数，所以用 Model(&model.Project{})
	s.db.Model(&model.Project{}).Count(&total)

	// 查询分页数据
	err := s.db.Preload("Targets").Order("created_at desc").Limit(pageSize).Offset(offset).Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}
