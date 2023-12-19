package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/util/gconv"
	"go_kratos_template/app/auth/internal/conf"
	auth "go_kratos_template/pkg/middleware/auth"
)

type Auth struct {
	Name  string
	ID    string
	Token string
}

type AuthRepo interface {
	Create(context.Context, *Auth) error
	Delete(context.Context, *Auth) error
	Update(context.Context, *Auth) error
	FindByAuth(context.Context, *Auth) (*Auth, error)
	ListByName(context.Context, []*Auth) ([]*Auth, error)
}

// AuthUseCase is a Auth useCase.
type AuthUseCase struct {
	repo AuthRepo
	log  *log.Helper
	j    *auth.JWT
}

func NewTAuthUseCase(repo AuthRepo, logger log.Logger, jwt *auth.JWT) *AuthUseCase {
	return &AuthUseCase{repo: repo, log: log.NewHelper(logger), j: jwt}
}
func (uc *AuthUseCase) GeneralToken(ctx context.Context, g *Auth) (*Auth, error) {
	g.Token = uc.j.GenerateToken(gconv.Int64(g.ID), g.Name)
	err := uc.repo.Create(ctx, g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (uc *AuthUseCase) UpdateToken(ctx context.Context, g *Auth) (*Auth, error) {
	g.Token = uc.j.GenerateToken(gconv.Int64(g.ID), g.Name)
	err := uc.repo.Update(ctx, g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (uc *AuthUseCase) DeleteToken(ctx context.Context, g *Auth) error {
	err := uc.repo.Delete(ctx, g)
	if err != nil {
		return err
	}
	return nil
}

func (uc *AuthUseCase) QueryToken(ctx context.Context, g *Auth) (*Auth, error) {
	result, err := uc.repo.FindByAuth(ctx, g)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (uc *AuthUseCase) QueryTokenList(ctx context.Context, g []*Auth) ([]*Auth, error) {
	result, err := uc.repo.ListByName(ctx, g)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (uc *AuthUseCase) CheckToken(ctx context.Context, g *Auth) (interface{}, error) {
	token, err := uc.j.ParseToken(g.Token)
	if err != nil {
		return nil, err
	}
	return token, err
}

func NewJwt(s *conf.Security) *auth.JWT {
	return auth.NewJwt(s.CookieName, s.JwtSecret, s.JwtTimeout.AsDuration())
}
