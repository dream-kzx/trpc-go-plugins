package nacos

import (
	"fmt"

	"trpc.group/trpc-go/trpc-go/errs"
	trpcRegistry "trpc.group/trpc-go/trpc-go/naming/registry"
	"trpc.group/trpc-go/trpc-go/plugin"

	"github.com/dream-kzx/trpc-go-plugins/nacos/config"
	"github.com/dream-kzx/trpc-go-plugins/nacos/registry"
)

type RegistryFactory struct{}

func (r *RegistryFactory) Type() string {
	return PluginTypeRegistry
}

func (r *RegistryFactory) Setup(name string, dec plugin.Decoder) error {
	var cfg config.RegistryConfig
	if err := dec.Decode(&cfg); err != nil {
		return errs.NewFrameError(errs.ErrorTypeFramework, fmt.Sprintf("nacos registry config error: %v", err))
	}

	for _, service := range cfg.Services {
		opt := registry.Options{
			ServerConfigs: cfg.ServerConfigs,
			LogDir:        cfg.LogDir,
			CacheDir:      cfg.CacheDir,
			Username:      cfg.Username,
			Password:      cfg.Password,
			Timeout:       cfg.Timeout,
			Service: registry.Service{
				Name:          service.Name,
				NamespaceId:   service.NamespaceId,
				Group:         service.Group,
				Cluster:       service.Cluster,
				DefaultWeight: service.DefaultWeight,
				Metadata:      service.Metadata,
			},
		}

		reg, err := registry.NewNacosRegistry(&opt)
		if err != nil {
			return errs.NewFrameError(errs.ErrorTypeFramework, fmt.Sprintf("nacos registry init error: %v", err))
		}
		trpcRegistry.Register(opt.Service.Name, reg)
	}
	return nil
}
