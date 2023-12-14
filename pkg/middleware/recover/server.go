package recover

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go_kratos_template/pkg/response"
	"time"
)

// Server is an server logging kratos_middleware.
func Server(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code       int32
				reason     string
				kind       string
				requestUrl string
				remoteAddr string
				operation  string
				header     map[string][]string
			)

			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				ht, _ := info.(*http.Transport)
				request := ht.Request()
				header = request.Header
				requestUrl = request.RequestURI
				remoteAddr = request.RemoteAddr
				kind = info.Kind().String()
				operation = info.Operation()
				// Notice 添加 trace_id 到 请求的 Header 中！！！
				traceID := otelTrace.SpanContextFromContext(ctx)
				//fmt.Println("Server中间件生成的traceId: ", traceID)
				request.Header.Set(response.TraceIdKey, traceID.TraceID().String())
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}

			//	输出info之外的信息
			level, stack := extractError(err)
			if level != log.LevelInfo {
				_ = log.WithContext(ctx, logger).Log(level,
					"kind", "server",
					"header", header,
					"request_url", requestUrl,
					"remote_addr", remoteAddr,
					"component", kind,
					"operation", operation,
					"args", extractArgs(req),
					"code", code,
					"reason", reason,
					"stack", stack,
					"latency", time.Since(startTime).Seconds(),
				)
			}

			return
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}
