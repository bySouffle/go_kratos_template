package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	v1 "go_kratos_template/api/template/v1"
	"go_kratos_template/pkg/ws"
	"gorm.io/gorm"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Template is a Template model.
type Template struct {
	Name string
	ID   string
}

// TemplateRepo is a Greater repo.
type TemplateRepo interface {
	Save(context.Context, *DataBase) (*DataBase, error)
	Update(context.Context, *DataBase) (*DataBase, error)
	FindByID(context.Context, int64) (*DataBase, error)
	ListByHello(context.Context, string) ([]*DataBase, error)
	ListAll(context.Context) ([]*DataBase, error)
}

// TemplateUseCase is a Template useCase.
type TemplateUseCase struct {
	WebSocketManager *ws.ClientManager
	repo             TemplateRepo
	log              *log.Helper
}

// NewTemplateUseCase new a Template useCase.
func NewTemplateUseCase(repo TemplateRepo, logger log.Logger, ws *ws.ClientManager) *TemplateUseCase {
	return &TemplateUseCase{repo: repo, log: log.NewHelper(logger), WebSocketManager: ws}
}

// CreateTemplate creates a Template, and returns the new Template.
func (uc *TemplateUseCase) CreateTemplate(ctx context.Context, g *Template) (*Template, error) {
	ctx, span := otel.Tracer("Biz").Start(ctx, "CreateTemplateBiz", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	uc.log.WithContext(ctx).Infof("CreateTemplate: %v", g.Name)
	data := DataBase{
		Name: g.Name,
	}

	span.SetAttributes(attribute.String("g", gconv.String(data)))

	value, err := uc.repo.Save(ctx, &data)
	if err != nil {
		return &Template{}, err
	}
	return &Template{
		Name: value.Name,
	}, nil
}

func (uc *TemplateUseCase) GetWsMsg(ctx context.Context, g *Template) (*Template, error) {
	ch := uc.WebSocketManager.GetClientMsgChan(g.Name, g.ID)
	select {
	case data := <-ch:
		return &Template{Name: string(data)}, nil
	default:
	}
	return &Template{}, nil
}

type DataBase struct {
	gorm.Model
	Name string
}

func (tb *DataBase) TableName() string {
	return "template_data"
}
