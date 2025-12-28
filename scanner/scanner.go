package scanner

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"subdomains_discovery/dns"
	"subdomains_discovery/output"
	"subdomains_discovery/utils"
	"subdomains_discovery/wildcard"
)

// Scanner 扫描器结构体
type Scanner struct {
	domain      string
	wordlist    string
	dnsClient   *dns.Client
	detector    *wildcard.Detector
	concurrency int
	results     []output.Result
	resultsMu   sync.Mutex
}

// NewScanner 创建新的扫描器
func NewScanner(domain, wordlist string, dnsClient *dns.Client, detector *wildcard.Detector, concurrency int) *Scanner {
	return &Scanner{
		domain:      domain,
		wordlist:    wordlist,
		dnsClient:   dnsClient,
		detector:    detector,
		concurrency: concurrency,
		results:     make([]output.Result, 0),
	}
}

// ReadWordlist 读取字典文件
func (s *Scanner) ReadWordlist() ([]string, error) {
	file, err := os.Open(s.wordlist)
	if err != nil {
		return nil, fmt.Errorf("无法打开字典文件 %s: %v", s.wordlist, err)
	}
	defer file.Close()

	var subdomains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			subdomains = append(subdomains, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取字典文件失败: %v", err)
	}

	if len(subdomains) == 0 {
		return nil, fmt.Errorf("字典文件为空")
	}

	return subdomains, nil
}

// Scan 执行扫描
func (s *Scanner) Scan(subdomains []string) {
	semaphore := make(chan struct{}, s.concurrency)
	var wg sync.WaitGroup

	// 创建进度条
	progress := utils.NewProgressBar(len(subdomains))

	for _, sub := range subdomains {
		wg.Add(1)
		go func(subdomain string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			fullDomain := subdomain + "." + s.domain
			s.scanSubdomain(fullDomain)

			// 更新进度
			progress.Increment()
			if progress.Current%10 == 0 || progress.Current == progress.Total {
				progress.Display()
			}
		}(sub)
	}

	wg.Wait()
	progress.Finish()
}

// scanSubdomain 扫描单个子域名
func (s *Scanner) scanSubdomain(domain string) {
	result, err := s.dnsClient.Lookup(domain)
	if err != nil || result == nil || len(result.IPAddresses) == 0 {
		return
	}

	// 检查是否为通配符
	if s.detector.IsWildcardResult(result.IPAddresses) {
		return
	}

	// 过滤内网IP
	filteredIPs := utils.FilterPrivateIPs(result.IPAddresses)
	if len(filteredIPs) == 0 {
		return
	}

	// 保存结果
	s.resultsMu.Lock()
	s.results = append(s.results, output.Result{
		Subdomain:   domain,
		IPAddresses: filteredIPs,
		CNAME:       result.CNAME,
		TTL:         result.TTL,
		RecordType:  result.RecordType,
		IsWildcard:  false,
		DNSServer:   result.DNSServer,
	})
	s.resultsMu.Unlock()

	fmt.Printf("\n[+] %s -> %v\n", domain, filteredIPs)
}

// GetResults 获取扫描结果
func (s *Scanner) GetResults() []output.Result {
	s.resultsMu.Lock()
	defer s.resultsMu.Unlock()
	return s.results
}
