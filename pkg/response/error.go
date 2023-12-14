package response

import (
	"github.com/go-kratos/kratos/v2/encoding"
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/errors"
	"time"
)

func ErrorEncoder(w nethttp.ResponseWriter, r *nethttp.Request, err error) {
	e := errors.FromError(err)

	data := &Response{
		Code:       int(e.Code),
		Reason:     e.Reason,
		Message:    e.Message,
		TraceId:    r.Header.Get(TraceIdKey), // 从请求头中获取 traceId
		ServerTime: time.Now().Unix(),
	}
	codec := encoding.GetCodec("json")
	body, err := codec.Marshal(data)
	if err != nil {
		return
	}

	w.WriteHeader(200) // http code forever 200
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
