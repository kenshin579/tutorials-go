package main

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed config/default.yaml
var defaultConfig []byte

func loadConfig(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultConfig, nil
	}
	return data, nil
}

func TestEmbed_Config_FallbackToDefault(t *testing.T) {
	// 존재하지 않는 경로 → 임베디드 기본값 사용
	data, err := loadConfig("/nonexistent/config.yaml")
	assert.NoError(t, err)
	assert.Contains(t, string(data), "port: 8080")
	assert.Contains(t, string(data), "name: my-app")
}

func TestEmbed_Config_UseExternalFile(t *testing.T) {
	// 임시 외부 설정 파일 생성
	tmpDir := t.TempDir()
	externalConfig := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(externalConfig, []byte("server:\n  port: 9090\n"), 0644)
	assert.NoError(t, err)

	// 외부 파일이 있으면 외부 파일 사용
	data, err := loadConfig(externalConfig)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "port: 9090")
}
