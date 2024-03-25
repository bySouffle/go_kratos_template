package server

import (
	"github.com/gogf/gf/v2/util/gconv"
	"go_kratos_template/app/template/internal/conf"
	"go_kratos_template/app/template/internal/task/timer"

	"github.com/go-kratos/kratos/v2/log"
	"go_kratos_template/pkg/cronjob"
)

func NewCronServer(c *conf.Server, t *timer.Timer, logger log.Logger) *cronjob.Server {
	cc := &cronjob.CronConfig{}
	err := gconv.Struct(c.Cron, cc)
	if err != nil {
		panic("[NewCronServer] 配置转换失败")
	}

	cron := cronjob.NewCronServer(cc, logger)
	cronjob.RegisterCronHandler(cron, cronjob.NewHandler("timer:echo", t.Echo))

	return cron
}

func NewCronRegister(c *conf.Server, logger log.Logger) *cronjob.Register {
	cc := &cronjob.CronConfig{}
	err := gconv.Struct(c.Cron, cc)
	if err != nil {
		panic("[NewCronRegister] 配置转换失败")
	}

	register := cronjob.NewRegisterServer(cc, logger)
	cronjob.RegisterCronJob(register, cronjob.NewJob("timer:echo", "@every 30s"))
	return register
}

func NewCronManager(c *conf.Server, logger log.Logger, source *CronConfigProvider) *cronjob.PeriodicManager {
	config := cronjob.CronConfig{}
	if err := gconv.Struct(c.Cron, &config); err != nil {
		return nil
	}
	return cronjob.NewPeriodicManager(&config, logger, source)
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
