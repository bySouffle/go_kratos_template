package main

import (
	"github.com/go-kratos/kratos/v2/registry"
	"go_kratos_template/app/template/internal/boot"
	"go_kratos_template/app/template/internal/conf"
	"go_kratos_template/pkg/client"
	"go_kratos_template/pkg/cronjob"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func newApp(info *conf.APPInfo, logger log.Logger, gs *grpc.Server, hs *http.Server, cs *cronjob.Server, rs *cronjob.Register, ms *client.Mqtt, cr registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(info.ID),
		kratos.Name(info.Name),
		kratos.Version(info.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
			cs,
			rs,
			ms,
		),
		kratos.Registrar(cr),
	)
}

func main() {
	c := boot.NewBootConf()
	c.Load()
	defer c.Close()
	bt := boot.NewBootStrap(c.C).Load().Setting(boot.WithUUID())
	logger := boot.NewBootLog(c.C).Load().Run()
	tp := boot.NewBootTrace(c.C).Load().Run()
	rc := boot.NewRegistry(c.C).Load()
	app, cleanup, err := wireApp(bt.Param.App, bt.Param.Server, bt.Param.Data, logger, tp, rc, bt.Param.General, bt.Param.Experiment, bt.Param.Security)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
