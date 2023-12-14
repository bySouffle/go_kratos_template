package cronjob

import (
	//"InspectionRobot/app/task_core_module/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
	"time"
)

type Job struct {
	Topic    string
	Schedule string
}

func NewJob(topic, schedule string) Job {
	return Job{
		Topic:    topic,
		Schedule: schedule,
	}
}

type Register struct {
	clientOpt asynq.RedisClientOpt
	Scheduler *asynq.Scheduler

	log *log.Helper
}

func (s *Register) Start(ctx context.Context) error {
	log.Infof("[cron] Scheduler use redis: %s db: %d", s.clientOpt.Addr, s.clientOpt.DB)

	if err := s.Scheduler.Run(); err != nil {
		log.Errorf("[cron] Scheduler start fail: %v", err)
		return err
	}

	return nil
}

func (s *Register) Stop(ctx context.Context) error {
	log.Info("[cron] Scheduler stopping")
	s.Scheduler.Shutdown()
	return nil
}

func NewRegisterServer(c *CronConfig, logger log.Logger) *Register {
	redisOpt := asynq.RedisClientOpt{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           int(c.Db),
		DialTimeout:  c.DialTimeout.AsDuration(),
		ReadTimeout:  c.ReadTimeout.AsDuration(),
		WriteTimeout: c.WriteTimeout.AsDuration(),
		PoolSize:     int(c.PoolSize),
	}

	return &Register{
		clientOpt: redisOpt,
		Scheduler: asynq.NewScheduler(
			redisOpt,
			&asynq.SchedulerOpts{Location: time.Local},
		),
		log: log.NewHelper(logger),
	}
}

func RegisterCronJob(s *Register, job Job) {
	task := asynq.NewTask(job.Topic, nil)
	_, err := s.Scheduler.Register(job.Schedule, task)
	if err != nil {
		panic(err)
	}
}
