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
	dnsRecordStore := store.NewDnsRecordStore(db)
	portStore := store.NewPortStore(db)

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
			Concurrency: 10,
			Logger:      log,
			// Queues 定义了不同队列的优先级。数字越大，优先级越高。
			// Worker 会优先从高优先级的队列中取任务。
			Queues: map[string]int{
				"critical":  6, // 比如用于处理紧急的0day扫描任务
				"default":   3, // 用于处理常规的、轻量级的任务
				"discovery": 1, // 优先级最低，用于处理耗时的初始发现任务
			},
		},
	)

	// 创建并注入依赖到 TaskProcessor
	processor := &worker.TaskProcessor{
		Log:            log,
		Cfg:            cfg,
		ProjectStore:   projectStore,
		AssetStore:     assetStore,
		AsynqClient:    asynqClient,
		DnsRecordStore: dnsRecordStore,
		PortStore:      portStore,
	}

	// 注册所有任务处理器
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeSubdomainDiscovery, processor.HandleSubdomainDiscoveryTask)
	mux.HandleFunc(tasks.TypeDnsResolve, processor.HandleDnsResolveTask)
	mux.HandleFunc(tasks.TypePortScan, processor.HandlePortScanTask)

	log.Info("Starting Asynq worker...")
	if err := asynqServer.Run(mux); err != nil {
		log.Fatalf("Could not run Asynq worker: %v", err)
	}

}
