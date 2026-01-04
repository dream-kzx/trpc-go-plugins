package registry

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"trpc.group/trpc-go/trpc-go/errs"
	trpcRegistry "trpc.group/trpc-go/trpc-go/naming/registry"
)

type NacosRegistry struct {
	client       naming_client.INamingClient
	opt          *Options
	localAddress localAddress
}

func NewNacosRegistry(opt *Options) (*NacosRegistry, error) {
	sConfigs := make([]constant.ServerConfig, 0, len(opt.ServerConfigs))
	for _, conf := range opt.ServerConfigs {
		var sConfig constant.ServerConfig
		if conf.GrpcPort != 0 {
			sConfig = *constant.NewServerConfig(
				conf.Ip, conf.Port, constant.WithGrpcPort(conf.GrpcPort))
		} else {
			sConfig = *constant.NewServerConfig(
				conf.Ip, conf.Port)
		}

		sConfigs = append(sConfigs, sConfig)
	}

	cConfig := constant.NewClientConfig(
		constant.WithNamespaceId(opt.Service.NamespaceId),
		constant.WithTimeoutMs(opt.Timeout),
		constant.WithUsername(opt.Username),
		constant.WithPassword(opt.Password),
		constant.WithLogDir(opt.LogDir),
		constant.WithCacheDir(opt.CacheDir),
	)

	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  cConfig,
		ServerConfigs: sConfigs,
	})
	if err != nil {
		return nil, err
	}

	return &NacosRegistry{client: client, opt: opt}, nil
}

func (r *NacosRegistry) Register(service string, opt ...trpcRegistry.Option) error {
	options := &trpcRegistry.Options{}
	for _, o := range opt {
		o(options)
	}

	if options.Address == "" {
		return errs.NewFrameError(errs.ErrorTypeCalleeFramework, "nacos registry address is empty")
	}

	parts := strings.Split(options.Address, ":")
	if len(parts) != 2 {
		return errs.NewFrameError(errs.ErrorTypeCalleeFramework, "invalid address format")
	}

	port, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return errs.NewFrameError(errs.ErrorTypeCalleeFramework, fmt.Sprintf("invalid port: %v", err))
	}

	param := vo.RegisterInstanceParam{
		Ip:          parts[0],
		Port:        port,
		ServiceName: service,
		GroupName:   r.opt.Service.Group,
		ClusterName: r.opt.Service.Cluster,
		Weight:      r.opt.Service.DefaultWeight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    r.opt.Service.Metadata,
	}

	success, err := r.client.RegisterInstance(param)
	if err != nil {
		return err
	}
	if !success {
		return errs.NewFrameError(errs.ErrorTypeCalleeFramework, "register instance failed")
	}

	r.localAddress = localAddress{
		ip:   parts[0],
		port: port,
	}
	return nil
}

func (r *NacosRegistry) Deregister(service string) error {
	success, err := r.client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          r.localAddress.ip,
		Port:        r.localAddress.port,
		Cluster:     r.opt.Service.Cluster,
		ServiceName: service,
		GroupName:   r.opt.Service.Group,
		Ephemeral:   true,
	})
	if err != nil {
		return err
	}
	if !success {
		return errs.NewFrameError(errs.ErrorTypeCalleeFramework, "deregister instance failed")
	}
	return nil
}

type localAddress struct {
	ip   string
	port uint64
}
