package response

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/encoding"
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/errors"
	"time"
)

type ErrorResponse struct {
	Code       int             `json:"code"`
	Reason     string          `json:"reason"`
	Message    interface{}     `json:"message"`
	TraceId    string          `json:"traceId"`
	ServerTime int64           `json:"serverTime"`
	MetaData   json.RawMessage `json:"metaData,omitempty"`
}

func ErrorEncoder(w nethttp.ResponseWriter, r *nethttp.Request, err error) {
	e := errors.FromError(err)

	data := &ErrorResponse{
		Code:       int(e.Code),
		Reason:     e.Reason,
		Message:    e.Message,
		TraceId:    r.Header.Get(TraceIdKey), // 从请求头中获取 traceId
		ServerTime: time.Now().Unix(),
	}

	if len(e.GetMetadata()) > 0 {
		metaData, err := json.Marshal(e.GetMetadata())
		if err == nil {
			data.MetaData = metaData
		}
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
