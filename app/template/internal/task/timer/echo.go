package timer

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
	"go_kratos_template/app/template/internal/biz"
	"time"
)

type Timer struct {
	log *log.Helper
	t   time.Time
	uc  *biz.TemplateUseCase
}

func NewTimer(logger log.Logger, useCase *biz.TemplateUseCase) *Timer {
	mLog := log.NewHelper(log.With(logger, "Test", "app/template/internal/task/timer"))
	return &Timer{log: mLog, uc: useCase}
}

func (r *Timer) Echo(ctx context.Context, task *asynq.Task) error {
	now := time.Now()
	r.uc.WebSocketManager.SendMsgToGroup("/ws/conn", []byte(now.Sub(r.t).String()))
	r.t = now
	return nil
}
