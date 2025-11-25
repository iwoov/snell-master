package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateAPIToken 生成节点使用的 API Token。
func GenerateAPIToken() (string, error) {
	return randomHex(32)
}

// GenerateSubscribeToken 生成订阅访问令牌。
func GenerateSubscribeToken() (string, error) {
	return randomHex(32)
}

func randomHex(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate random hex: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
