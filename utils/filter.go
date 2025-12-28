package utils

import "net"

// FilterPrivateIPs 过滤内网IP地址
func FilterPrivateIPs(ips []string) []string {
	var filtered []string
	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			continue
		}

		// 检查是否为私有IP
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified() {
			continue
		}

		// 检查是否为特殊地址
		if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			continue
		}

		filtered = append(filtered, ipStr)
	}
	return filtered
}

// GetMapKeys 获取map的所有键
func GetMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
