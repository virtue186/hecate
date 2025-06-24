package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"hecate/internal/app/tasks"
	"time"
)

func (p *TaskProcessor) HandlePortScanTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks.PortScanPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	p.Log.Infof("[Placeholder] Starting port scan for asset ID: %s", payload.AssetID)
	time.Sleep(5 * time.Second) // 模拟耗时
	p.Log.Infof("[Placeholder] Finished port scan for asset ID: %s", payload.AssetID)

	return nil
}
