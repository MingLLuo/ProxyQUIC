package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// ServerConfig 配置文件
type ServerConfig struct {
	Description string `json:"description"`
	ServerAddr  string `json:"server_address"`
	Http2Addr   string `json:"http2_address"`
	Http3Addr   string `json:"http3_address"`
	UseHTTPS    bool   `json:"use_https"`
}

// LoadServerConfig 从指定文件读取并解析配置
func LoadServerConfig(path string) (*ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}

	var cfg ServerConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config file error: %w", err)
	}
	return &cfg, nil
}
