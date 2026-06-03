package main

import (
	"log"
	"time"

	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/config"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/handler"
	customMiddleware "github.com/kenshin579/tutorials-go/web/sns-login-session/backend/middleware"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/model"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/provider"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/repository"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/service"
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
	if err := db.AutoMigrate(&model.User{}, &model.Session{}); err != nil {
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
	sessionRepo := repository.NewSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepo, 7*24*time.Hour) // 세션 7일
	authService := service.NewAuthService(providers, userRepo, sessionService)

	authHandler := handler.NewAuthHandler(authService, cfg.FrontendURL)
	userHandler := handler.NewUserHandler(authService)

	// Echo 서버 설정
	e := echo.New()

	// 미들웨어
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// 라우트 등록
	api := e.Group("/api")

	auth := api.Group("/auth")
	auth.GET("/:provider/url", authHandler.GetAuthURL)
	auth.GET("/:provider/callback", authHandler.HandleCallback)
	auth.POST("/logout", authHandler.Logout)

	user := api.Group("/user", customMiddleware.SessionAuth(sessionService))
	user.GET("/me", userHandler.GetMe)

	// 헬스체크
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	log.Printf("서버 시작: http://localhost:%s", cfg.ServerPort)
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}
