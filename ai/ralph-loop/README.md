# Ralph Loop Demo

Ralph Loop(Ralph Wiggum Technique) 패턴을 보여주는 데모 프로젝트입니다.

## 프로젝트 구조

```
ai/ralph-loop/
├── loop.sh              # Ralph Loop 실행 스크립트
├── PROMPT_plan.md       # Planning 모드 프롬프트
├── PROMPT_build.md      # Building 모드 프롬프트
├── AGENTS.md            # 에이전트 운영 가이드
├── CLAUDE.md            # Claude Code 설정
├── specs/               # 요구사항 스펙 파일
│   ├── api-server.md    # API 서버 스펙
│   └── health-check.md  # 헬스체크 스펙
├── prd.json             # PRD JSON (진행 상태 추적)
├── main.go              # HTTP API 서버
├── main_test.go         # 테스트
└── go.mod               # Go 모듈 (루트 모듈에 포함)
```

## 실행 방법

### 서버 실행

```bash
go run .
# http://localhost:8080/api/hello
# http://localhost:8080/health
```

### 테스트

```bash
go test -v ./...
```

### Ralph Loop 실행

```bash
# Planning 모드 (계획 수립, 최대 3회)
./loop.sh plan 3

# Building 모드 (구현, 최대 10회)
./loop.sh build 10

# Building 모드 (무제한)
./loop.sh
```

## Ralph Loop 워크플로우

1. **Phase 1 - Specs**: `specs/` 디렉토리에 요구사항 스펙 작성
2. **Phase 2 - Planning**: `./loop.sh plan 3` 으로 IMPLEMENTATION_PLAN.md 생성
3. **Phase 3 - Building**: `./loop.sh build 10` 으로 자율 구현

## Best Practices

각 파일에 DO/DON'T 주석이 포함되어 있습니다. 자세한 내용은 각 파일을 참조하세요.

## 관련 블로그

- [Ralph Loop 완벽 가이드 - AI 에이전트 자율 개발 패턴](https://blog.advenoh.pe.kr)
