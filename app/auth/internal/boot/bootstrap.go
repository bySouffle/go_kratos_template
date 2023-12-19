package boot

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/google/uuid"
	"go_kratos_template/app/auth/internal/conf"
)

type BootStrapOption func(*conf.Bootstrap)

type BootStrap struct {
	C     config.Config
	Param conf.Bootstrap
}

func NewBootStrap(c config.Config) *BootStrap {
	return &BootStrap{C: c}
}

func (b *BootStrap) Load() *BootStrap {
	if err := b.C.Scan(&b.Param); err != nil {
		panic(err)
	}

	b.Param.App.ID = uuid.New().String()

	return b
}

func (b *BootStrap) Setting(opts ...BootStrapOption) *BootStrap {
	for _, opt := range opts {
		opt(&b.Param)
	}
	return b
}

func WithID(id string) BootStrapOption {
	return func(bootstrap *conf.Bootstrap) {
		bootstrap.App.ID = id
	}
}

func WithUUID() BootStrapOption {
	return func(bootstrap *conf.Bootstrap) {
		bootstrap.App.ID = uuid.New().String()
	}
}
