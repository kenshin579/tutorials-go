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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
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
		Sub               string `json:"sub"`
		PreferredUsername string `json:"preferred_username"`
		Email             string `json:"email"`
		GivenName         string `json:"given_name"`
		FamilyName        string `json:"family_name"`
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, introspectURL, strings.NewReader(data.Encode()))
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userURL, nil)
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
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
