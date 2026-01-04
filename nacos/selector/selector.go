package selector

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/naming/discovery"
	trpcRegistry "trpc.group/trpc-go/trpc-go/naming/registry"
	trpcSelector "trpc.group/trpc-go/trpc-go/naming/selector"

	"github.com/dream-kzx/trpc-go-plugins/nacos/config"
)

type NacosSelector struct {
	client naming_client.INamingClient
	cfg    *config.SelectorConfig
}

func NewNacosSelector(cfg *config.SelectorConfig) (*NacosSelector, error) {
	sConfigs := make([]constant.ServerConfig, 0, len(cfg.ServerConfigs))
	for _, cnf := range cfg.ServerConfigs {
		if cnf.GrpcPort != 0 {
			sConfigs = append(sConfigs, *constant.NewServerConfig(cnf.Ip, cnf.Port, constant.WithGrpcPort(cnf.GrpcPort)))
		} else {
			sConfigs = append(sConfigs, *constant.NewServerConfig(cnf.Ip, cnf.Port))
		}

	}

	cc := constant.NewClientConfig(
		constant.WithNamespaceId(cfg.NamespaceId),
		constant.WithTimeoutMs(cfg.Timeout),
		constant.WithUsername(cfg.Username),
		constant.WithPassword(cfg.Password),
		constant.WithLogDir(cfg.LogDir),
		constant.WithCacheDir(cfg.CacheDir),
	)

	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  cc,
		ServerConfigs: sConfigs,
	})
	if err != nil {
		return nil, err
	}
	return &NacosSelector{client: client, cfg: cfg}, nil
}

func (n *NacosSelector) Select(serviceName string, opt ...trpcSelector.Option) (*trpcRegistry.Node, error) {
	instances, err := n.client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   n.cfg.Group,
		Clusters:    []string{n.cfg.Cluster},
		HealthyOnly: true,
	})

	if err != nil {
		return nil, errs.NewFrameError(errs.ErrorTypeFramework, fmt.Sprintf("nacos selector error: %v", err))
	}

	if len(instances) == 0 {
		return nil, errs.NewFrameError(errs.ErrorTypeBusiness, "no instance available")
	}

	// 负载均衡示例：随机（可根据 cfg.LoadBalanceType 切换为轮询、权重等）
	var inst model.Instance
	switch n.cfg.LoadBalanceType {
	case "random":
		inst = instances[rand.Intn(len(instances))]
	// 添加其他类型...
	default:
		inst = instances[rand.Intn(len(instances))]
	}

	metadata := make(map[string]interface{})
	for k, v := range inst.Metadata {
		metadata[k] = v
	}

	return &trpcRegistry.Node{
		ServiceName: serviceName,
		Address:     fmt.Sprintf("%s:%d", inst.Ip, inst.Port),
		Weight:      int(inst.Weight),
		Metadata:    metadata,
	}, nil
}

func (n *NacosSelector) Report(node *trpcRegistry.Node, cost time.Duration, err error) error {
	return nil
}

func (n *NacosSelector) List(serviceName string, opt ...discovery.Option) (nodes []*trpcRegistry.Node, err error) {
	instances, err := n.client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   n.cfg.Group,
		Clusters:    []string{n.cfg.Cluster},
		HealthyOnly: true,
	})

	if err != nil {
		return nil, errs.NewFrameError(errs.ErrorTypeFramework, fmt.Sprintf("nacos selector error: %v", err))
	}

	if len(instances) == 0 {
		return []*trpcRegistry.Node{}, nil
	}

	for _, inst := range instances {
		metadata := make(map[string]interface{})
		for k, v := range inst.Metadata {
			metadata[k] = v
		}

		nodes = append(nodes, &trpcRegistry.Node{
			ServiceName: serviceName,
			Address:     fmt.Sprintf("%s:%d", inst.Ip, inst.Port),
			Weight:      int(inst.Weight),
			Metadata:    metadata,
		})
	}

	return nodes, nil
}
