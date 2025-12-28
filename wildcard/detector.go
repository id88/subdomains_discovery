package wildcard

import (
	"fmt"
	"time"

	"subdomains_discovery/dns"
	"subdomains_discovery/utils"
)

// Detector 通配符检测器
type Detector struct {
	domain      string
	dnsClient   *dns.Client
	wildcardIPs map[string]bool
}

// NewDetector 创建新的通配符检测器
func NewDetector(domain string, dnsClient *dns.Client) *Detector {
	return &Detector{
		domain:      domain,
		dnsClient:   dnsClient,
		wildcardIPs: make(map[string]bool),
	}
}

// Detect 检测通配符DNS
func (d *Detector) Detect() {
	fmt.Println("[*] 检测通配符DNS...")

	// 生成10个随机子域名进行测试
	for i := 0; i < 10; i++ {
		randomSub := fmt.Sprintf("%s.%s", utils.RandomString(16), d.domain)
		result, err := d.dnsClient.Lookup(randomSub)

		if err == nil && result != nil {
			for _, ip := range result.IPAddresses {
				d.wildcardIPs[ip] = true
			}
		}

		time.Sleep(50 * time.Millisecond) // 避免触发速率限制
	}

	if len(d.wildcardIPs) > 0 {
		fmt.Printf("[!] 检测到通配符DNS，以下IP将被过滤: %v\n", utils.GetMapKeys(d.wildcardIPs))
	} else {
		fmt.Println("[*] 未检测到通配符DNS")
	}
}

// IsWildcard 检查IP是否为通配符IP
func (d *Detector) IsWildcard(ip string) bool {
	return d.wildcardIPs[ip]
}

// IsWildcardResult 检查结果是否包含通配符IP
func (d *Detector) IsWildcardResult(ips []string) bool {
	for _, ip := range ips {
		if d.wildcardIPs[ip] {
			return true
		}
	}
	return false
}
