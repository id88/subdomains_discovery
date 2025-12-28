package dns

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// LookupResult DNS查询结果
type LookupResult struct {
	IPAddresses []string
	CNAME       string
	TTL         uint32
	RecordType  string
	DNSServer   string
}

// Lookup 执行DNS查询
func (c *Client) Lookup(domain string) (*LookupResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	resultsChan := make(chan *LookupResult, len(c.Servers))
	errChan := make(chan error, len(c.Servers))

	var wg sync.WaitGroup
	for _, server := range c.Servers {
		wg.Add(1)
		go func(srv string) {
			defer wg.Done()

			start := time.Now()
			result, err := c.querySingleServer(srv, domain)
			elapsed := time.Since(start)

			// 更新统计
			stat := c.GetStat(srv)
			if err == nil && result != nil && len(result.IPAddresses) > 0 {
				stat.RecordSuccess(elapsed)
				select {
				case resultsChan <- result:
				default:
					// 已有成功结果
				}
			} else {
				stat.RecordFailure(elapsed)
				errChan <- err
			}
		}(server)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	select {
	case result := <-resultsChan:
		return result, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("查询超时")
	case err := <-errChan:
		// 所有服务器都失败
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("所有DNS服务器查询失败")
	}
}

// querySingleServer 单DNS服务器查询
func (c *Client) querySingleServer(server, domain string) (*LookupResult, error) {
	client := &dns.Client{
		Timeout: c.Timeout,
		Net:     "udp",
	}

	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	// 查询A记录
	resp, _, err := client.Exchange(msg, net.JoinHostPort(server, "53"))
	if err != nil {
		return nil, err
	}

	if resp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("DNS错误: %s", dns.RcodeToString[resp.Rcode])
	}

	result := &LookupResult{
		DNSServer: server,
	}

	for _, ans := range resp.Answer {
		switch rr := ans.(type) {
		case *dns.A:
			result.IPAddresses = append(result.IPAddresses, rr.A.String())
			if result.TTL == 0 {
				result.TTL = rr.Hdr.Ttl
				result.RecordType = "A"
			}
		case *dns.CNAME:
			result.CNAME = strings.TrimSuffix(rr.Target, ".")
			result.TTL = rr.Hdr.Ttl
			result.RecordType = "CNAME"

			// 如果找到CNAME，需要解析CNAME的目标
			target := strings.TrimSuffix(rr.Target, ".")
			cnameResult, err := c.querySingleServer(server, target)
			if err == nil && cnameResult != nil && len(cnameResult.IPAddresses) > 0 {
				result.IPAddresses = append(result.IPAddresses, cnameResult.IPAddresses...)
				if result.TTL > cnameResult.TTL {
					result.TTL = cnameResult.TTL
				}
			}
		case *dns.AAAA:
			// 忽略IPv6，只关注IPv4
			if len(result.IPAddresses) == 0 && result.TTL == 0 {
				result.TTL = rr.Hdr.Ttl
				result.RecordType = "AAAA"
			}
		}
	}

	if len(result.IPAddresses) == 0 {
		return nil, fmt.Errorf("未找到A记录")
	}

	return result, nil
}
