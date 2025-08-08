# Keycloak 기반 인증 시스템 구현 가이드

## 개요
이 문서는 PRD.md의 요구사항에 따라 Keycloak 기반 인증 시스템을 구현하는 상세한 가이드입니다.

## 1. Keycloak 설정 개선

### 1.1 Keycloak 환경 확인

현재 Keycloak이 이미 Docker로 실행 중입니다:
- **Keycloak URL**: http://localhost:8080
- **Admin Console**: http://localhost:8080/admin
- **Admin 계정**: admin / admin

기존 `infra/docker_run.sh` 스크립트로 실행된 상태입니다.

### 1.2 Keycloak 클라이언트 설정

#### 1.2.1 React 클라이언트 설정
- **Client ID**: `myclient`
- **Client Protocol**: `openid-connect`
- **Access Type**: `public`
- **Standard Flow Enabled**: `ON`
- **Valid Redirect URIs**: 
  - `http://localhost:3000/*`
  - `http://localhost:3000`
- **Web Origins**: 
  - `http://localhost:3000`
  - `+` (모든 origin 허용 - 개발용)

#### 1.2.2 백엔드 클라이언트 설정
- **Client ID**: `mybackend`
- **Client Protocol**: `openid-connect`
- **Access Type**: `confidential`
- **Service Accounts Enabled**: `ON`
- **Valid Redirect URIs**: 
  - `http://localhost:8081/api/*`
  - `http://localhost:8081/api`
- **Web Origins**: 
  - `http://localhost:3000` (프론트엔드에서 API 호출)

### 1.3 Realm 설정
- **Realm Name**: `myrealm`
- **Login Theme**: `keycloak`
- **Account Theme**: `keycloak`
- **Admin Theme**: `keycloak`
- **Email Theme**: `keycloak`

## 2. 백엔드 구현 (Golang + Echo + Clean Architecture)

### 2.1 프로젝트 구조 생성

```bash
mkdir -p keycloak/backend/{cmd/server,internal/{domain,usecase,repository,handler},pkg/middleware}
cd keycloak/backend
go mod init github.com/kenshin579/tutorials-go/keycloak/backend
```

### 2.2 의존성 설치

```bash
# Go 1.24.5 사용
go mod init github.com/kenshin579/tutorials-go/keycloak/backend
go mod tidy

# 필요한 의존성들
go get github.com/labstack/echo/v4@latest
go get github.com/golang-jwt/jwt/v5@latest
```

### 2.3 Domain Layer

```go
// internal/domain/user.go
package domain

import "context"

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserRepository interface {
	GetUserByID(ctx context.Context, userID string) (*User, error)
	GetUserInfo(ctx context.Context, token string) (*User, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
}

type UserUseCase interface {
	GetUserInfo(ctx context.Context, token string) (*User, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
}
```

```go
// internal/domain/auth.go
package domain

import "context"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type AuthRepository interface {
	Login(ctx context.Context, username, password string) (*AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}
```

### 2.4 Repository Layer

```go
// internal/repository/keycloak_repository.go
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/domain"
)

type KeycloakRepository struct {
	baseURL      string
	realm        string
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

func NewKeycloakRepository(baseURL, realm, clientID, clientSecret string) *KeycloakRepository {
	return &KeycloakRepository{
		baseURL:      baseURL,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{},
	}
}

// GetUserInfo - Keycloak UserInfo 엔드포인트 호출
func (r *KeycloakRepository) GetUserInfo(ctx context.Context, token string) (*domain.User, error) {
	userInfoURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", r.baseURL, r.realm)
	
	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call userinfo endpoint: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userinfo endpoint returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var userInfo struct {
		Sub                string `json:"sub"`
		PreferredUsername  string `json:"preferred_username"`
		Email              string `json:"email"`
		GivenName          string `json:"given_name"`
		FamilyName         string `json:"family_name"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo response: %w", err)
	}
	
	return &domain.User{
		ID:        userInfo.Sub,
		Username:  userInfo.PreferredUsername,
		Email:     userInfo.Email,
		FirstName: userInfo.GivenName,
		LastName:  userInfo.FamilyName,
	}, nil
}

// ValidateToken - Keycloak Token Introspection 엔드포인트 호출
func (r *KeycloakRepository) ValidateToken(ctx context.Context, token string) (bool, error) {
	introspectURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect", r.baseURL, r.realm)
	
	data := url.Values{}
	data.Set("token", token)
	data.Set("client_id", r.clientID)
	data.Set("client_secret", r.clientSecret)
	
	req, err := http.NewRequestWithContext(ctx, "POST", introspectURL, strings.NewReader(data.Encode()))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to call token introspection endpoint: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("token introspection endpoint returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var result struct {
		Active bool `json:"active"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode introspection response: %w", err)
	}
	
	return result.Active, nil
}

// GetUserByID - Keycloak Admin API를 사용하여 사용자 정보 조회
func (r *KeycloakRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	// Admin 토큰 먼저 획득
	adminToken, err := r.getAdminToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin token: %w", err)
	}
	
	userURL := fmt.Sprintf("%s/admin/realms/%s/users/%s", r.baseURL, r.realm, userID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", userURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call admin API: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("admin API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var user struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}
	
	return &domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

// getAdminToken - Keycloak Admin 토큰 획득
func (r *KeycloakRepository) getAdminToken(ctx context.Context) (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", r.baseURL)
	
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", "admin-cli")
	data.Set("client_secret", r.clientSecret)
	
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create admin token request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call admin token endpoint: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("admin token endpoint returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode admin token response: %w", err)
	}
	
	return tokenResp.AccessToken, nil
}
```

### 2.5 UseCase Layer

```go
// internal/usecase/user_usecase.go
package usecase

import (
	"context"
	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/domain"
)

type UserUseCaseImpl struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &UserUseCaseImpl{
		userRepo: userRepo,
	}
}

func (u *UserUseCaseImpl) GetUserInfo(ctx context.Context, token string) (*domain.User, error) {
	return u.userRepo.GetUserInfo(ctx, token)
}

func (u *UserUseCaseImpl) ValidateToken(ctx context.Context, token string) (bool, error) {
	return u.userRepo.ValidateToken(ctx, token)
}
```

### 2.6 Handler Layer

```go
// internal/handler/user_handler.go
package handler

import (
	"context"
	"net/http"
	"strings"
	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/domain"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(userUseCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) GetUserInfo(c echo.Context) error {
	token := extractToken(c)
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
	}
	
	ctx := c.Request().Context()
	user, err := h.userUseCase.GetUserInfo(ctx, token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}
	
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ValidateToken(c echo.Context) error {
	token := extractToken(c)
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
	}
	
	ctx := c.Request().Context()
	valid, err := h.userUseCase.ValidateToken(ctx, token)
	if err != nil || !valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}
	
	return c.JSON(http.StatusOK, map[string]bool{"valid": true})
}

func extractToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	
	return authHeader
}
```

### 2.7 Middleware

```go
// pkg/middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"
	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/domain"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(userUseCase domain.UserUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := extractToken(c)
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
			}
			
			ctx := c.Request().Context()
			valid, err := userUseCase.ValidateToken(ctx, token)
			if err != nil || !valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			
			// 토큰을 컨텍스트에 저장
			c.Set("token", token)
			
			return next(c)
		}
	}
}

func extractToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	
	return authHeader
}
```

### 2.8 Configuration

간단한 구현을 위해 하드코딩된 설정을 사용합니다:

```go
// pkg/config/config.go
package config

type Config struct {
	Server   ServerConfig
	Keycloak KeycloakConfig
}

type ServerConfig struct {
	Port string
}

type KeycloakConfig struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8081",
		},
		Keycloak: KeycloakConfig{
			BaseURL:      "http://localhost:8080",
			Realm:        "myrealm",
			ClientID:     "mybackend",
			ClientSecret: "your-client-secret", // Keycloak에서 생성된 클라이언트 시크릿으로 변경
		},
	}
}
```

### 2.9 Main Application

```go
// cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/domain"
	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/handler"
	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/repository"
	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/usecase"
	"github.com/kenshin579/tutorials-go/keycloak/backend/pkg/config"
	"github.com/kenshin579/tutorials-go/keycloak/backend/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 설정 로드
	cfg := config.NewConfig()
	
	// Echo 인스턴스 생성
	e := echo.New()
	
	// 미들웨어 설정
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
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
	protected.Use(middleware.AuthMiddleware(userUseCase))
	protected.GET("/user", userHandler.GetUserInfo)
	
	// 인증이 필요없는 라우트
	api.GET("/validate", userHandler.ValidateToken)
	
	// 서버 시작
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(e.Start(":" + cfg.Server.Port))
}
```

### 2.10 설정

설정은 `pkg/config/config.go` 파일에서 하드코딩되어 있습니다.

주요 설정값:
- **서버 포트**: 8081
- **Keycloak URL**: http://localhost:8080
- **Realm**: myrealm
- **Client ID**: mybackend
- **Client Secret**: Keycloak에서 생성된 실제 시크릿으로 변경 필요

## 3. 프론트엔드 구현 (React)

### 3.1 프로젝트 생성

```bash
npx create-react-app keycloak/frontend
cd keycloak/frontend
npm install keycloak-js axios react-router-dom
```

### 3.2 Keycloak 설정

```javascript
// src/services/keycloak.js
import Keycloak from 'keycloak-js';

const keycloakConfig = {
  url: 'http://localhost:8080',
  realm: 'myrealm',
  clientId: 'myclient'
};

const keycloak = new Keycloak(keycloakConfig);

export const initKeycloak = () => {
  return new Promise((resolve, reject) => {
    keycloak.init({
      onLoad: 'check-sso',
      silentCheckSsoRedirectUri: window.location.origin + '/silent-check-sso.html',
      pkceMethod: 'S256'
    })
    .then((authenticated) => {
      resolve({ authenticated, keycloak });
    })
    .catch((error) => {
      reject(error);
    });
  });
};

export default keycloak;
```

### 3.3 API 서비스

```javascript
// src/services/api.js
import axios from 'axios';
import keycloak from './keycloak';

const API_BASE_URL = 'http://localhost:8081/api';

const api = axios.create({
  baseURL: API_BASE_URL,
});

// 요청 인터셉터 - 토큰 추가
api.interceptors.request.use(
  (config) => {
    if (keycloak.token) {
      config.headers.Authorization = `Bearer ${keycloak.token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 응답 인터셉터 - 토큰 갱신
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      try {
        await keycloak.updateToken(30);
        error.config.headers.Authorization = `Bearer ${keycloak.token}`;
        return api.request(error.config);
      } catch (refreshError) {
        keycloak.logout();
        return Promise.reject(refreshError);
      }
    }
    return Promise.reject(error);
  }
);

export const getUserInfo = () => api.get('/protected/user');
export const validateToken = () => api.get('/validate');

export default api;
```

### 3.4 컴포넌트 구현

```jsx
// src/components/Login.js
import React from 'react';
import keycloak from '../services/keycloak';

const Login = () => {
  const handleLogin = () => {
    keycloak.login();
  };

  return (
    <div className="login-container">
      <h2>Welcome to Keycloak Auth Demo</h2>
      <p>Please login to continue</p>
      <button onClick={handleLogin} className="login-button">
        Login with Keycloak
      </button>
    </div>
  );
};

export default Login;
```

```jsx
// src/components/Logout.js
import React from 'react';
import keycloak from '../services/keycloak';

const Logout = () => {
  const handleLogout = () => {
    keycloak.logout();
  };

  return (
    <div className="logout-container">
      <button onClick={handleLogout} className="logout-button">
        Logout
      </button>
    </div>
  );
};

export default Logout;
```

```jsx
// src/components/UserInfo.js
import React, { useState, useEffect } from 'react';
import { getUserInfo } from '../services/api';
import keycloak from '../services/keycloak';

const UserInfo = () => {
  const [userInfo, setUserInfo] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const response = await getUserInfo();
        setUserInfo(response.data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    if (keycloak.authenticated) {
      fetchUserInfo();
    }
  }, []);

  if (loading) return <div>Loading user info...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!userInfo) return <div>No user info available</div>;

  return (
    <div className="user-info">
      <h3>User Information</h3>
      <div className="user-details">
        <p><strong>Username:</strong> {userInfo.username}</p>
        <p><strong>Email:</strong> {userInfo.email}</p>
        <p><strong>First Name:</strong> {userInfo.firstName}</p>
        <p><strong>Last Name:</strong> {userInfo.lastName}</p>
        <p><strong>Full Name:</strong> {userInfo.firstName} {userInfo.lastName}</p>
      </div>
    </div>
  );
};

export default UserInfo;
```

### 3.5 메인 앱 컴포넌트

```jsx
// src/App.js
import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { initKeycloak } from './services/keycloak';
import Login from './components/Login';
import UserInfo from './components/UserInfo';
import Logout from './components/Logout';
import './App.css';

function App() {
  const [keycloak, setKeycloak] = useState(null);
  const [authenticated, setAuthenticated] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    initKeycloak()
      .then(({ authenticated, keycloak }) => {
        setKeycloak(keycloak);
        setAuthenticated(authenticated);
        setLoading(false);
      })
      .catch((error) => {
        console.error('Keycloak init failed:', error);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <h1>Keycloak Authentication Demo</h1>
          {authenticated && <Logout />}
        </header>
        
        <main>
          <Routes>
            <Route 
              path="/" 
              element={
                authenticated ? (
                  <div>
                    <UserInfo />
                  </div>
                ) : (
                  <Login />
                )
              } 
            />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
```

### 3.6 스타일링

```css
/* src/App.css */
.App {
  text-align: center;
  padding: 20px;
}

.App-header {
  background-color: #282c34;
  padding: 20px;
  color: white;
  margin-bottom: 20px;
}

.login-container, .user-info, .logout-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

.login-button, .logout-button {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 5px;
  cursor: pointer;
  font-size: 16px;
}

.login-button:hover, .logout-button:hover {
  background-color: #0056b3;
}

.user-details {
  text-align: left;
  background-color: #f8f9fa;
  padding: 20px;
  border-radius: 5px;
  margin-top: 20px;
}

.user-details p {
  margin: 10px 0;
}
```

## 4. 실행 방법

### 4.1 환경 시작

```bash
# Keycloak은 이미 실행 중 (localhost:8080)

# 백엔드 시작
cd keycloak/backend
go run cmd/server/main.go

# 프론트엔드 시작
cd ../frontend
npm start
```

### 4.2 접속 URL
- **Keycloak Admin Console**: http://localhost:8080
- **React 앱**: http://localhost:3000
- **백엔드 API**: http://localhost:8081

### 4.3 테스트 시나리오
1. React 앱 접속
2. "Login with Keycloak" 버튼 클릭
3. Keycloak 로그인 페이지에서 사용자 정보 입력
4. 로그인 후 사용자 정보 확인
5. 로그아웃 테스트

## 5. 추가 개선사항

### 5.1 보안 강화
- HTTPS 설정
- 토큰 만료 시간 조정

### 5.2 기능 확장
- 사용자 등록
- 비밀번호 변경
- 프로필 관리

### 5.3 모니터링
- 로그 설정
- 헬스체크 엔드포인트

이 구현 가이드를 따라하면 PRD.md의 모든 요구사항을 만족하는 Keycloak 기반 인증 시스템을 구축할 수 있습니다.
