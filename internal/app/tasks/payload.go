package tasks

import "github.com/google/uuid"

const (
	TypeSubdomainDiscovery = "discovery:subdomain"
	TypePortScan           = "discovery:portscan" // 新增
)

// SubdomainDiscoveryPayload 的载荷现在是 TargetID
type SubdomainDiscoveryPayload struct {
	TargetID uuid.UUID
}

// PortScanPayload 端口扫描任务的载荷
type PortScanPayload struct { // 新增
	AssetID uuid.UUID
}
