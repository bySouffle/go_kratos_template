package cronjob

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
)

type Handler struct {
	Topic    string
	CallBack func(context.Context, *asynq.Task) error
}

func NewHandler(topic string, callback func(context.Context, *asynq.Task) error) Handler {
	return Handler{
		Topic:    topic,
		CallBack: callback,
	}
}

type Server struct {
	clientOpt asynq.RedisClientOpt
	Mux       *asynq.ServeMux
	Srv       *asynq.Server

	log *log.Helper
}

func (s *Server) Start(ctx context.Context) error {
	log.Infof("[cron] server use redis: %s db: %d", s.clientOpt.Addr, s.clientOpt.DB)

	if err := s.Srv.Run(s.Mux); err != nil {
		log.Errorf("[cron] server start fail: %v", err)
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[cron] server stopping")
	s.Srv.Shutdown()
	return nil
}

func NewCronServer(c *CronConfig, logger log.Logger) *Server {
	redisOpt := asynq.RedisClientOpt{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           int(c.Db),
		DialTimeout:  c.DialTimeout.AsDuration(),
		ReadTimeout:  c.ReadTimeout.AsDuration(),
		WriteTimeout: c.WriteTimeout.AsDuration(),
		PoolSize:     int(c.PoolSize),
	}

	cronSrv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: int(c.Concurrency),
			// Optionally specify multiple queues with different priority.
			//Queues: map[string]int{
			//	"critical": 6,
			//	"default":  3,
			//	"low":      1,
			//},
			// See the godoc for other configuration options
		},
	)

	return &Server{
		clientOpt: redisOpt,
		Mux:       asynq.NewServeMux(),
		Srv:       cronSrv,
		log:       log.NewHelper(logger),
	}
}

func RegisterCronHandler(s *Server, handler Handler) {
	s.Mux.HandleFunc(handler.Topic, handler.CallBack)
}
