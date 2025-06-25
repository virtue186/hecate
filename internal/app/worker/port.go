package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"hecate/internal/app/model"
	"hecate/internal/app/tasks"
	"hecate/internal/pkg/runner"
)

func (p *TaskProcessor) HandlePortScanTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks.PortScanPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	p.Log.Infof("Starting port scan for asset ID: %s", payload.AssetID)

	// 1. 获取资产信息
	asset, err := p.AssetStore.FindByID(payload.AssetID)
	if err != nil {
		p.Log.WithError(err).Errorf("Failed to find asset by ID: %s", payload.AssetID)
		return err
	}

	// 2. 调用 Runner 执行扫描
	openPorts, err := runner.RunNaabu(asset.Value, p.Cfg)
	if err != nil {
		p.Log.WithError(err).Errorf("Naabu run failed for asset: %s", asset.Value)
		return err
	}
	p.Log.Infof("Discovered %d open ports for %s", len(openPorts), asset.Value)

	if len(openPorts) == 0 {
		return nil
	}

	// 3. 准备数据模型
	var newPorts []*model.Port
	for _, port := range openPorts {
		newPorts = append(newPorts, &model.Port{
			AssetID:    asset.ID,
			PortNumber: port.Port,
			Protocol:   port.Protocol,
			Service:    "", // 服务名称暂时留空
		})
	}

	// 4. 存入数据库
	if err := p.PortStore.CreateBatch(newPorts); err != nil {
		p.Log.WithError(err).Error("Failed to batch create ports")
		return err
	}

	p.Log.Infof("Successfully saved %d open ports for asset %s", len(newPorts), asset.Value)

	// TODO: 下一步将在这里创建服务扫描任务

	return nil
}
