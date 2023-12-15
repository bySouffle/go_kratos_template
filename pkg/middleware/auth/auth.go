package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type JWT struct {
	cookieName string
	jwtSecret  string
	jwtTimeout time.Duration
}

type Claims struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

func NewJwt(cookieName string, jwtSecret string, duration time.Duration) *JWT {
	return &JWT{
		cookieName: cookieName,
		jwtSecret:  jwtSecret,
		jwtTimeout: duration,
	}
}

func (j *JWT) GenerateToken(id int64, name string) string {
	nowTime := time.Now()
	expireTime := nowTime.Add(j.jwtTimeout)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":   id,
		"Name": name,
		"nbf":  expireTime.Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(j.jwtSecret))
	if err != nil {
		panic(err)
	}
	return tokenString
}

func (j *JWT) JWTAuth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				tokenString := tr.RequestHeader().Get("Authorization")
				auths := strings.SplitN(tokenString, " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], "Token") {
					return nil, errors.New("jwt token missing")
				}

				token, err := jwt.Parse(auths[1], func(token *jwt.Token) (interface{}, error) {
					// Don't forget to validate the alg is what you expect:
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(j.jwtSecret), nil
				})

				if err != nil {
					return nil, err
				}

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					// put CurrentUser into ctx
					if u, ok := claims["ID"]; ok {
						ctx = j.WithContext(ctx, &Claims{ID: u.(int64)})
					}
				} else {
					return nil, errors.New("token Invalid")
				}
			}
			return handler(ctx, req)
		}
	}
}
func (j *JWT) FromContext(ctx context.Context) *Claims {
	return ctx.Value(j.cookieName).(*Claims)
}

func (j *JWT) WithContext(ctx context.Context, session *Claims) context.Context {
	return context.WithValue(ctx, j.cookieName, session)
}
