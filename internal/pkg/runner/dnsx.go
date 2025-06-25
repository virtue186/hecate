package runner

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	"hecate/internal/pkg/config"
	"hecate/internal/pkg/logger"
)

// DnsxResult 包含了解析出的A记录和CNAME记录
type DnsxResult struct {
	A     []string
	AAAA  []string
	CNAME []string
}

// RunDnsx 执行 dnsx 并返回解析结果
func RunDnsx(domain string, cfg *config.Config) (*DnsxResult, error) {
	log := logger.GetLogger()

	// 配置 dnsx 客户端选项
	// Configure dnsx client options.
	options := dnsx.Options{
		BaseResolvers: cfg.Tools.Dnsx.Resolvers,
		MaxRetries:    cfg.Tools.Dnsx.Retries,
		QuestionTypes: []uint16{
			dns.TypeA,
			dns.TypeAAAA,
			dns.TypeCNAME,
		},
	}
	// 创建一个新的 dnsx 客户端
	// Create a new dnsx client.
	dnsClient, err := dnsx.New(options)
	if err != nil {
		log.Errorf("failed to create dnsx client: %v", err)
		return nil, fmt.Errorf("failed to create dnsx client: %w", err)
	}

	// 使用 QueryOne 方法进行单个域名的查询
	// Use the QueryOne method for a single domain lookup.
	rawResp, err := dnsClient.QueryOne(domain)
	if err != nil {
		log.Errorf("dnsx query failed for domain %s: %v", domain, err)
		return nil, fmt.Errorf("dnsx query failed for domain %s: %w", domain, err)
	}

	if rawResp == nil {
		log.Infof("no dns records found for domain %s", domain)
		return &DnsxResult{}, nil // 返回空结果而不是 nil
	}

	// 从 dnsx 的结果结构体中填充我们自己的结果结构体
	// Populate our result struct from the dnsx result struct.
	result := &DnsxResult{
		A:     rawResp.A,
		AAAA:  rawResp.AAAA,
		CNAME: rawResp.CNAME,
	}

	return result, nil
}
