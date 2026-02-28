package main

import (
	"log"

	"github.com/kenshin579/tutorials-go/web/sns-login/backend/config"
	"github.com/kenshin579/tutorials-go/web/sns-login/backend/handler"
	customMiddleware "github.com/kenshin579/tutorials-go/web/sns-login/backend/middleware"
	"github.com/kenshin579/tutorials-go/web/sns-login/backend/model"
	"github.com/kenshin579/tutorials-go/web/sns-login/backend/provider"
	"github.com/kenshin579/tutorials-go/web/sns-login/backend/repository"
	"github.com/kenshin579/tutorials-go/web/sns-login/backend/service"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 설정 로드
	cfg := config.Load()

	// SQLite 연결
	db, err := gorm.Open(sqlite.Open("data/app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("DB 연결 실패:", err)
	}

	// 테이블 자동 마이그레이션
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("마이그레이션 실패:", err)
	}

	// 의존성 초기화
	googleProvider := provider.NewGoogleProvider(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
	)

	providers := map[string]provider.OAuthProvider{
		"google": googleProvider,
	}

	userRepo := repository.NewUserRepository(db)
	tokenService := service.NewTokenService(cfg.JWTSecret)
	authService := service.NewAuthService(providers, userRepo, tokenService)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(authService)

	// Echo 서버 설정
	e := echo.New()

	// 미들웨어
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// 라우트 등록
	api := e.Group("/api")

	// 인증 라우트 (공개)
	auth := api.Group("/auth")
	auth.GET("/:provider/url", authHandler.GetAuthURL)
	auth.GET("/:provider/callback", authHandler.HandleCallback)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)

	// 사용자 라우트 (인증 필요)
	user := api.Group("/user", customMiddleware.JWTAuth(tokenService))
	user.GET("/me", userHandler.GetMe)

	// 헬스체크
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	log.Printf("서버 시작: http://localhost:%s", cfg.ServerPort)
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}
