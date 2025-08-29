package main

import (
	"log"
	"net/http"

	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/handler"
	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/repository"
	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/usecase"
	appconfig "github.com/kenshin579/tutorials-go/keycloak/backend/pkg/config"
	appmw "github.com/kenshin579/tutorials-go/keycloak/backend/pkg/middleware"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

func main() {
	// 설정 로드
	cfg := appconfig.NewConfig()

	// Echo 인스턴스 생성
	e := echo.New()

	// 미들웨어 설정 (Echo 기본)
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Repository 생성
	keycloakRepo := repository.NewKeycloakRepository(
		cfg.Keycloak.BaseURL,
		cfg.Keycloak.Realm,
		cfg.Keycloak.ClientID,
		cfg.Keycloak.ClientSecret,
	)

	// UseCase 생성
	userUseCase := usecase.NewUserUseCase(keycloakRepo)

	// Handler 생성
	userHandler := handler.NewUserHandler(userUseCase)

	// 라우트 설정
	api := e.Group("/api")

	// 인증이 필요한 라우트
	protected := api.Group("/protected")
	protected.Use(appmw.AuthMiddleware(userUseCase))
	protected.GET("/user", userHandler.GetUserInfo)

	// 인증이 필요없는 라우트
	api.GET("/validate", userHandler.ValidateToken)

	// 서버 시작
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(e.Start(":" + cfg.Server.Port))
}
