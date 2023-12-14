package service

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	pb "go_kratos_template/api/template/v1"
	"go_kratos_template/app/template/internal/biz"
)

type TemplateService struct {
	pb.UnimplementedTemplateServer
	uc *biz.TemplateUseCase
}

func NewTemplateService(uc *biz.TemplateUseCase) *TemplateService {
	return &TemplateService{
		uc: uc,
	}
}

func (t *TemplateService) CreateTemplate(ctx context.Context, req *pb.CreateTemplateRequest) (*pb.CreateTemplateReply, error) {
	ctx, span := otel.Tracer("Service").Start(ctx, "CreateTemplate", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	span.SetAttributes(attribute.String("req", gconv.String(req)))

	data := biz.Template{Name: req.Name}
	template, err := t.uc.CreateTemplate(ctx, &data)
	if err != nil {
		return nil, err
	}
	return &pb.CreateTemplateReply{Name: template.Name}, nil
}
func (t *TemplateService) UpdateTemplate(ctx context.Context, req *pb.UpdateTemplateRequest) (*pb.UpdateTemplateReply, error) {
	return &pb.UpdateTemplateReply{}, nil
}
func (t *TemplateService) DeleteTemplate(ctx context.Context, req *pb.DeleteTemplateRequest) (*pb.DeleteTemplateReply, error) {
	return &pb.DeleteTemplateReply{}, nil
}
func (t *TemplateService) GetTemplate(ctx context.Context, req *pb.GetTemplateRequest) (*pb.GetTemplateReply, error) {
	msg, err := t.uc.GetWsMsg(ctx, &biz.Template{Name: req.Name, ID: req.Id})
	if err != nil {
		return nil, err
	}

	return &pb.GetTemplateReply{Name: msg.Name}, nil
}
func (t *TemplateService) ListTemplate(ctx context.Context, req *pb.ListTemplateRequest) (*pb.ListTemplateReply, error) {
	return &pb.ListTemplateReply{}, nil
}
func (t *TemplateService) WSTemplate(ctx context.Context, req *pb.WSTemplateRequest) (*pb.WSTemplateReply, error) {
	//	TODO 中间件中无法获取 http.ResponseWriter,在这里无法转换为websocket transport/http/context.go+93
	//reqCtx, ok := httpctx.FromContext(ctx)
	//if !ok {
	//	return nil, errors.BadRequest("custom_error", "not a http.Context")
	//}
	//reqs := reqCtx
	//
	//upGrader := websocket.Upgrader{}
	//c, _ := upGrader.Upgrade(nil, reqs, nil)
	//err := c.WriteMessage(websocket.TextMessage, []byte("websocket connect"))
	//if err != nil {
	//	return nil, err
	//}
	//return nil, nil
	if ip, ok := ctx.Value("ip").(string); ok {
		return &pb.WSTemplateReply{Name: ip}, nil
	}
	return &pb.WSTemplateReply{}, nil
}
