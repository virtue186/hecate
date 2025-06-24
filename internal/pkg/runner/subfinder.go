package runner

import (
	"context"
	"fmt"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"hecate/internal/pkg/config"
	"hecate/internal/pkg/logger"
)

func RunSubfinder(domain string, cfg *config.Config) (map[string]map[string]struct{}, error) {
	log := logger.GetLogger()
	// 配置 subfinder runner
	options := &runner.Options{
		Threads:            cfg.Tools.Subfinder.Threads,
		Timeout:            cfg.Tools.Subfinder.Timeout,
		MaxEnumerationTime: cfg.Tools.Subfinder.MaxEnumerationTime,
		ProviderConfig:     cfg.Tools.Subfinder.ProviderConfigFile,
		All:                cfg.Tools.Subfinder.AllSources,
		RemoveWildcard:     true,
	}

	subfinderRunner, err := runner.NewRunner(options)
	if err != nil {
		// 3. [修正] 修正了原始代码中的致命错误：创建失败必须立刻返回
		log.Errorf("failed to create subfinder runner: %v", err)
		return nil, fmt.Errorf("failed to create subfinder runner: %w", err)
	}

	sourceMap, err := subfinderRunner.EnumerateSingleDomainWithCtx(context.Background(), domain, nil)
	if err != nil {
		log.Errorf("subfinder enumeration failed for domain %s: %v", domain, err)
		return nil, fmt.Errorf("subfinder enumeration failed: %w", err)
	}

	return sourceMap, nil
}
