package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hecate/internal/app/dto"
	"hecate/internal/app/service"
	"hecate/internal/pkg/response"
	"strconv"
)

type ProjectHandler struct {
	service service.ProjectService
	log     *logrus.Logger
}

func NewProjectHandler(service service.ProjectService, log *logrus.Logger) *ProjectHandler {
	return &ProjectHandler{service: service, log: log}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	// 1. 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Errorf("Failed to bind request: %v", err)
		response.ValidationError(c, "")
		return
	}

	// 2. 调用 service 层处理业务逻辑
	projectResp, err := h.service.CreateProject(&req)
	if err != nil {
		h.log.Errorf("Failed to create project: %v", err)
		// 在这里可以更精细地处理错误，比如判断是否是数据库唯一键冲突
		response.InternalError(c, "")
		return
	}

	h.log.Infof("Project %s created successfully", projectResp.Name)
	// 3. 返回成功响应
	response.Created(c, projectResp)
}

func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	// 1. 从 URL 路径中获取 ID
	shortID := c.Param("id")

	// 2. 调用 Service
	project, err := h.service.GetProjectByID(shortID)
	if err != nil {
		// 判断错误类型
		if err == gorm.ErrRecordNotFound {
			h.log.Warnf("Project with short_id %s not found", shortID)
			response.ValidationError(c, "Projectid not found")
			return
		}
		// 如果是 shortuuid 的解码错误，通常是格式问题
		// 为了简化，我们统一返回 400
		h.log.Errorf("Failed to get project for short_id %s: %v", shortID, err)
		response.ValidationError(c, "")
		return
	}

	// 3. 返回成功响应
	response.Success(c, project)
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	// 1. 从查询参数中获取分页信息
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // 设置默认值和最大值
	}

	// 2. 调用 Service
	paginatedResponse, err := h.service.ListProjects(page, pageSize)
	if err != nil {
		h.log.Errorf("Failed to list projects: %v", err)
		response.ValidationError(c, "")
		return
	}

	// 3. 返回成功响应
	response.Success(c, paginatedResponse)
}
