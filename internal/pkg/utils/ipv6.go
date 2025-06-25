package utils

import "net"

func IsIPv6(ipStr string) bool {
	ip := net.ParseIP(ipStr) // 尝试解析字符串为 IP 地址
	if ip == nil {
		return false // 如果解析失败，则不是一个有效的 IP 地址
	}
	return ip.To4() == nil && ip.To16() != nil
}
