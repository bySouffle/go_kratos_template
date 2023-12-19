package data

import (
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go_kratos_template/app/auth/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	slog "log"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewAuthRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db    *gorm.DB
	redis *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: NewDB(c), redis: NewRedis(c)}, cleanup, nil
}

func NewDB(c *conf.Data) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)
	if c.Database.Driver == "mysql" {
		db, err = gorm.Open(mysql.Open(c.Database.Source), gormConfig(c))

	} else if c.Database.Driver == "sqlite" {
		db, err = gorm.Open(sqlite.Open(c.Database.Source), gormConfig(c))
	}
	sql, sqlErr := db.DB()
	if sqlErr != nil {
		panic("failed to connect database")
	}
	sql.SetMaxOpenConns(int(c.Database.MaxOpenConn))
	sql.SetMaxIdleConns(int(c.Database.MaxIdleConn))
	sql.SetConnMaxLifetime(c.Database.ConnMaxLifeTime.AsDuration())
	// Initialize otel plugin with options
	plugin := tracing.NewPlugin(
	// include any options here
	)

	// set trace
	err = db.Use(plugin)
	if err != nil {
		panic("failed to Use database tracing plugin")
	}

	if err != nil {
		panic("failed to connect database")
	}

	if err != nil {
		log.Errorf("failed opening connection to %v: %v", c.Database.Source, err)
		panic("failed to connect database")
	}

	return db
}

func NewRedis(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})
	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}

	return rdb
}

func gormConfig(c *conf.Data) *gorm.Config {
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true} // 禁止外键约束, 生产环境不建议使用外键约束
	// 打印所有SQL
	if c.Database.ShowLog {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}
	// 只打印慢查询
	if c.Database.SlowThreshold.AsDuration().Nanoseconds() > 0 {
		config.Logger = logger.New(
			//将标准输出作为Writer
			slog.New(os.Stdout, "\r\n", slog.LstdFlags),
			logger.Config{
				//设定慢查询时间阈值
				SlowThreshold: c.Database.SlowThreshold.AsDuration(), // nolint: golint
				Colorful:      true,
				//设置日志级别，只有指定级别以上会输出慢查询日志
				LogLevel: logger.Warn,
			},
		)
	}

	return config
}
