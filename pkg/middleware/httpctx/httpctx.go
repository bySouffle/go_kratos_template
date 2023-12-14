package httpctx

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	"net"
)

type HttpC struct{}

func FromContext(ctx context.Context) (*http.Request, bool) {
	h, ok := ctx.Value("ctx").(*http.Request)
	return h, ok
}

// GetHttpReqContext	获取原始请求的ctx
func GetHttpReqContext() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// Do something on entering
			//ctx = context.WithValue(ctx, HttpC{}, ctx)
			if httpCtx, ok := http.RequestFromServerContext(ctx); ok {
				host, _, err := net.SplitHostPort(httpCtx.RemoteAddr)
				if err != nil {
					return nil, err
				}
				ctx = context.WithValue(ctx, "ip", host)
				ctx = context.WithValue(ctx, "ctx", httpCtx)

			}
			defer func() {
				// Do something on exiting
			}()
			return handler(ctx, req)
		}
	}
}
