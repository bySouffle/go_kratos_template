package boot

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go_kratos_template/app/template/internal/conf"
	otlp "go_kratos_template/pkg/trace/otlptrace"
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

func (t *Trace) Run() *sdktrace.TracerProvider {
	var tp *sdktrace.TracerProvider

	resourceRet, errRet := t.resource()
	if errRet != nil {
		log.Errorw("trace_resource_error", errRet)
		panic(errRet)
	}
	if t.Param.Experiment.EnableTrace {
		//	OTLP
		if t.Param.Experiment.Trace.Exporter == otlp.OTLPHTTP ||
			t.Param.Experiment.Trace.Exporter == otlp.OTLPGRPC {
			c := &otlp.Config{}
			err := gconv.Structs(t.Param.Experiment.Trace, c)
			if err != nil {
				panic("[trace] OTLPGRPC config conv error")
			}
			tp, _ := c.InitOpenTelemetry()
			return tp
		}
		//	jeager file stdout
		exporter, err := t.exporter(t.Param.Experiment.Trace.Exporter) // Notice 配置的 jaeger
		if err != nil {
			log.Errorw("trace_exporter_error", err)
			panic(err)
		}
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resourceRet),
		)
	} else {
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithResource(resourceRet),
		)
	}
	otel.SetTracerProvider(tp)
	return tp
}

func (t *Trace) exporter(types string) (sdktrace.SpanExporter, error) {
	var exporterRet sdktrace.SpanExporter
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
			semconv.ServiceName(t.Param.App.Name),                         //实例名称
			attribute.String("ID", t.Param.App.ID),                        // ID
			attribute.String("environment", t.Param.App.Environment),      // 相关环境
			attribute.String("Version", t.Param.App.Version),              //版本
			attribute.String("token", t.Param.Experiment.Trace.GrpcToken), //token
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
