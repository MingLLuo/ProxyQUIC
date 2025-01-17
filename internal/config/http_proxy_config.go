package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// HttpProxyConfig 配置文件
type HttpProxyConfig struct {
	Description string `json:"description"`
	ProxyAddr   string `json:"proxy_address"`
}

// LoadHttpProxyConfig 从指定文件读取并解析配置
func LoadHttpProxyConfig(path string) (*HttpProxyConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}

	var cfg HttpProxyConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config file error: %w", err)
	}
	return &cfg, nil
}
