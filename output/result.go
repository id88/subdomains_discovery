package output

// Result 子域名扫描结果结构体
type Result struct {
	Subdomain   string
	IPAddresses []string
	CNAME       string
	TTL         uint32
	RecordType  string
	IsWildcard  bool
	DNSServer   string
}
