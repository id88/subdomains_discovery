package output

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// SaveToCSV 保存结果到CSV文件
func SaveToCSV(results []Result, filename string) error {
	if len(results) == 0 {
		return fmt.Errorf("没有结果可保存")
	}

	// #nosec G304 -- filename 来自命令行参数，用户有意控制输出路径
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("无法创建输出文件: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入标题行
	headers := []string{
		"subdomain",
		"ip_address",
		"cname",
		"ttl",
		"record_type",
		"wildcard",
		"dns_server",
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("写入CSV标题失败: %v", err)
	}

	// 写入数据行
	for _, result := range results {
		ipStr := strings.Join(result.IPAddresses, ";")
		wildcardStr := "false"
		if result.IsWildcard {
			wildcardStr = "true"
		}

		record := []string{
			result.Subdomain,
			ipStr,
			result.CNAME,
			fmt.Sprintf("%d", result.TTL),
			result.RecordType,
			wildcardStr,
			result.DNSServer,
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("写入CSV数据失败: %v", err)
		}
	}

	return nil
}
