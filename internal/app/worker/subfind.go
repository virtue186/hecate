package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"hecate/internal/app/model"
	"hecate/internal/app/tasks"
	"hecate/internal/pkg/runner"
	"strings"
)

// HandleSubdomainDiscoveryTask 是子域名发现任务的具体处理逻辑
func (p *TaskProcessor) HandleSubdomainDiscoveryTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks.SubdomainDiscoveryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	p.Log.Infof("Starting subdomain discovery for target ID: %s", payload.TargetID)

	// 1. 从数据库获取目标信息FindTargetByID
	target, err := p.ProjectStore.FindTargetByID(payload.TargetID)
	if err != nil {
		p.Log.WithError(err).Error("Failed to find target by ID")
		return err
	}

	// 2. 调用 runner 执行 subfinder
	sourceMap, err := runner.RunSubfinder(target.Value, p.Cfg)
	if err != nil {
		p.Log.WithError(err).Errorf("Subfinder run failed for domain: %s", target.Value)
		return err
	}
	p.Log.Infof("Discovered %d unique subdomains for %s", len(sourceMap), target.Value)

	if len(sourceMap) == 0 {
		p.Log.Info("No subdomains found, task finished.")
		return nil
	}

	// 3. 将结果转换为 Asset 模型
	var newAssets []*model.Asset
	for host, sources := range sourceMap {
		var sourceSlice []string
		for sourceName := range sources {
			sourceSlice = append(sourceSlice, sourceName)
		}
		newAssets = append(newAssets, &model.Asset{
			ProjectID: target.ProjectID,
			Value:     host,
			Type:      model.AssetTypeSubdomain,
			Source:    strings.Join(sourceSlice, ","),
		})
	}

	// 4. 批量存入数据库 (GORM会回填ID)
	if err := p.AssetStore.CreateBatch(newAssets); err != nil {
		p.Log.WithError(err).Error("Failed to batch create assets")
		return err
	}
	p.Log.Infof("Successfully saved assets for project %s", target.ProjectID)

	// 5. 任务链：为新发现的资产创建下一阶段任务
	enqueuedCount := 0
	for _, asset := range newAssets {
		if asset.ID == uuid.Nil {
			continue
		}
		portScanTask, err := tasks.NewPortScanTask(asset.ID, "default")
		if err != nil {
			p.Log.WithError(err).Errorf("Failed to create port scan task for asset %s", asset.Value)
			continue
		}
		_, err = p.AsynqClient.Enqueue(portScanTask)
		if err != nil {
			p.Log.WithError(err).Errorf("Failed to enqueue port scan task for asset %s", asset.Value)
			continue
		}
		enqueuedCount++
	}

	p.Log.Infof("Enqueued %d port scan tasks for newly discovered assets.", enqueuedCount)

	return nil
}
