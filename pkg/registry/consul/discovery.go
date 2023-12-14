package consul

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

type RegistryConfig struct {
	Address   string
	Scheme    string
	Discovery string
}

type Discovery struct {
	registry.Discovery
}

func NewDiscovery(conf *RegistryConfig) Discovery {
	c := api.DefaultConfig()
	c.Address = conf.Address
	c.Scheme = conf.Scheme
	cli, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return Discovery{r}
}

func (d *Discovery) GetEndPointHTTP(ctx context.Context, endpoint string) (string, error) {
	service, err := d.Discovery.GetService(ctx, endpoint)
	fmt.Print(service)
	if err != nil {
		return "", err
	}
	for _, item := range service {
		for _, endpoint := range item.Endpoints {
			if len(endpoint) >= 4 && endpoint[:4] == "http" {
				return endpoint, nil
			}
		}
	}
	return "", fmt.Errorf("endpoint不存在")
}

func (d *Discovery) GetEndPointGRPC(ctx context.Context, endpoint string) (string, error) {
	service, err := d.Discovery.GetService(ctx, endpoint)
	fmt.Print(service)
	if err != nil {
		return "", err
	}
	for _, item := range service {
		for _, endpoint := range item.Endpoints {
			if len(endpoint) >= 4 && endpoint[:4] == "grpc" {
				return endpoint, nil
			}
		}
	}
	return "", fmt.Errorf("endpoint不存在")
}
