package service

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go_kratos_template/pkg/registry/consul"
	"go_kratos_template/pkg/response"
	"runtime/debug"
)

// 单独注册的路由
func (t *TemplateService) KratosSingleHandler(w http.ResponseWriter, c *http.Request) {
	ctx, span := otel.Tracer("Service").Start(c.Context(), "TemplateService", trace.WithSpanKind(trace.SpanKindProducer))
	// Notice 需要单独写recover panic
	defer func() {
		if err := recover(); err != nil {
			stackInfo := string(debug.Stack())
			span.SetAttributes(attribute.String("errStack", stackInfo))
			errStr, _ := json.Marshal(response.Response{
				Code:       500,
				Reason:     "single-handler Internal Server Error!",
				Message:    nil,
				ServerTime: 0,
				TraceId:    otelTrace.SpanContextFromContext(ctx).TraceID().String(),
			})
			w.Write(gconv.Bytes(errStr))
		}
		span.End()
	}()
	cr := &consul.RegistryConfig{
		Scheme:    "http",
		Address:   "127.0.0.1:8500",
		Discovery: "template",
	}
	ds := consul.NewDiscovery(cr)
	pointHTTP, err := ds.GetEndPointHTTP(ctx, cr.Discovery)
	if err != nil {
		panic("")
	}
	pointGRPC, err := ds.GetEndPointGRPC(ctx, cr.Discovery)
	if err != nil {
		panic("")
	}
	//panic("single-handler")

	retStr, _ := json.Marshal(response.Response{
		Code:   200,
		Reason: "单独注册的路由!",
		Message: map[string]string{
			"HTTP": pointHTTP,
			"GRPC": pointGRPC,
		},
		ServerTime: 0,
		TraceId:    otelTrace.SpanContextFromContext(ctx).TraceID().String(),
	})
	w.Write(gconv.Bytes(retStr))
}
