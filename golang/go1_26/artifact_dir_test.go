package go1_26_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestArtifactDir(t *testing.T) {
	// ArtifactDir: 테스트 아티팩트를 저장할 디렉토리 반환
	dir := t.ArtifactDir()
	fmt.Printf("Artifact directory: %s\n", dir)

	// 테스트 결과 파일 저장
	content := []byte("test output data: Go 1.26 artifact example")
	outputPath := filepath.Join(dir, "output.txt")
	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		t.Fatalf("failed to write artifact: %v", err)
	}

	// 파일이 정상적으로 저장되었는지 확인
	readBack, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read artifact: %v", err)
	}
	if string(readBack) != string(content) {
		t.Errorf("artifact content mismatch")
	}

	fmt.Printf("Artifact saved to: %s\n", outputPath)
}

func TestArtifactDirMultipleFiles(t *testing.T) {
	dir := t.ArtifactDir()

	// 여러 아티팩트 파일 저장
	files := map[string]string{
		"result.json": `{"status": "pass", "count": 42}`,
		"log.txt":     "2026-02-14 test log entry",
	}

	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", name, err)
		}
		fmt.Printf("Saved: %s\n", path)
	}

	// 디렉토리 내 파일 목록 확인
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 2 {
		t.Errorf("expected at least 2 files, got %d", len(entries))
	}
}
