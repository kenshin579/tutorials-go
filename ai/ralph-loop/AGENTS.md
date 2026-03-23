# AGENTS.md - 에이전트 운영 가이드

<!--
  이 파일은 Ralph Loop 실행 중 발견된 운영 지식을 축적하는 곳입니다.
  매 루프 반복마다 자동으로 로드되어 에이전트에게 컨텍스트를 제공합니다.

  ✅ 필수 항목: 빌드/테스트 명령어, 검증 규칙
  📌 권장 항목: 발견된 패턴, 주의사항
-->

## 빌드 & 테스트 명령어 [필수]

```bash
# 빌드
go build -o server .

# 테스트 실행
go test -v ./...

# 테스트 + 커버리지
go test -cover ./...

# 린트 (설치 필요: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
golangci-lint run ./...
```

## 검증 규칙 [필수]

- 커밋 전 `go test ./...` 통과 필수
- `go build .` 성공 필수
- `go vet ./...` 경고 없어야 함

## 프로젝트 구조 [필수]

```
ai/ralph-loop/
├── main.go          # HTTP 서버 (엔트리포인트)
├── main_test.go     # 테스트
├── specs/           # 요구사항 스펙 (1 파일 = 1 관심사)
├── prd.json         # 진행 상태 추적
└── IMPLEMENTATION_PLAN.md  # 구현 계획 (Ralph가 생성/관리)
```

## 코딩 컨벤션 [필수]

- Go 표준 프로젝트 레이아웃 준수
- 에러는 즉시 처리 (`if err != nil`)
- JSON 태그는 camelCase가 아닌 snake_case 사용
- 테스트 함수명: `TestXxx_설명` 형식

## 발견된 패턴 [권장]

<!--
  Ralph Loop 실행 중 발견한 패턴이나 주의사항을 여기에 추가하세요.
  예시:
  - HTTP 핸들러는 writeJSON 헬퍼를 사용할 것
  - 새 엔드포인트 추가 시 routes() 메서드에 등록할 것
-->

- `writeJSON()` 헬퍼로 JSON 응답 통일
- 새 엔드포인트는 `Server.routes()`에 등록
- 환경변수 `PORT`로 포트 설정 (기본값: 8080)
