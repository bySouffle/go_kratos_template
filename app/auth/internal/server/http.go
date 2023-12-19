package server

import (
	"crypto/tls"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/gorilla/handlers"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	v1 "go_kratos_template/api/auth/v1"
	"go_kratos_template/app/auth/internal/conf"
	"go_kratos_template/app/auth/internal/service"
	"go_kratos_template/pkg/middleware/httpctx"
	"go_kratos_template/pkg/middleware/recover"
	"go_kratos_template/pkg/response"
	"net/http/pprof"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, g *conf.General, experiment *conf.Experiment, authSrv *service.AuthService, logger log.Logger, tp *tracesdk.TracerProvider) *http.Server {
	var opts = []http.ServerOption{
		http.ResponseEncoder(response.RespEncoder), //	success resp: 有读取请求头中的traceId的操作
		http.ErrorEncoder(response.ErrorEncoder),   //	err resp: 有读取请求头中的traceId的操作
		http.Middleware(
			//recover.Recovery(),
			//metrics.Server(
			//	metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
			//	metrics.WithRequests(prom.NewCounter(_metricRequests)),
			//),
			tracing.Server(tracing.WithTracerProvider(tp)), //	Notice 必须把它写到 recover.Server(logger) 前面 才能生成正常的traceId
			recover.Server(logger),                         //	Notice 自己写的中间件: 请求头加上traceId、在日志中打印请求的信息...
			recover.RecoverMiddleware(),                    //	Notice 自定义的Recover中间件: 优点是可以将错误栈信息加到trace中方便排查问题
			//logging.Server(logger),	//	请求信息输出
			httpctx.GetHttpReqContext(),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
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
		opts = append(opts, http.TLSConfig(tlsConf))
	}

	srv := http.NewServer(opts...)
	h := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", h) ///q/swagger-ui/
	v1.RegisterAuthHTTPServer(srv, authSrv)

	if experiment.EnablePprof {
		RegisterPprof(srv)
	}

	return srv
}
func RegisterPprof(s *http.Server) {
	s.HandleFunc("/debug/pprof", pprof.Index)
	s.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.HandleFunc("/debug/pprof/trace", pprof.Trace)
	s.HandleFunc("/debug/allocs", pprof.Handler("allocs").ServeHTTP)
	s.HandleFunc("/debug/block", pprof.Handler("block").ServeHTTP)
	s.HandleFunc("/debug/goroutine", pprof.Handler("goroutine").ServeHTTP)
	s.HandleFunc("/debug/heap", pprof.Handler("heap").ServeHTTP)
	s.HandleFunc("/debug/mutex", pprof.Handler("mutex").ServeHTTP)
	s.HandleFunc("/debug/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}
