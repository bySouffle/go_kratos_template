package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go_kratos_template/app/template/internal/conf"
)

const (
	TaskPeriodPrefix = "cronjob:task_period:schedule"
)

type PeriodicTaskConfigContainer map[string]string

type CronConfigProvider struct {
	redis *redis.Client
}

func NewCronConfigProvider(redis *redis.Client, cron *conf.Cron) *CronConfigProvider {
	var (
		ctx = context.Background()
	)
	b, err := redis.Get(ctx, TaskPeriodPrefix).Bytes()
	if err != nil || string(b) == "null" {
		initConf, err := json.Marshal(cron.TaskPeriod)
		if err != nil {
			panic("[server] NewCronConfigProvider init config error")
		}
		if err := redis.Set(ctx, TaskPeriodPrefix, initConf, 0).Err(); err != nil {
			panic("[server] NewCronConfigProvider Set config error")
		}
	}

	return &CronConfigProvider{
		redis: redis,
	}
}

func (p *CronConfigProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {

	var (
		c   PeriodicTaskConfigContainer
		ctx = context.Background()
	)

	b, err := p.redis.Get(ctx, TaskPeriodPrefix).Bytes()
	if err != nil || string(b) == "null" {
		fmt.Printf("获取任务周期失败: value[%v] err[%v]", b, err)
		return nil, err
	}

	if err := json.Unmarshal(b, &c); err != nil {
		fmt.Printf("解析任务周期失败:%v", err)
		return nil, err
	}

	var configs []*asynq.PeriodicTaskConfig
	for topic, cron := range c {
		configs = append(configs, &asynq.PeriodicTaskConfig{Cronspec: cron, Task: asynq.NewTask(topic, nil)})
	}
	return configs, nil
}

func (p *CronConfigProvider) AddCron(ctx context.Context, topic string, cron string) error {
	b, err := p.redis.Get(ctx, TaskPeriodPrefix).Bytes()
	if err != nil {
		return err
	}

	c := PeriodicTaskConfigContainer{}
	if err := json.Unmarshal(b, &c); err != nil {
		fmt.Printf("解析任务周期失败:%v", err)
		return err
	}

	c[topic] = cron

	confB, encErr := json.Marshal(c)
	if encErr != nil {
		return err
	}

	return p.redis.Set(ctx, TaskPeriodPrefix, confB, 0).Err()

}

func (p *CronConfigProvider) UpdateCron(ctx context.Context, topic string, cron string) error {
	b, err := p.redis.Get(ctx, TaskPeriodPrefix).Bytes()
	if err != nil {
		return err
	}

	c := PeriodicTaskConfigContainer{}
	if err := json.Unmarshal(b, &c); err != nil {
		fmt.Printf("解析任务周期失败:%v", err)
		return err
	}

	c[topic] = cron
	confB, encErr := json.Marshal(c)
	if encErr != nil {
		return err
	}

	return p.redis.Set(ctx, TaskPeriodPrefix, confB, 0).Err()

}

func (p *CronConfigProvider) DelCron(ctx context.Context, topic string) error {
	b, err := p.redis.Get(ctx, TaskPeriodPrefix).Bytes()
	if err != nil {
		return err
	}

	c := PeriodicTaskConfigContainer{}
	if err := json.Unmarshal(b, &c); err != nil {
		fmt.Printf("解析任务周期失败:%v", err)
		return err
	}
	delete(c, topic)
	confB, encErr := json.Marshal(c)
	if encErr != nil {
		return err
	}

	return p.redis.Set(ctx, TaskPeriodPrefix, confB, 0).Err()

}

func (p *CronConfigProvider) GetCron(ctx context.Context) (PeriodicTaskConfigContainer, error) {
	b, err := p.redis.Get(ctx, TaskPeriodPrefix).Bytes()
	if err != nil {
		return nil, err
	}
	c := PeriodicTaskConfigContainer{}
	return c, json.Unmarshal(b, &c)
}
