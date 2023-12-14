package server

import (
	"go_kratos_template/app/template/internal/conf"
	"go_kratos_template/pkg/registry/consul"
)

func NewDiscovery(registry *conf.Registry) consul.Discovery {
	r := consul.NewDiscovery(DiscoveryConfigConv(registry))
	return r
}

func DiscoveryConfigConv(registry *conf.Registry) *consul.RegistryConfig {
	return &consul.RegistryConfig{
		Address:   registry.Endpoint.Address,
		Scheme:    registry.Endpoint.Scheme,
		Discovery: registry.Endpoint.Discovery,
	}
}
