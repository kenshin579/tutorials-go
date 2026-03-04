package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/kenshin579/tutorials-go/golang/middleware/custom"
)

var signingKey = []byte("my-secret-key-for-demo")

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	e := echo.New()

	// 글로벌 미들웨어 체인 (순서 중요: Recover → RequestID → Logger → CORS)
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(custom.ZapLogger(logger))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// 공개 엔드포인트
	e.GET("/health", healthHandler)
	e.POST("/login", loginHandler)

	// 인증 필요 API 그룹
	api := e.Group("/api")
	api.Use(custom.JWTAuth(custom.JWTConfig{
		SigningKey: signingKey,
	}))
	api.GET("/profile", profileHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func loginHandler(c echo.Context) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "잘못된 요청")
	}

	// 데모용 하드코딩 인증
	if req.Username != "admin" || req.Password != "password" {
		return echo.NewHTTPError(http.StatusUnauthorized, "잘못된 인증 정보")
	}

	// JWT 토큰 생성
	claims := &custom.Claims{
		UserID:   "user-123",
		Username: req.Username,
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "토큰 생성 실패")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}

func profileHandler(c echo.Context) error {
	claims := c.Get("user").(*custom.Claims)
	return c.JSON(http.StatusOK, map[string]string{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"role":     claims.Role,
	})
}
