package nacos

import (
	trpcDiscovery "trpc.group/trpc-go/trpc-go/naming/discovery"
	"trpc.group/trpc-go/trpc-go/plugin"

	"github.com/dream-kzx/trpc-go-plugins/nacos/config"
	"github.com/dream-kzx/trpc-go-plugins/nacos/selector"
)

type DiscoveryFactory struct {
}

func (s DiscoveryFactory) Type() string {
	return PluginTypeSelector
}

func (s DiscoveryFactory) Setup(name string, dec plugin.Decoder) error {
	var cfg config.SelectorConfig
	if err := dec.Decode(&cfg); err != nil {
		return err
	}

	dis, err := selector.NewNacosSelector(&cfg)
	if err != nil {
		return err
	}

	trpcDiscovery.Register(name, dis)
	return nil
}
