package cronjob

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
	"time"
)

// PeriodicManager A scheduled task whose period can be modified
type PeriodicManager struct {
	clientOpt asynq.RedisClientOpt
	Manager   *asynq.PeriodicTaskManager

	log *log.Helper
}

func NewPeriodicManager(c *CronConfig, logger log.Logger, conf asynq.PeriodicTaskConfigProvider) *PeriodicManager {
	redisOpt := asynq.RedisClientOpt{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           int(c.Db),
		DialTimeout:  c.DialTimeout.AsDuration(),
		ReadTimeout:  c.ReadTimeout.AsDuration(),
		WriteTimeout: c.WriteTimeout.AsDuration(),
		PoolSize:     int(c.PoolSize),
	}

	manager, err := asynq.NewPeriodicTaskManager(
		asynq.PeriodicTaskManagerOpts{
			RedisConnOpt:               redisOpt,
			PeriodicTaskConfigProvider: conf,
			SyncInterval:               10 * time.Second,
		},
	)
	if err != nil {
		panic("[cronjob] NewPeriodicTaskManager failed")
	}

	return &PeriodicManager{
		clientOpt: redisOpt,
		Manager:   manager,
		log:       log.NewHelper(logger),
	}
}

func (m *PeriodicManager) Start(ctx context.Context) error {
	log.Infof("[cron] PeriodicManager use redis: %s db: %d", m.clientOpt.Addr, m.clientOpt.DB)
	if err := m.Manager.Start(); err != nil {
		log.Errorf("[cron] PeriodicManager start fail: %v", err)
		return err
	}
	return nil
}

func (m *PeriodicManager) Stop(ctx context.Context) error {
	log.Info("[cron] PeriodicManager stopping")
	m.Manager.Shutdown()
	return nil
}
