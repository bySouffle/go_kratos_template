package boot

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"go_kratos_template/app/auth/internal/conf"
	plog "go_kratos_template/pkg/log"
)

type LogOption func(*conf.Logger)

type Log struct {
	C     config.Config
	Param conf.Logger
}

func NewBootLog(c config.Config) *Log {
	return &Log{C: c}
}

func (l *Log) Load() *Log {
	if err := l.C.Scan(&l.Param); err != nil {
		panic(err)
	}
	return l
}

func (l *Log) Setting(opts ...LogOption) *Log {
	for _, opt := range opts {
		opt(&l.Param)
	}
	return l
}

func (l *Log) Run() log.Logger {
	lConf := plog.Config{
		Development:       l.Param.Development,
		DisableCaller:     l.Param.DisableCaller,
		DisableStacktrace: l.Param.DisableStacktrace,
		Encoding:          l.Param.Encoding,
		Level:             l.Param.Level,
		Name:              l.Param.Name,
		Writers:           l.Param.Writers,
		LoggerDir:         l.Param.LoggerDir,
		LogRollingPolicy:  l.Param.LogRollingPolicy,
		LogBackupCount:    uint(l.Param.LogBackupCount),
	}

	logger := plog.Init(&lConf)
	return logger
}
