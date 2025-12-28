package config

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Config 配置结构体
type Config struct {
	Domain      string
	Wordlist    string
	DNSServers  []string
	Output      string
	Timeout     int
	Concurrency int
}

// ParseFlags 解析命令行参数
func ParseFlags() *Config {
	cfg := &Config{}
	var dnsServers string
	var defaultWordlist string

	// 获取当前目录
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	defaultWordlist = cwd + "/dict.txt"

	flag.StringVar(&cfg.Domain, "d", "", "目标域名 (例如: example.com)")
	flag.StringVar(&cfg.Wordlist, "w", defaultWordlist, "字典文件路径 (默认: ./dict.txt)")
	flag.StringVar(&dnsServers, "dnss", "114.114.114.114,223.5.5.5", "DNS服务器，逗号分隔 (默认: 114.114.114.114,223.5.5.5)")
	flag.StringVar(&cfg.Output, "o", "", "输出文件路径 (默认: subdomains_域名_时间戳.csv)")
	flag.IntVar(&cfg.Timeout, "t", 3, "超时时间(秒)")
	flag.IntVar(&cfg.Concurrency, "c", 100, "并发数")

	flag.Parse()

	// 必需参数检查
	if cfg.Domain == "" {
		fmt.Println("错误: 必须指定目标域名")
		fmt.Println("用法:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 处理DNS服务器列表
	cfg.DNSServers = parseDNSServers(dnsServers)
	if len(cfg.DNSServers) == 0 {
		cfg.DNSServers = []string{"114.114.114.114", "223.5.5.5"}
	}

	// 设置默认输出文件名
	if cfg.Output == "" {
		timestamp := time.Now().Format("20060102_150405")
		cfg.Output = fmt.Sprintf("subdomains_%s_%s.csv", cfg.Domain, timestamp)
	}

	// 确保域名没有协议前缀
	cfg.Domain = strings.TrimPrefix(cfg.Domain, "http://")
	cfg.Domain = strings.TrimPrefix(cfg.Domain, "https://")
	cfg.Domain = strings.TrimSuffix(cfg.Domain, "/")

	return cfg
}

// parseDNSServers 解析DNS服务器列表
func parseDNSServers(input string) []string {
	servers := strings.Split(input, ",")
	validServers := make([]string, 0, len(servers))

	for _, server := range servers {
		server = strings.TrimSpace(server)
		if server == "" {
			continue
		}
		// 验证IP地址格式
		if net.ParseIP(server) != nil {
			validServers = append(validServers, server)
		} else {
			fmt.Printf("[!] 警告: 无效的DNS服务器地址: %s\n", server)
		}
	}

	return validServers
}
