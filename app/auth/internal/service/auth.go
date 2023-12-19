package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"go_kratos_template/app/auth/internal/biz"

	pb "go_kratos_template/api/auth/v1"
)

type AuthService struct {
	pb.UnimplementedAuthServer
	uc *biz.AuthUseCase
}

func NewAuthService(uc *biz.AuthUseCase) *AuthService {
	return &AuthService{uc: uc}
}

func (s *AuthService) CreateAuth(ctx context.Context, req *pb.CreateAuthRequest) (*pb.CreateAuthReply, error) {
	token, err := s.uc.GeneralToken(ctx, &biz.Auth{
		Name: req.Name,
		ID:   req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateAuthReply{
		Token: token.Token,
	}, nil
}
func (s *AuthService) UpdateAuth(ctx context.Context, req *pb.UpdateAuthRequest) (*pb.UpdateAuthReply, error) {
	token, err := s.uc.UpdateToken(ctx, &biz.Auth{
		Name: req.Name,
		ID:   req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateAuthReply{
		Token: token.Token,
	}, nil
}
func (s *AuthService) DeleteAuth(ctx context.Context, req *pb.DeleteAuthRequest) (*pb.DeleteAuthReply, error) {
	err := s.uc.DeleteToken(ctx, &biz.Auth{
		Name: req.Name,
		ID:   req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteAuthReply{
		Ok: "1",
	}, nil
}
func (s *AuthService) GetAuth(ctx context.Context, req *pb.GetAuthRequest) (*pb.GetAuthReply, error) {
	token, err := s.uc.QueryToken(ctx, &biz.Auth{
		Name: req.Name,
		ID:   req.Id,
	})
	if err != nil {
		return &pb.GetAuthReply{}, err
	}

	return &pb.GetAuthReply{
		Token: token.Token,
	}, nil
}
func (s *AuthService) ListAuth(ctx context.Context, req *pb.ListAuthRequest) (*pb.ListAuthReply, error) {
	var auths []*biz.Auth

	if err := gconv.Structs(req.User, auths); err != nil {
		return nil, err
	}
	g.Dump(auths)

	token, err := s.uc.QueryTokenList(ctx, auths)
	if err != nil {
		return nil, err
	}

	result := &pb.ListAuthReply{}

	if err := gconv.Structs(token, result.Token); err != nil {
		return nil, err
	}
	g.Dump(result.Token)

	return result, nil
}

func (s *AuthService) CheckAuth(ctx context.Context, req *pb.CheckAuthRequest) (*pb.CheckAuthReply, error) {
	token, err := s.uc.CheckToken(ctx, &biz.Auth{
		Token: req.Token,
	})
	if err != nil {
		return nil, err
	}

	resp := &pb.CheckAuthReply{}
	if err := gconv.Structs(token, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
