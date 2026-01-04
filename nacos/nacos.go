package nacos

import "trpc.group/trpc-go/trpc-go/plugin"

const (
	PluginTypeRegistry = "registry"
	PluginTypeSelector = "selector"
	PluginNameNacos    = "nacos"
)

func init() {
	plugin.Register(PluginNameNacos, &RegistryFactory{})
	plugin.Register(PluginNameNacos, &SelectorFactory{})
	plugin.Register(PluginTypeRegistry, &RegistryFactory{})
}
