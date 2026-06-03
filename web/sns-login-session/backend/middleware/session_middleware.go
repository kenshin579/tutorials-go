package middleware

import (
	"net/http"

	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/service"
	"github.com/labstack/echo/v4"
)

const SessionCookieName = "session_id"

// SessionAuth는 세션 쿠키를 검증하는 미들웨어
func SessionAuth(sessionService *service.SessionService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(SessionCookieName)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "세션 쿠키가 필요합니다")
			}
			userID, err := sessionService.Validate(cookie.Value)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 세션")
			}
			c.Set("user_id", userID)
			return next(c)
		}
	}
}
