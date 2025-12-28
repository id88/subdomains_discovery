package dns

import (
	"sync"
	"time"
)

// Client DNS客户端结构体
type Client struct {
	Servers  []string
	Timeout  time.Duration
	Strategy string // parallel, fallback
	Stats    map[string]*Stat
	statsMu  sync.RWMutex
}

// NewClient 创建新的DNS客户端
func NewClient(servers []string, timeout time.Duration) *Client {
	client := &Client{
		Servers:  servers,
		Timeout:  timeout,
		Strategy: "parallel",
		Stats:    make(map[string]*Stat),
	}

	// 初始化统计信息
	for _, server := range servers {
		client.Stats[server] = &Stat{}
	}

	return client
}

// GetStat 获取指定服务器的统计信息
func (c *Client) GetStat(server string) *Stat {
	c.statsMu.RLock()
	defer c.statsMu.RUnlock()
	return c.Stats[server]
}
