package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go_kratos_template/pkg/response"
	"go_kratos_template/pkg/ws"
	"io"
	"runtime/debug"
)

func (t *TemplateService) KratosWSHandler(w http.ResponseWriter, c *http.Request) {
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
			_, err := w.Write(gconv.Bytes(errStr))
			if err != nil {
				return
			}
		}
		span.End()
	}()

	client := ws.NewClient(ws.WithURL("/ws/conn"), ws.WithUUID())
	err := client.UpGrader(w, c)
	if err != nil {
		log.Errorf("UpGrader failed: %v", err)
		return
	}
	if err := client.Run(); err != nil {
		client.CloseChan <- struct{}{}
		log.Errorf("%v", err)
		return
	}
	regErr := t.uc.WebSocketManager.Register(client)
	if regErr != nil {
		client.CloseChan <- struct{}{}
		return
	}
	client.SendMsg <- []byte("zzzzzz")
}

func (t *TemplateService) KratosWSClose(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("Service").Start(r.Context(), "TemplateService", trace.WithSpanKind(trace.SpanKindProducer))
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	var reqData map[string]interface{}
	if err := json.Unmarshal(body, &reqData); err != nil {
		return
	}
	client := ws.Client{}

	if _, ok := reqData["close"].(string); ok {
		t.uc.WebSocketManager.Stop()
		w.Write(body)
		return
	}

	if url, ok := reqData["url"].(string); ok {
		client.URL = url
	}
	if uuid, ok := reqData["UUID"].(string); ok {
		client.ID = uuid
	}
	err = t.uc.WebSocketManager.LogOut(&client)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%v", err)))
	}
	w.Write(body)
}
