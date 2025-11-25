package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GeneratePSK 生成 32 字节随机密钥，并使用 Base64 编码。
func GeneratePSK() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate psk: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}
