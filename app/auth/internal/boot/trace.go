package boot

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go_kratos_template/app/auth/internal/conf"
	"os"
)

type TraceOption func(*conf.Bootstrap)

type Trace struct {
	C     config.Config
	Param conf.Bootstrap
}

func NewBootTrace(c config.Config) *Trace {
	return &Trace{
		C: c,
	}
}

func (t *Trace) Load() *Trace {
	if err := t.C.Scan(&t.Param); err != nil {
		panic(err)
	}
	return t
}

func (t *Trace) Run() *tracesdk.TracerProvider {
	var tp *tracesdk.TracerProvider

	resourceRet, errRet := t.resource()
	if errRet != nil {
		log.Errorw("trace_resource_error", errRet)
		panic(errRet)
	}
	if t.Param.Experiment.EnableTrace {
		exporter, err := t.exporter(t.Param.Experiment.Trace.Exporter) // Notice 配置的 jaeger
		if err != nil {
			log.Errorw("trace_exporter_error", err)
			panic(err)
		}
		tp = tracesdk.NewTracerProvider(
			tracesdk.WithBatcher(exporter),
			tracesdk.WithResource(resourceRet),
		)
	} else {
		tp = tracesdk.NewTracerProvider(
			tracesdk.WithResource(resourceRet),
		)
	}
	otel.SetTracerProvider(tp)
	return tp
}

func (t *Trace) exporter(types string) (tracesdk.SpanExporter, error) {
	var exporterRet tracesdk.SpanExporter
	var err error

	switch types {
	case "stdout":
		exporterRet, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint())
	case "file":
		var osFile *os.File
		osFile, err = os.Create(t.Param.Experiment.Trace.TraceFilePath)
		exporterRet, err = stdouttrace.New(
			stdouttrace.WithWriter(osFile))
	case "jaeger":
		exporterRet, err = jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(t.Param.Experiment.Trace.CollectorEndpoint),
			))
	}
	return exporterRet, err
}

func (t *Trace) resource() (*resource.Resource, error) {
	resourceRet, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(t.Param.App.Name),                     //实例名称
			attribute.String("ID", t.Param.App.ID),                    // ID
			attribute.String("environment", t.Param.App.Environment),  // 相关环境
			attribute.String("Version", t.Param.App.Version),          //版本
			attribute.String("token", t.Param.Experiment.Trace.Token), //token
		),
	)
	if err != nil {
		log.Infof("[Trace] errResource: %v", err)
	}
	return resourceRet, err
}

func (t *Trace) Setting(opts ...TraceOption) *Trace {
	for _, opt := range opts {
		opt(&t.Param)
	}
	return t
}
