package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hecate/internal/app/router"
	"hecate/internal/app/service"
	"hecate/internal/app/store"
	"hecate/internal/pkg/config"
	"hecate/internal/pkg/database"
	"hecate/internal/pkg/logger"
)

func main() {
	// 加载配置文件
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
	}
	// 加载日志系统
	if err := logger.Init(&cfg.Log); err != nil {
		fmt.Println("Error initializing logger:", err)
	}
	log := logger.GetLogger() // 获取 logger 实例
	// 初始化数据库
	log.Info("Initializing database connection...")
	if _, err := database.InitDB(&cfg.Database); err != nil {
		log.WithError(err).Fatal("Failed to initialize database")
	}

	db := database.GetDB()
	projectStore := store.NewProjectStore(db)
	projectService := service.NewProjectService(projectStore)
	r := gin.Default()

	router.RegisterRoutes(r, projectService, log)
	log.Info("API routes registered.")

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Infof("Starting web server")

}
