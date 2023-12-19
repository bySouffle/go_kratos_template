package boot

import (
	"github.com/go-kratos/kratos/v2/config"
	"go_kratos_template/app/auth/internal/conf"
)

type RegistryOption func(*conf.Registry)

type Registry struct {
	C     config.Config
	Param conf.Registry
}

func NewRegistry(c config.Config) *Registry {
	return &Registry{C: c}
}

func (d *Registry) Load() *conf.Registry {
	if err := d.C.Scan(&d.Param); err != nil {
		panic(err)
	}
	return &d.Param
}

func (d *Registry) Setting(opts ...RegistryOption) *Registry {
	for _, opt := range opts {
		opt(&d.Param)
	}
	return d
}
