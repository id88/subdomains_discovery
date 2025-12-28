package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

// RandomString 生成指定长度的随机字符串
func RandomString(length int) string {
	b := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range b {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			// 如果随机数生成失败，使用默认字符
			b[i] = charset[0]
			continue
		}
		b[i] = charset[num.Int64()]
	}

	return string(b)
}
