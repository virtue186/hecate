package asynq_client

import (
	"github.com/hibiken/asynq"
	"hecate/internal/pkg/config"
)

var client *asynq.Client

// InitClient 初始化任务队列客户端
func InitClient(cfg *config.RedisConfig) {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	client = asynq.NewClient(redisOpt)
}

// GetClient 返回 asynq 客户端实例
func GetClient() *asynq.Client {
	if client == nil {
		panic("Asynq client is not initialized")
	}
	return client
}
