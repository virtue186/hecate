package main

import (
	"github.com/hibiken/asynq"
	"hecate/internal/app/store"
	"hecate/internal/app/tasks"
	"hecate/internal/app/worker"
	"hecate/internal/pkg/asynq_client"
	"hecate/internal/pkg/config"
	"hecate/internal/pkg/database"
	"hecate/internal/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	if err := logger.Init(&cfg.Log); err != nil {
		panic(err)
	}
	log := logger.GetLogger()

	// Worker 也需要初始化数据库连接
	if _, err := database.InitDB(&cfg.Database); err != nil {
		log.WithError(err).Fatal("Failed to initialize database for worker")
	}
	db := database.GetDB()

	// 创建 store 实例
	projectStore := store.NewProjectStore(db)
	assetStore := store.NewAssetStore(db)

	// 为 Worker 创建 Asynq Client 实例 (用于任务链)
	asynq_client.InitClient(&cfg.Redis)
	asynqClient := asynq_client.GetClient()
	defer asynqClient.Close()

	// 创建 Asynq Server 实例 (用于处理任务)
	// 2. 配置并启动 Asynq 服务器
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
	asynqServer := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				"default": 10, // "default" 队列的并发度为 10
			},
		},
	)

	// 创建并注入依赖到 TaskProcessor
	processor := &worker.TaskProcessor{
		Log:          log,
		Cfg:          cfg,
		ProjectStore: projectStore,
		AssetStore:   assetStore,
		AsynqClient:  asynqClient,
	}

	// 注册所有任务处理器
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeSubdomainDiscovery, processor.HandleSubdomainDiscoveryTask)
	mux.HandleFunc(tasks.TypePortScan, processor.HandlePortScanTask) // 注册占位符

	log.Info("Starting Asynq worker...")
	if err := asynqServer.Run(mux); err != nil {
		log.Fatalf("Could not run Asynq worker: %v", err)
	}

}
