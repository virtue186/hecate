package worker

import (
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"hecate/internal/app/store"
	"hecate/internal/pkg/config"
)

// TaskProcessor 所有任务处理器需要的依赖
type TaskProcessor struct {
	Log            *logrus.Logger
	Cfg            *config.Config
	ProjectStore   store.ProjectStore
	AssetStore     store.AssetStore
	AsynqClient    *asynq.Client
	DnsRecordStore store.DnsRecordStore
	PortStore      store.PortStore
}
