package service

import (
	"testing"
)

func TestGenerateTokenPair_TokenTypes(t *testing.T) {
	s := NewTokenService("test-secret")

	pair, err := s.GenerateTokenPair(42)
	if err != nil {
		t.Fatalf("토큰 생성 실패: %v", err)
	}

	accessClaims, err := s.ValidateToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("access 토큰 검증 실패: %v", err)
	}
	if accessClaims.TokenType != "access" {
		t.Errorf("access 토큰 타입 기대값 access, 실제 %q", accessClaims.TokenType)
	}
	if accessClaims.UserID != 42 {
		t.Errorf("UserID 기대값 42, 실제 %d", accessClaims.UserID)
	}

	refreshClaims, err := s.ValidateToken(pair.RefreshToken)
	if err != nil {
		t.Fatalf("refresh 토큰 검증 실패: %v", err)
	}
	if refreshClaims.TokenType != "refresh" {
		t.Errorf("refresh 토큰 타입 기대값 refresh, 실제 %q", refreshClaims.TokenType)
	}
}
