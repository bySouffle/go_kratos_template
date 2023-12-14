// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/sdk/trace"
	"go_kratos_template/app/template/internal/biz"
	"go_kratos_template/app/template/internal/conf"
	"go_kratos_template/app/template/internal/data"
	"go_kratos_template/app/template/internal/server"
	"go_kratos_template/app/template/internal/service"
	"go_kratos_template/app/template/internal/task/timer"
	"go_kratos_template/pkg/ws"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(appInfo *conf.APPInfo, confServer *conf.Server, confData *conf.Data, logger log.Logger, tracerProvider *trace.TracerProvider, registry *conf.Registry, general *conf.General, experiment *conf.Experiment) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	templateRepo := data.NewTemplateRepo(dataData, logger)
	clientManager := ws.NewNewClientManagerWithRun(logger)
	templateUseCase := biz.NewTemplateUseCase(templateRepo, logger, clientManager)
	templateService := service.NewTemplateService(templateUseCase)
	grpcServer := server.NewGRPCServer(confServer, general, templateService, logger)
	httpServer := server.NewHTTPServer(confServer, general, experiment, templateService, logger, tracerProvider)
	timerTimer := timer.NewTimer(logger, templateUseCase)
	cronjobServer := server.NewCronServer(confServer, timerTimer, logger)
	register := server.NewCronRegister(confServer, logger)
	mqtt := server.NewMqttClient(confServer, logger)
	registrar := server.NewRegistrar(registry)
	app := newApp(appInfo, logger, grpcServer, httpServer, cronjobServer, register, mqtt, registrar)
	return app, func() {
		cleanup()
	}, nil
}
