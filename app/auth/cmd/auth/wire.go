//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"go_kratos_template/app/auth/internal/biz"
	"go_kratos_template/app/auth/internal/conf"
	"go_kratos_template/app/auth/internal/data"
	"go_kratos_template/app/auth/internal/server"
	"go_kratos_template/app/auth/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.APPInfo, *conf.Server, *conf.Data, log.Logger, *tracesdk.TracerProvider, *conf.Registry, *conf.General, *conf.Experiment, *conf.Security) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
