package boot

import (
	"github.com/go-kratos/kratos/v2/config"
	"go_kratos_template/app/auth/internal/conf"
)

type DeviceOption func(*conf.Device)

type Device struct {
	C     config.Config
	Param conf.Device
}

func NewDevice(c config.Config) *Device {
	return &Device{C: c}
}

func (d *Device) Load() {
	if err := d.C.Scan(&d.Param); err != nil {
		panic(err)
	}
}

func (d *Device) Setting(opts ...DeviceOption) *Device {
	for _, opt := range opts {
		opt(&d.Param)
	}
	return d
}
