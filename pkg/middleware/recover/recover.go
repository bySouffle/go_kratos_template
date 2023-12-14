package recover

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"runtime"
)

// RecoverMiddleware 创建一个用于处理panic的中间件 可以讲错误栈存入trace中方便排查问题
func RecoverMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// Notice span
			ctx2, span := otel.Tracer("Middleware").Start(ctx, "RecoverMiddleware", trace.WithSpanKind(trace.SpanKindInternal))
			defer func() {
				if rerr := recover(); rerr != nil {
					// 构建错误日志
					buf := make([]byte, 64<<10)
					n := runtime.Stack(buf, false)
					// Notice log...
					log.Error(ctx2, "自定义中间件捕获到了异常!", rerr, " stack: ", string(buf[:n]))
					// Notice 可以把错误栈写到trace中！！！
					span.SetAttributes(attribute.String("InternalServerError", gconv.String(buf)))
					// Notice 设置错误信息,返回错误的响应！
					err = errors.New("自定义中间件捕获-Internal Server Error")
				}
				span.End()
			}()
			return handler(ctx, req)
		}
	}
}
