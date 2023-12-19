package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go_kratos_template/app/auth/internal/biz"
	"go_kratos_template/pkg/cache"
	"time"
)

type authRepo struct {
	data *Data
	log  *log.Helper
}

func (a authRepo) Create(ctx context.Context, auth *biz.Auth) error {
	key, err := cache.BuildCacheKey(auth.Name, auth.ID)
	if err != nil {
		return err
	}
	return a.data.redis.Set(ctx, key, auth.Token, time.Duration(0)).Err()
}

func (a authRepo) Delete(ctx context.Context, auth *biz.Auth) error {
	key, err := cache.BuildCacheKey(auth.Name, auth.ID)
	if err != nil {
		return err
	}
	return a.data.redis.Del(ctx, key).Err()
}

func (a authRepo) Update(ctx context.Context, auth *biz.Auth) error {
	key, err := cache.BuildCacheKey(auth.Name, auth.ID)
	if err != nil {
		return err
	}
	return a.data.redis.Set(ctx, key, auth.Token, time.Duration(0)).Err()
}

func (a authRepo) FindByAuth(ctx context.Context, auth *biz.Auth) (*biz.Auth, error) {
	key, err := cache.BuildCacheKey(auth.Name, auth.ID)
	if err != nil {
		return nil, err
	}
	val, redisErr := a.data.redis.Get(ctx, key).Result()

	if redisErr != nil {
		return nil, redisErr
	}
	result := &biz.Auth{
		Token: val,
	}
	return result, nil
}

func (a authRepo) ListByName(ctx context.Context, auths []*biz.Auth) ([]*biz.Auth, error) {
	for _, auth := range auths {
		key, err := cache.BuildCacheKey(auth.Name, auth.ID)
		if err != nil {
			continue
		}
		val, redisErr := a.data.redis.Get(ctx, key).Result()

		if redisErr != nil {
			continue
		}
		auth.Token = val
	}
	return auths, nil
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
