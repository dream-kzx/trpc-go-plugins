package registry

import "github.com/dream-kzx/trpc-go-plugins/nacos/config"

type Options struct {
	ServerConfigs []config.ServerConfig `yaml:"server_configs"`
	LogDir        string
	CacheDir      string
	Username      string
	Password      string
	Timeout       uint64
	Service       Service
}

type Service struct {
	Name          string
	NamespaceId   string
	Group         string
	Cluster       string
	DefaultWeight float64
	Metadata      map[string]string
}
