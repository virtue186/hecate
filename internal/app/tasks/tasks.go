package tasks

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// NewSubdomainDiscoveryTask 创建一个新的子域名发现任务
func NewSubdomainDiscoveryTask(targetID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(SubdomainDiscoveryPayload{TargetID: targetID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSubdomainDiscovery, payload), nil
}

// NewPortScanTask 创建一个新的端口扫描任务
func NewPortScanTask(targetID uuid.UUID) (*asynq.Task, error) { // 新增
	payload, err := json.Marshal(PortScanPayload{TargetID: targetID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePortScan, payload), nil
}
