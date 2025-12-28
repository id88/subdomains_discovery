package main

import (
	"fmt"
	"time"

	"subdomains_discovery/config"
	"subdomains_discovery/dns"
	"subdomains_discovery/output"
	"subdomains_discovery/scanner"
	"subdomains_discovery/wildcard"
)

func main() {
	// 解析命令行参数
	cfg := config.ParseFlags()

	// 初始化DNS客户端
	dnsClient := dns.NewClient(cfg.DNSServers, time.Duration(cfg.Timeout)*time.Second)

	// 检测通配符
	detector := wildcard.NewDetector(cfg.Domain, dnsClient)
	detector.Detect()

	// 创建扫描器
	scan := scanner.NewScanner(cfg.Domain, cfg.Wordlist, dnsClient, detector, cfg.Concurrency)

	// 读取字典文件
	subdomains, err := scan.ReadWordlist()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("[*] 开始扫描域名: %s\n", cfg.Domain)
	fmt.Printf("[*] 字典大小: %d\n", len(subdomains))
	fmt.Printf("[*] DNS服务器: %v\n", cfg.DNSServers)
	fmt.Printf("[*] 并发数: %d\n", cfg.Concurrency)
	fmt.Printf("[*] 超时时间: %d秒\n", cfg.Timeout)
	fmt.Println("[*] 扫描开始...")

	startTime := time.Now()

	// 执行扫描
	scan.Scan(subdomains)

	// 获取结果
	results := scan.GetResults()

	// 保存结果
	if len(results) > 0 {
		if err := output.SaveToCSV(results, cfg.Output); err != nil {
			fmt.Printf("错误: %v\n", err)
		} else {
			fmt.Printf("[*] 结果保存在: %s\n", cfg.Output)
		}
	} else {
		fmt.Println("[!] 未发现任何子域名")
	}

	elapsed := time.Since(startTime)
	fmt.Printf("[*] 扫描完成! 耗时: %v\n", elapsed)
	fmt.Printf("[*] 发现子域名: %d 个\n", len(results))
}
