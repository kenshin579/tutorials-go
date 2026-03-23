# CLAUDE.md - Ralph Loop용 설정

<!--
  Ralph Loop에서 Claude Code가 매 반복마다 읽는 설정 파일입니다.
  프로젝트 컨텍스트와 규칙을 여기에 정의합니다.
-->

## 프로젝트 개요

간단한 Go HTTP API 서버 프로젝트. Ralph Loop 데모용.

## 핵심 규칙

### 작업 방식
- IMPLEMENTATION_PLAN.md를 참조하여 작업 선택
- **한 번에 하나의 작업만** 수행
- 구현 전 코드베이스 검색으로 중복 확인

### 코드 품질
- `go test ./...` 통과 필수
- `go build .` 성공 필수
- 플레이스홀더/TODO 금지 — 완전한 구현만

### 커밋 규칙
- 테스트 통과 후에만 커밋
- 커밋 메시지: 한국어, `[#이슈번호] 간결한 설명` 형식

## 빌드 & 테스트

```bash
go build -o server .
go test -v ./...
```
