package config

type RegistryConfig struct {
	ServerConfigs []ServerConfig `yaml:"server_configs"`

	Username            string    `yaml:"username"`
	Password            string    `yaml:"password"`
	Timeout             uint64    `yaml:"timeout"`
	LogDir              string    `yaml:"log_dir"`
	CacheDir            string    `yaml:"cache_dir"`
	NotLoadCacheAtStart bool      `yaml:"not_load_cache_at_start"`
	Services            []Service `yaml:"services"`
}

type ServerConfig struct {
	Ip       string `yaml:"ip"`
	Port     uint64 `yaml:"port"`
	GrpcPort uint64 `yaml:"grpc_port"`
}

type Service struct {
	Name          string            `yaml:"name"`
	NamespaceId   string            `yaml:"namespace_id"`
	Group         string            `yaml:"group"`
	Cluster       string            `yaml:"cluster"`
	DefaultWeight float64           `yaml:"default_weight"`
	Metadata      map[string]string `yaml:"metadata"`
}

type SelectorConfig struct {
	ServerConfigs   []ServerConfig `yaml:"server_configs"`
	NamespaceId     string         `yaml:"namespace_id"`
	Group           string         `yaml:"group"`
	Cluster         string         `yaml:"cluster"`
	Timeout         uint64         `yaml:"timeout"`
	LoadBalanceType string         `yaml:"load_balance_type"`
	Username        string         `yaml:"username"`
	Password        string         `yaml:"password"`
	LogDir          string         `yaml:"log_dir"`
	CacheDir        string         `yaml:"cache_dir"`
}
