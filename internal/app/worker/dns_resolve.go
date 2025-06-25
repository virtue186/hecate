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
	"hecate/internal/pkg/utils"
)

func (p *TaskProcessor) HandleDnsResolveTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks.DnsResolvePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	p.Log.Infof("Starting dns resolve for asset ID: %s", payload.AssetID)

	// 1. 获取域名资产信息
	domainAsset, err := p.AssetStore.FindByID(payload.AssetID)
	if err != nil {
		p.Log.WithError(err).Errorf("Failed to find asset by ID: %s", payload.AssetID)
		return err
	}

	// 2. 调用 runner 执行 dnsx
	dnsResult, err := runner.RunDnsx(domainAsset.Value, p.Cfg)
	if err != nil {
		p.Log.WithError(err).Errorf("Dnsx run failed for asset: %s", domainAsset.Value)
		return err
	}

	var dnsRecords []*model.DNS_Record
	var ipAssets []*model.Asset

	// 3. 处理并保存 CNAME 记录
	for _, cname := range dnsResult.CNAME {
		dnsRecords = append(dnsRecords, &model.DNS_Record{
			AssetID: domainAsset.ID, Host: domainAsset.Value, Type: "CNAME", Value: cname, Source: "dnsx"})
	}

	// 4. 处理并保存 A/AAAA 记录，并创建新的IP资产
	allIPs := append(dnsResult.A, dnsResult.AAAA...)
	for _, ip := range allIPs {
		recordType := "A"
		if utils.IsIPv6(ip) {
			recordType = "AAAA"
		}
		dnsRecords = append(dnsRecords, &model.DNS_Record{
			AssetID: domainAsset.ID, Host: domainAsset.Value, Type: recordType, Value: ip, Source: "dnsx"})

		ipAssets = append(ipAssets, &model.Asset{
			ProjectID: domainAsset.ProjectID, Value: ip, Type: model.AssetTypeIP, Source: "dnsx"})
	}

	// 5. 批量存入数据库
	if err := p.DnsRecordStore.CreateBatch(dnsRecords); err != nil {
		p.Log.WithError(err).Error("Failed to batch create dns records")
	}
	if err := p.AssetStore.CreateBatch(ipAssets); err != nil {
		p.Log.WithError(err).Error("Failed to batch create ip assets")
	}
	p.Log.Infof("Saved %d DNS records and %d IP assets for %s", len(dnsRecords), len(ipAssets), domainAsset.Value)

	// 6. 任务链：为新发现的IP资产分发端口扫描任务
	enqueuedCount := 0
	for _, ipAsset := range ipAssets {
		if ipAsset.ID == uuid.Nil {
			continue
		}
		portScanTask, err := tasks.NewPortScanTask(ipAsset.ID, "default")
		if err != nil {
			continue
		}
		_, err = p.AsynqClient.Enqueue(portScanTask)
		if err != nil {
			continue
		}
		enqueuedCount++
	}
	p.Log.Infof("Enqueued %d port scan tasks for newly discovered IP assets.", enqueuedCount)

	return nil
}
