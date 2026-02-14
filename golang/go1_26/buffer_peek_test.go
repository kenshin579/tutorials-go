package go1_26_test

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBufferPeek(t *testing.T) {
	buf := bytes.NewBufferString("hello world")

	// Peek: 버퍼를 진행시키지 않고 다음 n 바이트를 확인
	peeked, err := buf.Peek(5)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Peek(5): %q\n", peeked) // "hello"
	fmt.Printf("Len after Peek: %d\n", buf.Len()) // 11 (변경 없음)

	if string(peeked) != "hello" {
		t.Errorf("expected 'hello', got '%s'", peeked)
	}

	// Peek 후에도 버퍼 크기 유지
	if buf.Len() != 11 {
		t.Errorf("buffer length should not change after Peek: expected 11, got %d", buf.Len())
	}

	// Read와 비교: Read는 버퍼를 진행시킴
	readBuf := make([]byte, 5)
	n, err := buf.Read(readBuf)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Read(%d): %q\n", n, readBuf[:n]) // "hello"
	fmt.Printf("Len after Read: %d\n", buf.Len()) // 6

	if buf.Len() != 6 {
		t.Errorf("buffer length should change after Read: expected 6, got %d", buf.Len())
	}
}

func TestBufferPeekThenRead(t *testing.T) {
	buf := bytes.NewBufferString("Go 1.26")

	// Peek으로 먼저 확인
	header, err := buf.Peek(2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Header: %q\n", header) // "Go"

	// 전체 내용 읽기 (Peek이 영향 없음)
	all := buf.String()
	fmt.Printf("Full: %q\n", all) // "Go 1.26"

	if all != "Go 1.26" {
		t.Errorf("expected 'Go 1.26', got '%s'", all)
	}
}

func TestBufferPeekEmpty(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	// 빈 버퍼에서 Peek
	_, err := buf.Peek(1)
	if err == nil {
		t.Error("expected error on empty buffer Peek")
	}
	fmt.Printf("Empty buffer Peek error: %v\n", err)
}
