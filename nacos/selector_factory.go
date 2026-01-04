package nacos

import (
	trpcSelector "trpc.group/trpc-go/trpc-go/naming/selector"
	"trpc.group/trpc-go/trpc-go/plugin"

	"github.com/dream-kzx/trpc-go-plugins/nacos/config"
	"github.com/dream-kzx/trpc-go-plugins/nacos/selector"
)

type SelectorFactory struct {
}

func (s SelectorFactory) Type() string {
	return PluginTypeSelector
}

func (s SelectorFactory) Setup(name string, dec plugin.Decoder) error {
	var cfg config.SelectorConfig
	if err := dec.Decode(&cfg); err != nil {
		return err
	}

	sl, err := selector.NewNacosSelector(&cfg)
	if err != nil {
		return err
	}

	trpcSelector.Register(name, sl)
	return nil
}
