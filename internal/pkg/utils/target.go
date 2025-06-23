package utils

import (
	"net"
	"regexp"
)

// TargetType 定义了目标类型的枚举
type TargetType int

const (
	Unknown TargetType = iota
	Domain
	IP
	CIDR
	// CompanyName 等可以后续添加
)

// 正则表达式用于匹配域名
var domainRegex = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

// DetermineTargetType 判断给定字符串的目标类型
func DetermineTargetType(target string) TargetType {
	// 尝试解析为 IP 地址
	ip := net.ParseIP(target)
	if ip != nil {
		return IP
	}

	// 尝试解析为 CIDR
	_, _, err := net.ParseCIDR(target)
	if err == nil {
		return CIDR
	}

	// 尝试匹配为域名
	if domainRegex.MatchString(target) {
		return Domain
	}

	return Unknown
}
