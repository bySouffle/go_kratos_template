package server

import (
	"go_kratos_template/app/template/internal/conf"
	"go_kratos_template/app/template/internal/task/timer"

	"github.com/go-kratos/kratos/v2/log"
	"go_kratos_template/pkg/cronjob"
)

func NewCronServer(c *conf.Server, t *timer.Timer, logger log.Logger) *cronjob.Server {
	cron := cronjob.NewCronServer(ConfigConv(c.Cron), logger)
	cronjob.RegisterCronHandler(cron, cronjob.NewHandler("timer:echo", t.Echo))

	return cron
}

func NewCronRegister(c *conf.Server, logger log.Logger) *cronjob.Register {
	register := cronjob.NewRegisterServer(ConfigConv(c.Cron), logger)
	cronjob.RegisterCronJob(register, cronjob.NewJob("timer:echo", "@every 30s"))
	return register
}

func ConfigConv(c *conf.Server_Cron) *cronjob.CronConfig {
	return &cronjob.CronConfig{
		Network:      c.Network,
		Addr:         c.Addr,
		Db:           c.Db,
		Password:     c.Password,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		MinIdleConn:  c.MinIdleConn,
		PoolSize:     c.PoolSize,
		PoolTimeout:  c.PoolTimeout,
		Concurrency:  c.Concurrency,
	}
}
