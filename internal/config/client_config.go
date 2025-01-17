package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// ClientConfig 配置文件
type ClientConfig struct {
	Description   string `json:"description"`
	ClientAddr    string `json:"client_address"`
	ServerAddr    string `json:"server_address"`
	ProxyAddr     string `json:"proxy_address"`
	ClientMessage string `json:"client_message"`
	UseHTTPS      bool   `json:"use_https"`
}

// LoadClientConfig 从指定文件读取并解析配置
func LoadClientConfig(path string) (*ClientConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}

	var cfg ClientConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config file error: %w", err)
	}
	return &cfg, nil
}
