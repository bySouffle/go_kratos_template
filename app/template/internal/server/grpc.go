package server

import (
	"crypto/tls"
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	template "go_kratos_template/api/template/v1"
	"go_kratos_template/app/template/internal/conf"
	"go_kratos_template/app/template/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, g *conf.General, templateSrv *service.TemplateService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prom.NewCounter(_metricRequests)),
			),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	if g.SSL == true {
		cert, err := tls.LoadX509KeyPair(g.Cert, g.Key)
		if err != nil {
			panic(err)
		}
		tlsConf := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		}
		opts = append(opts, grpc.TLSConfig(tlsConf))
	}

	srv := grpc.NewServer(opts...)
	template.RegisterTemplateServer(srv, templateSrv)
	return srv
}
