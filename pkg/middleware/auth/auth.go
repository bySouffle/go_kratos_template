package auth

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	jwt.RegisteredClaims
}

func NewJwt(cookieName string, jwtSecret string, duration time.Duration) *JWT {
	return &JWT{
		cookieName: cookieName,
		jwtSecret:  jwtSecret,
		jwtTimeout: duration,
	}
}

func (j *JWT) GenerateToken(id int64, name string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test",
			Subject:   "sub",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.jwtTimeout)),
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间
			ID:        gconv.String(id),
		},
		ID:   id,
		Name: name,
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(j.jwtSecret))
	if err != nil {
		panic(err)
	}
	return tokenString
}

func (j *JWT) Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(j.jwtSecret), nil // 这是我的secret
	}
}

func (j *JWT) ParseToken(token string) (*Claims, error) {
	parseToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parseToken.Claims.(*Claims); ok && parseToken.Valid {
		return claims, err
	} else {
		return nil, err
	}
}

func (j *JWT) JWTAuth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			ctx2, span := otel.Tracer("Middleware").Start(ctx, "JWT", trace.WithSpanKind(trace.SpanKindInternal))
			defer span.End()
			if tr, ok := transport.FromServerContext(ctx2); ok {
				tokenString := tr.RequestHeader().Get("Authorization")
				auths := strings.SplitN(tokenString, " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], "Token") {
					return nil, errors.New("jwt token missing")
				}
				span.SetAttributes(attribute.String("token", auths[1]))
				token, err := jwt.ParseWithClaims(auths[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(j.jwtSecret), nil
				})
				//token, err := jwt.Parse(auths[1], func(token *jwt.Token) (interface{}, error) {
				//	// Don't forget to validate the alg is what you expect:
				//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				//		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				//	}
				//	return []byte(j.jwtSecret), nil
				//})

				if err != nil {
					return nil, err
				}

				if claims, ok := token.Claims.(*Claims); ok && token.Valid {
					// put CurrentUser into ctx
					ctx = j.WithContext(ctx, claims)
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
