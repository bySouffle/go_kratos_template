package boot

import (
	"github.com/go-kratos/kratos/v2/config"
	"go_kratos_template/app/template/internal/conf"
)

type CronOption func(*conf.Cron)

type Cron struct {
	C     config.Config
	Param conf.Cron
}

func NewCron(c config.Config) *Cron {
	return &Cron{C: c}
}

func (d *Cron) Load() *conf.Cron {
	if err := d.C.Scan(&d.Param); err != nil {
		panic(err)
	}
	return &d.Param
}

func (d *Cron) Setting(opts ...CronOption) *Cron {
	for _, opt := range opts {
		opt(&d.Param)
	}
	return d
}
