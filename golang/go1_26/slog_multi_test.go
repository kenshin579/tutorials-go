package go1_26_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

func TestSlogNewMultiHandler(t *testing.T) {
	var textBuf, jsonBuf bytes.Buffer

	textHandler := slog.NewTextHandler(&textBuf, &slog.HandlerOptions{Level: slog.LevelInfo})
	jsonHandler := slog.NewJSONHandler(&jsonBuf, &slog.HandlerOptions{Level: slog.LevelInfo})

	// NewMultiHandler: 여러 핸들러에 동시 출력
	multi := slog.NewMultiHandler(textHandler, jsonHandler)
	logger := slog.New(multi)

	logger.Info("user login", "user", "alice", "ip", "192.168.1.1")

	fmt.Println("--- Text Output ---")
	fmt.Print(textBuf.String())
	fmt.Println("--- JSON Output ---")
	fmt.Print(jsonBuf.String())

	// 두 핸들러 모두 출력되었는지 확인
	if !strings.Contains(textBuf.String(), "user login") {
		t.Error("text handler should contain 'user login'")
	}
	if !strings.Contains(jsonBuf.String(), "user login") {
		t.Error("json handler should contain 'user login'")
	}
}

func TestSlogMultiHandlerWithLevels(t *testing.T) {
	var infoBuf, errorBuf bytes.Buffer

	infoHandler := slog.NewTextHandler(&infoBuf, &slog.HandlerOptions{Level: slog.LevelInfo})
	errorHandler := slog.NewTextHandler(&errorBuf, &slog.HandlerOptions{Level: slog.LevelError})

	multi := slog.NewMultiHandler(infoHandler, errorHandler)
	logger := slog.New(multi)

	// Info 레벨 메시지
	logger.Info("info message")
	// Error 레벨 메시지
	logger.Error("error message")

	fmt.Println("--- Info Handler ---")
	fmt.Print(infoBuf.String())
	fmt.Println("--- Error Handler ---")
	fmt.Print(errorBuf.String())

	// Info 핸들러: 두 메시지 모두 출력
	if !strings.Contains(infoBuf.String(), "info message") {
		t.Error("info handler should contain 'info message'")
	}
	if !strings.Contains(infoBuf.String(), "error message") {
		t.Error("info handler should also contain 'error message'")
	}

	// Error 핸들러: error 메시지만 출력
	if strings.Contains(errorBuf.String(), "info message") {
		t.Error("error handler should not contain 'info message'")
	}
	if !strings.Contains(errorBuf.String(), "error message") {
		t.Error("error handler should contain 'error message'")
	}
}
