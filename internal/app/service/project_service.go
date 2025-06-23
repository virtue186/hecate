package service

import (
	"github.com/hibiken/asynq"
	"github.com/lithammer/shortuuid/v4"
	"hecate/internal/app/dto"
	"hecate/internal/app/model"
	"hecate/internal/app/store"
	"hecate/internal/app/tasks"
	"hecate/internal/pkg/logger"
	"hecate/internal/pkg/utils"
	"time"
)

type ProjectService interface {
	CreateProject(req *dto.CreateProjectRequest) (*dto.ProjectResponse, error)
	GetProjectByID(shortID string) (*dto.ProjectResponse, error)     // 新增
	ListProjects(page, pageSize int) (*dto.PaginatedResponse, error) // 新增

}

func NewProjectService(store store.ProjectStore, client *asynq.Client) ProjectService {
	return &projectService{
		store:       store,
		asynqClient: client, // 新增
	}

}

type projectService struct {
	store       store.ProjectStore
	asynqClient *asynq.Client
}

func (s *projectService) CreateProject(req *dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	// 1. 将 DTO 转换为 Model
	project := &model.Project{
		Name:          req.Name,
		Description:   req.Description,
		ProjectSource: req.ProjectSource,
		Rules:         req.Rules,
		Targets:       make([]model.Target, 0, len(req.Targets)+len(req.ExcludedTargets)),
	}

	for _, t := range req.Targets {
		project.Targets = append(project.Targets, model.Target{Value: t, IsExcluded: false})
	}
	for _, t := range req.ExcludedTargets {
		project.Targets = append(project.Targets, model.Target{Value: t, IsExcluded: true})
	}

	// 2. 调用 Store 层进行数据库操作
	if err := s.store.Create(project); err != nil {
		// 这里可以根据 error 类型进行不同的处理, 比如判断是否是唯一键冲突
		return nil, err
	}

	log := logger.GetLogger()
	for _, target := range project.Targets {
		// 不为排除目标分派任务
		if target.IsExcluded {
			continue
		}

		targetType := utils.DetermineTargetType(target.Value)
		var task *asynq.Task
		var err error

		switch targetType {
		case utils.Domain:
			log.Infof("Dispatching subdomain discovery for domain: %s", target.Value)
			task, err = tasks.NewSubdomainDiscoveryTask(target.ID)
		case utils.IP, utils.CIDR:
			log.Infof("Dispatching port scan for IP/CIDR: %s", target.Value)
			task, err = tasks.NewPortScanTask(target.ID)
		default:
			log.Warnf("Unknown target type for value: %s, skipping task dispatch.", target.Value)
		}

		if err != nil {
			log.WithError(err).Errorf("Failed to create task for target: %s", target.Value)
			continue // 继续处理下一个目标
		}

		if task != nil {
			info, err := s.asynqClient.Enqueue(task)
			if err != nil {
				log.WithError(err).Errorf("Failed to enqueue task for target: %s", target.Value)
			} else {
				log.Infof("Enqueued task ID: %s for target: %s (%s)", info.ID, target.Value, info.Type)
			}
		}
	}

	return modelToResponse(project), nil
}

func modelToResponse(p *model.Project) *dto.ProjectResponse {
	resp := &dto.ProjectResponse{
		ID:            shortuuid.DefaultEncoder.Encode(p.ID), // 编码
		Name:          p.Name,
		Description:   p.Description,
		ProjectSource: p.ProjectSource,
		Rules:         p.Rules,
		CreatedAt:     p.CreatedAt.Format(time.RFC3339),
		Targets:       make([]dto.TargetResponse, len(p.Targets)),
	}
	for i, t := range p.Targets {
		resp.Targets[i] = dto.TargetResponse{
			Value:      t.Value,
			IsExcluded: t.IsExcluded,
		}
	}
	return resp

}

func (s *projectService) GetProjectByID(shortID string) (*dto.ProjectResponse, error) {
	// 1. 解码 ShortUUID -> UUID
	id, err := shortuuid.DefaultEncoder.Decode(shortID)
	if err != nil {
		return nil, err
	}

	// 2. 调用 Store
	project, err := s.store.FindByID(id)
	if err != nil {
		return nil, err // 将 gorm.ErrRecordNotFound 等错误直接传递上去
	}

	// 3. 转换模型为响应 DTO
	return modelToResponse(project), nil
}

func (s *projectService) ListProjects(page, pageSize int) (*dto.PaginatedResponse, error) {
	projects, total, err := s.store.FindAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	// 将 []model.Project 转换为 []dto.ProjectResponse
	responses := make([]*dto.ProjectResponse, len(projects))
	for i, p := range projects {
		responses[i] = modelToResponse(&p)
	}

	return &dto.PaginatedResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Data:     responses,
	}, nil
}
