package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	jwthelper "github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/pkg/jwt"
)

const ctxUserID = "user_id"

// JWTAuth는 Authorization: Bearer <token> 헤더를 검증하고 user_id를 컨텍스트에 주입한다.
// 토큰이 누락되거나 검증 실패 시 401을 반환한다.
func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
			}
			tok := strings.TrimPrefix(h, "Bearer ")
			claims, err := jwthelper.Parse(tok, secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}
			c.Set(ctxUserID, claims.UserID)
			return next(c)
		}
	}
}

// UserIDFrom은 JWTAuth 미들웨어가 주입한 user_id를 컨텍스트에서 추출한다.
// 미들웨어가 적용되지 않은 라우트에서 호출하면 0을 반환한다.
func UserIDFrom(c echo.Context) uint {
	v, _ := c.Get(ctxUserID).(uint)
	return v
}
