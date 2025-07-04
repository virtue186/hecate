package tasks

import "github.com/google/uuid"

const (
	TypeSubdomainDiscovery = "discovery:subdomain"
	TypePortScan           = "discovery:portscan"
	TypeDnsResolve         = "discovery:resolve"
)

// SubdomainDiscoveryPayload 的载荷现在是 TargetID
type SubdomainDiscoveryPayload struct {
	TargetID uuid.UUID
}

// PortScanPayload 端口扫描任务的载荷
type PortScanPayload struct { // 新增
	AssetID uuid.UUID
}

type DnsResolvePayload struct { // 新增
	AssetID uuid.UUID
}
