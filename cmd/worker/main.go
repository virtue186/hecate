package main

import (
	"github.com/hibiken/asynq"
	"hecate/internal/app/tasks"
	"hecate/internal/app/worker"
	"hecate/internal/pkg/config"
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

	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				"default": 10, // "default" 队列的并发度为 10
			},
		},
	)

	processor := &worker.TaskProcessor{Log: log}
	mux := asynq.NewServeMux()

	// [修改] 注册所有任务处理器
	mux.HandleFunc(tasks.TypeSubdomainDiscovery, processor.HandleSubdomainDiscoveryTask)
	mux.HandleFunc(tasks.TypePortScan, processor.HandlePortScanTask) // 新增注册

	log.Info("Starting Asynq worker...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Could not run Asynq worker: %v", err)
	}

}
