package runner

import (
	"context"
	"fmt"
	"github.com/projectdiscovery/naabu/v2/pkg/result"
	"github.com/projectdiscovery/naabu/v2/pkg/runner"
	"hecate/internal/pkg/config"
	"hecate/internal/pkg/logger"
	"sync"
	"time"
)

// PortResult 现在只包含 naabu 直接发现的信息
type PortResult struct {
	Port     int
	Protocol string
}

// RunNaabu 执行 naabu 并返回开放的端口列表
func RunNaabu(host string, cfg *config.Config) ([]PortResult, error) {
	log := logger.GetLogger()
	naabuCfg := cfg.Tools.Naabu

	var results []PortResult
	// [核心修改] 使用互斥锁来保护对 results 切片的并发写入
	var mu sync.Mutex

	// 1. [核心修改] 完全遵照官方示例和文档配置 Options
	options := runner.Options{
		Host:     []string{host},
		Ports:    naabuCfg.Ports,
		Rate:     naabuCfg.Rate,
		Timeout:  time.Duration(naabuCfg.Timeout),
		Retries:  naabuCfg.Retries,
		ScanType: naabuCfg.ScanType,
		Silent:   true, // 不输出 banner 等信息
		// [新增] 应用 CDN 排除配置
		ExcludeCDN: naabuCfg.ExcludeCdn,
		// [核心修改] 使用 OnResult 回调函数来处理结果
		OnResult: func(hr *result.HostResult) {
			mu.Lock() // 在写入共享切片前加锁
			for _, port := range hr.Ports {
				results = append(results, PortResult{
					Port:     port.Port,
					Protocol: port.Protocol.String(),
				})
			}
			mu.Unlock()
		},
	}

	naabuRunner, err := runner.NewRunner(&options)
	if err != nil {
		log.Errorf("failed to create naabu runner: %v", err)
		return nil, fmt.Errorf("failed to create naabu runner: %w", err)
	}
	defer naabuRunner.Close()

	// 2. [核心修改] 调用 RunEnumeration 开始扫描
	if err := naabuRunner.RunEnumeration(context.Background()); err != nil {
		log.Errorf("naabu enumeration failed for host %s: %v", host, err)
		return nil, fmt.Errorf("naabu enumeration failed: %w", err)
	}

	return results, nil
}
