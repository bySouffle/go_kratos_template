package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go_kratos_template/app/template/internal/biz"
)

type templateRepo struct {
	data *Data
	log  *log.Helper
}

func (t templateRepo) Save(ctx context.Context, template *biz.DataBase) (*biz.DataBase, error) {
	ctx, span := otel.Tracer("Data").Start(ctx, "templateRepoSave", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	span.SetAttributes(attribute.String("req", gconv.String(template)))
	return template, t.data.db.Create(template).Error
}

func (t templateRepo) Update(ctx context.Context, template *biz.DataBase) (*biz.DataBase, error) {
	//TODO implement me
	panic("implement me")
}

func (t templateRepo) FindByID(ctx context.Context, i int64) (*biz.DataBase, error) {
	//TODO implement me
	panic("implement me")
}

func (t templateRepo) ListByHello(ctx context.Context, s string) ([]*biz.DataBase, error) {
	//TODO implement me
	panic("implement me")
}

func (t templateRepo) ListAll(ctx context.Context) ([]*biz.DataBase, error) {
	//TODO implement me
	panic("implement me")
}

// NewTemplateRepo .
func NewTemplateRepo(data *Data, logger log.Logger) biz.TemplateRepo {
	return &templateRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
