package response

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/encoding"
	"google.golang.org/protobuf/proto"
	nethttp "net/http"
	"time"
)

const TraceIdKey = "TraceId"

type Response struct {
	Code       int         `json:"code"`
	Reason     string      `json:"reason"`
	Message    interface{} `json:"message"`
	TraceId    string      `json:"traceId"`
	ServerTime int64       `json:"serverTime"`
}

// RespEncoder
// @see https://go-kratos.dev/docs/component/encoding
func RespEncoder(w nethttp.ResponseWriter, r *nethttp.Request, i interface{}) error {
	codec := encoding.GetCodec("json")
	messageMap := make(map[string]interface{})
	messageStr, _ := codec.Marshal(i.(proto.Message))
	_ = codec.Unmarshal(messageStr, &messageMap)

	resp := Response{
		Code:       200,
		Reason:     "",
		TraceId:    r.Header.Get(TraceIdKey), // 从请求头中获取 traceId
		ServerTime: time.Now().Unix(),
	}

	message, err := codec.Marshal(i)
	_ = json.Unmarshal(message, &resp.Message)
	if err != nil {
		return err
	}

	data, err := codec.Marshal(resp)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return nil
}
