# Todo Web Application Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Echo (Go 1.26) + React 19 (Vite + TypeScript) 기반 in-memory Todo 웹 애플리케이션을 TDD 사이클로 처음부터 구현하면서 superpowers skill 사이클을 풀로 체험한다.

**Architecture:** 두 프로세스 분리 (BE :8080 / FE :5173, Vite proxy로 `/api/*` 포워딩). 백엔드는 `server` (라우팅) → `todo.Handler` (HTTP I/O) → `todo.Store` (도메인) 3계층. Store는 `sync.RWMutex` 보호 `map[string]Todo`. PATCH는 `map[string]json.RawMessage`로 1차 디코딩하여 키 존재/null 판별. 프론트는 단일 `useTodos` 훅이 useReducer로 상태 캡슐화, mutation 후 항상 list refetch.

**Tech Stack:**
- Backend: Go 1.26.0 (root `tutorials-go/go.mod`), Echo v4, google/uuid, testify/assert
- Frontend: React 19, Vite, TypeScript 5, Vitest, React Testing Library, MSW
- Spec 참조: `docs/superpowers/specs/2026-04-30-todo-app-design.md`

**Scope reference**: 모든 작업은 `/Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/` 하위에서 수행. Go 모듈 루트는 `tutorials-go/`. 모든 git 명령은 `tutorials-go/` 디렉토리에서 실행.

**커밋 규칙 (`CLAUDE.md` 준수):**
- 형식: `[feat/todo-app] <type>: <간결한 한국어 설명>`
- 타입: feat, fix, docs, style, refactor, test, chore
- 본 plan에서 모든 커밋은 `feat/todo-app` 브랜치에 적재

---

## Pre-flight

### Task A: 피처 브랜치 생성 + spec/plan 커밋

**Files:**
- Existing: `tutorials-go/ai/superpowers/todo/docs/superpowers/specs/2026-04-30-todo-app-design.md`
- Existing: `tutorials-go/ai/superpowers/todo/docs/superpowers/plans/2026-04-30-todo-app-plan.md` (본 문서)

- [ ] **Step 1: 현재 브랜치 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git status --short
git branch --show-current
```

Expected: 현재 master, working tree에 spec/plan 두 파일이 untracked로 보임.

- [ ] **Step 2: 피처 브랜치 생성 및 체크아웃**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git checkout -b feat/todo-app
```

Expected: `Switched to a new branch 'feat/todo-app'`

- [ ] **Step 3: spec + plan 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/docs/superpowers/specs/2026-04-30-todo-app-design.md \
        ai/superpowers/todo/docs/superpowers/plans/2026-04-30-todo-app-plan.md
git commit -m "$(cat <<'EOF'
[feat/todo-app] docs: superpowers todo app spec과 implementation plan 추가

* 학습용 Todo 웹 애플리케이션 (Echo + React + in-memory) 설계 문서
* superpowers skill 사이클 체험을 목적으로 함

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

Expected: 커밋 성공, 2개 파일.

- [ ] **Step 4: 커밋 검증**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git log --oneline -1
git status --short
```

Expected: 직전 커밋 표시, working tree clean.

---

## Phase 0: 프로젝트 스캐폴드

### Task 0.1: 디렉토리/Makefile/README 스켈레톤

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/README.md`
- Create: `tutorials-go/ai/superpowers/todo/Makefile`
- Create: `tutorials-go/ai/superpowers/todo/.gitignore`

- [ ] **Step 1: 디렉토리 골격 생성**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo
mkdir -p backend/server backend/todo frontend
```

Expected: 디렉토리 4개 생성됨.

- [ ] **Step 2: README 스켈레톤 작성**

`tutorials-go/ai/superpowers/todo/README.md` 내용:

```markdown
# Todo Web Application (superpowers 학습용)

Echo (Go) + React (Vite, TypeScript) 기반 in-memory Todo 앱. superpowers plugin skill 사이클을 풀로 체험하기 위한 학습 샘플.

## 실행

```bash
# Backend (포트 8080)
make dev-be

# Frontend (포트 5173, /api/* → 8080 proxy)
make dev-fe
```

브라우저에서 http://localhost:5173 접속.

## 빌드/테스트

```bash
make build           # frontend 프로덕션 빌드
make preview-fe      # 빌드 결과 :4173 프리뷰
make test            # backend + frontend 테스트
make test-be         # backend만
make test-fe         # frontend만
```

## 정책

- **데이터 영속성 없음**: 서버 재시작 시 모든 todo 손실 (in-memory).
- **시작 순서**: BE 먼저 띄운 후 FE. BE 부재 시 FE는 에러 배너 표시.
- **동시 PATCH**: last-write-wins (버전 필드 미도입).
- **타입 동기화**: `frontend/src/types.ts`는 백엔드 JSON 모델과 수기 동기화. BE 변경 시 같이 수정.

## 문서

- 설계: `docs/superpowers/specs/2026-04-30-todo-app-design.md`
- 구현 plan: `docs/superpowers/plans/2026-04-30-todo-app-plan.md`
```

- [ ] **Step 3: Makefile 작성**

`tutorials-go/ai/superpowers/todo/Makefile` 내용:

```makefile
.PHONY: dev-be dev-fe build preview-fe test test-be test-fe install clean

BE_DIR := backend
FE_DIR := frontend

install:
	cd $(FE_DIR) && npm install

dev-be:
	cd $(BE_DIR) && go run .

dev-fe:
	cd $(FE_DIR) && npm run dev

build:
	cd $(FE_DIR) && npm run build

preview-fe: build
	cd $(FE_DIR) && npm run preview

test: test-be test-fe

test-be:
	cd $(BE_DIR) && go test -race ./... && go vet ./...

test-fe:
	cd $(FE_DIR) && npm test -- --run

clean:
	rm -rf $(FE_DIR)/dist $(FE_DIR)/node_modules
```

- [ ] **Step 4: .gitignore 작성**

`tutorials-go/ai/superpowers/todo/.gitignore` 내용:

```
# frontend
frontend/node_modules/
frontend/dist/
frontend/coverage/

# misc
.DS_Store
*.log
```

- [ ] **Step 5: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/README.md \
        ai/superpowers/todo/Makefile \
        ai/superpowers/todo/.gitignore
git commit -m "$(cat <<'EOF'
[feat/todo-app] chore: 프로젝트 스켈레톤 (README, Makefile, .gitignore)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

Expected: 3개 파일 커밋.

---

### Task 0.2: 백엔드 Go 스캐폴드 + 의존성

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/backend/main.go`
- Modify: `tutorials-go/go.mod` (Echo + uuid + testify 추가)
- Modify: `tutorials-go/go.sum`

- [ ] **Step 1: main.go 스텁 작성**

`tutorials-go/ai/superpowers/todo/backend/main.go` 내용:

```go
// Package main is the entrypoint for the superpowers todo learning sample.
// It wires the Echo server with an in-memory todo store and starts listening on :8080.
package main

import (
	"log"

	"github.com/kenshin579/tutorials-go/ai/superpowers/todo/backend/server"
	"github.com/kenshin579/tutorials-go/ai/superpowers/todo/backend/todo"
)

func main() {
	store := todo.NewStore()
	srv := server.New(store)
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
```

이 파일은 아직 컴파일되지 않는다 (server, todo 패키지 미구현). Phase 1~4 완료 후 컴파일됨.

- [ ] **Step 2: 의존성 다운로드**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go get github.com/labstack/echo/v4
go get github.com/google/uuid
go get github.com/stretchr/testify
go mod tidy
```

Expected: `go.mod`에 echo/v4, uuid, testify 항목 추가됨.

- [ ] **Step 3: 의존성 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
grep -E "echo|uuid|testify" go.mod
```

Expected: 세 의존성 모두 출력.

- [ ] **Step 4: 커밋 (main.go는 아직 컴파일 안 되지만 의존성 추가는 커밋)**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/backend/main.go go.mod go.sum
git commit -m "$(cat <<'EOF'
[feat/todo-app] chore: 백엔드 main.go 스텁 + Echo/uuid/testify 의존성 추가

* main.go는 server/todo 패키지 구현 후 컴파일됨

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 0.3: 프론트엔드 Vite + React + TypeScript 스캐폴드

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/frontend/package.json`
- Create: `tutorials-go/ai/superpowers/todo/frontend/vite.config.ts`
- Create: `tutorials-go/ai/superpowers/todo/frontend/tsconfig.json`
- Create: `tutorials-go/ai/superpowers/todo/frontend/tsconfig.node.json`
- Create: `tutorials-go/ai/superpowers/todo/frontend/index.html`
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/main.tsx`
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/App.tsx`
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/index.css`
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/__tests__/setup.ts`

- [ ] **Step 1: package.json 작성**

`tutorials-go/ai/superpowers/todo/frontend/package.json` 내용:

```json
{
  "name": "superpowers-todo-frontend",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "preview": "vite preview",
    "test": "vitest"
  },
  "dependencies": {
    "react": "^19.0.0",
    "react-dom": "^19.0.0"
  },
  "devDependencies": {
    "@testing-library/react": "^16.0.0",
    "@testing-library/user-event": "^14.5.0",
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "@vitejs/plugin-react": "^4.3.0",
    "jsdom": "^24.0.0",
    "msw": "^2.4.0",
    "typescript": "^5.5.0",
    "vite": "^5.4.0",
    "vitest": "^2.1.0"
  }
}
```

- [ ] **Step 2: tsconfig 작성**

`tutorials-go/ai/superpowers/todo/frontend/tsconfig.json`:

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "useDefineForClassFields": true,
    "lib": ["ES2022", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "types": ["vitest/globals"]
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

`tutorials-go/ai/superpowers/todo/frontend/tsconfig.node.json`:

```json
{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true,
    "strict": true
  },
  "include": ["vite.config.ts"]
}
```

- [ ] **Step 3: vite.config.ts 작성 (proxy + vitest)**

`tutorials-go/ai/superpowers/todo/frontend/vite.config.ts`:

```ts
/// <reference types="vitest" />
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: false,
      },
    },
  },
  preview: { port: 4173 },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/__tests__/setup.ts'],
  },
})
```

- [ ] **Step 4: index.html 작성**

`tutorials-go/ai/superpowers/todo/frontend/index.html`:

```html
<!doctype html>
<html lang="ko">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Superpowers Todo</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

- [ ] **Step 5: 진입점 + 빈 App + CSS 작성**

`tutorials-go/ai/superpowers/todo/frontend/src/main.tsx`:

```tsx
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import './index.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
```

`tutorials-go/ai/superpowers/todo/frontend/src/App.tsx`:

```tsx
export default function App() {
  return <h1>Superpowers Todo (work in progress)</h1>
}
```

`tutorials-go/ai/superpowers/todo/frontend/src/index.css`:

```css
:root {
  font-family: system-ui, -apple-system, sans-serif;
  color-scheme: light dark;
}

body {
  margin: 0;
  padding: 1rem;
  max-width: 640px;
  margin-inline: auto;
}
```

- [ ] **Step 6: Vitest setup (빈 셋업, MSW는 Phase 5에서 추가)**

`tutorials-go/ai/superpowers/todo/frontend/src/__tests__/setup.ts`:

```ts
// Phase 5에서 MSW server 등록 추가됨.
// 현재는 testing-library 자동 cleanup만 활성화 (vitest globals + jsdom으로 충분).
```

- [ ] **Step 7: 의존성 설치**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm install
```

Expected: `node_modules/` 생성, `package-lock.json` 생성. 경고가 있어도 에러만 없으면 OK.

- [ ] **Step 8: 부팅 검증**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm run build
```

Expected: `dist/` 생성, 컴파일 에러 없음.

- [ ] **Step 9: 커밋 (`package-lock.json` 포함)**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/package.json \
        ai/superpowers/todo/frontend/package-lock.json \
        ai/superpowers/todo/frontend/vite.config.ts \
        ai/superpowers/todo/frontend/tsconfig.json \
        ai/superpowers/todo/frontend/tsconfig.node.json \
        ai/superpowers/todo/frontend/index.html \
        ai/superpowers/todo/frontend/src/main.tsx \
        ai/superpowers/todo/frontend/src/App.tsx \
        ai/superpowers/todo/frontend/src/index.css \
        ai/superpowers/todo/frontend/src/__tests__/setup.ts
git commit -m "$(cat <<'EOF'
[feat/todo-app] chore: Vite + React + TypeScript 프론트엔드 스캐폴드

* /api/* → :8080 proxy 설정
* Vitest + jsdom 테스트 환경 셋업

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Phase 1: 도메인 모델

### Task 1: `todo.go` — Priority, Todo, NewTodo, Patch, Validate, errors

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/backend/todo/todo.go`
- Create: `tutorials-go/ai/superpowers/todo/backend/todo/todo_test.go`

- [ ] **Step 1: 실패 테스트 작성 — Priority.IsValid + Validate (제목 검증 + priority 검증 + dueDate)**

`tutorials-go/ai/superpowers/todo/backend/todo/todo_test.go`:

```go
package todo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPriority_IsValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		p    Priority
		want bool
	}{
		{"low", PriorityLow, true},
		{"medium", PriorityMedium, true},
		{"high", PriorityHigh, true},
		{"empty", Priority(""), false},
		{"unknown", Priority("urgent"), false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, tc.p.IsValid())
		})
	}
}

func TestNewTodo_Validate(t *testing.T) {
	t.Parallel()
	future := time.Now().Add(24 * time.Hour)
	tooLong := ""
	for i := 0; i < 201; i++ {
		tooLong += "a"
	}

	tests := []struct {
		name      string
		input     NewTodo
		wantField string // 빈 문자열이면 통과 기대
	}{
		{"title only", NewTodo{Title: "buy milk"}, ""},
		{"title with priority", NewTodo{Title: "x", Priority: PriorityHigh}, ""},
		{"title with future due", NewTodo{Title: "x", DueDate: &future}, ""},
		{"empty title", NewTodo{Title: ""}, "title"},
		{"whitespace title", NewTodo{Title: "   "}, "title"},
		{"title length 201", NewTodo{Title: tooLong}, "title"},
		{"invalid priority", NewTodo{Title: "x", Priority: Priority("urgent")}, "priority"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.input.Validate()
			if tc.wantField == "" {
				assert.NoError(t, err)
				return
			}
			var verr *ValidationError
			if assert.ErrorAs(t, err, &verr) {
				assert.Equal(t, tc.wantField, verr.Field)
			}
		})
	}
}

func TestNewTodo_Validate_TitleLengthBoundaries(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		titleLen  int
		wantValid bool
	}{
		{"len 1", 1, true},
		{"len 200", 200, true},
		{"len 0", 0, false},
		{"len 201", 201, false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			title := ""
			for i := 0; i < tc.titleLen; i++ {
				title += "a"
			}
			err := NewTodo{Title: title}.Validate()
			if tc.wantValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test ./ai/superpowers/todo/backend/todo/...
```

Expected: 컴파일 에러 — `Priority`, `NewTodo`, `ValidationError` 등 미정의.

- [ ] **Step 3: `todo.go` 최소 구현**

`tutorials-go/ai/superpowers/todo/backend/todo/todo.go`:

```go
// Package todo defines the domain model and in-memory store for the
// superpowers todo learning sample. The store is concurrency-safe via
// sync.RWMutex, and the package exposes only value types so that
// callers cannot mutate stored state through returned references.
package todo

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Priority is the importance level of a todo item.
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// IsValid reports whether p is one of the defined Priority constants.
func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

// Todo is the domain entity. ID, CreatedAt, UpdatedAt are server-managed.
type Todo struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Completed bool       `json:"completed"`
	Priority  Priority   `json:"priority"`
	DueDate   *time.Time `json:"dueDate"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// NewTodo is the input type for Store.Add. Title is required.
// Priority defaults to PriorityMedium when empty. DueDate is optional.
type NewTodo struct {
	Title    string
	Priority Priority
	DueDate  *time.Time
}

// Validate returns a *ValidationError when input is rejected, otherwise nil.
// Title is trimmed before length check (1-200 chars).
// Priority must be empty (defaulted later) or one of the defined constants.
func (n NewTodo) Validate() error {
	title := strings.TrimSpace(n.Title)
	if title == "" {
		return &ValidationError{Field: "title", Message: "title is required"}
	}
	if len(title) > 200 {
		return &ValidationError{Field: "title", Message: "title must be at most 200 characters"}
	}
	if n.Priority != "" && !n.Priority.IsValid() {
		return &ValidationError{Field: "priority", Message: fmt.Sprintf("priority %q is invalid", n.Priority)}
	}
	return nil
}

// Patch describes a partial update to a Todo.
// A nil pointer means "field absent in request" (no change).
// ClearDueDate true means the request explicitly set dueDate to null.
type Patch struct {
	Title        *string
	Completed    *bool
	Priority     *Priority
	DueDate      *time.Time
	ClearDueDate bool
}

// Validate returns *ValidationError when patch is invalid.
// An empty patch (no fields set) is considered an error to match the API contract.
func (p Patch) Validate() error {
	if p.Title == nil && p.Completed == nil && p.Priority == nil && p.DueDate == nil && !p.ClearDueDate {
		return &ValidationError{Field: "", Message: "request body must contain at least one field"}
	}
	if p.Title != nil {
		t := strings.TrimSpace(*p.Title)
		if t == "" {
			return &ValidationError{Field: "title", Message: "title is required"}
		}
		if len(t) > 200 {
			return &ValidationError{Field: "title", Message: "title must be at most 200 characters"}
		}
	}
	if p.Priority != nil && !p.Priority.IsValid() {
		return &ValidationError{Field: "priority", Message: fmt.Sprintf("priority %q is invalid", *p.Priority)}
	}
	return nil
}

// ValidationError indicates a field-level validation failure.
// Field may be empty for whole-body errors.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements error.
func (e *ValidationError) Error() string { return e.Message }

// ErrNotFound is returned by Store.Update and Store.Delete when the id is unknown.
var ErrNotFound = errors.New("todo not found")
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test -race ./ai/superpowers/todo/backend/todo/...
```

Expected: 모든 테스트 통과 (3개 부모 + 16개 subtest).

- [ ] **Step 5: vet 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go vet ./ai/superpowers/todo/backend/todo/...
```

Expected: 출력 없음 (clean).

- [ ] **Step 6: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/backend/todo/todo.go ai/superpowers/todo/backend/todo/todo_test.go
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: 도메인 모델 (Priority, Todo, NewTodo, Patch, Validate)

* TDD: 검증 규칙을 표 기반 테스트로 먼저 작성
* title trim 후 1~200자, priority enum, ValidationError, ErrNotFound

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Phase 2: Store

### Task 2.1: `store.go` — CRUD 기본 (Add / Get / Update / Delete)

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/backend/todo/store.go`
- Create: `tutorials-go/ai/superpowers/todo/backend/todo/store_test.go`

- [ ] **Step 1: 실패 테스트 작성 — CRUD 시나리오**

`tutorials-go/ai/superpowers/todo/backend/todo/store_test.go`:

```go
package todo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return NewStore()
}

func TestStore_Add_AssignsIDAndTimestamps(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	before := time.Now()
	got := s.Add(NewTodo{Title: "buy milk"})
	after := time.Now()

	assert.NotEmpty(t, got.ID)
	assert.Equal(t, "buy milk", got.Title)
	assert.False(t, got.Completed)
	assert.Equal(t, PriorityMedium, got.Priority, "default priority")
	assert.Nil(t, got.DueDate)
	assert.False(t, got.CreatedAt.Before(before))
	assert.False(t, got.CreatedAt.After(after))
	assert.Equal(t, got.CreatedAt, got.UpdatedAt)
}

func TestStore_Add_TrimsTitle(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	got := s.Add(NewTodo{Title: "  buy milk  "})
	assert.Equal(t, "buy milk", got.Title)
}

func TestStore_Add_RespectsPriorityWhenProvided(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	got := s.Add(NewTodo{Title: "x", Priority: PriorityHigh})
	assert.Equal(t, PriorityHigh, got.Priority)
}

func TestStore_Get_ReturnsCopy(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	added := s.Add(NewTodo{Title: "x"})
	got, ok := s.Get(added.ID)
	assert.True(t, ok)
	assert.Equal(t, added, got)

	got.Title = "mutated"
	again, _ := s.Get(added.ID)
	assert.Equal(t, "x", again.Title, "store copy must not be affected by external mutation")
}

func TestStore_Get_NotFound(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	_, ok := s.Get("nope")
	assert.False(t, ok)
}

func TestStore_Update_PartialFields(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	added := s.Add(NewTodo{Title: "x"})
	time.Sleep(time.Millisecond) // updatedAt 비교를 위한 1ms 차이

	completed := true
	newTitle := "y"
	got, err := s.Update(added.ID, Patch{Title: &newTitle, Completed: &completed})
	assert.NoError(t, err)
	assert.Equal(t, "y", got.Title)
	assert.True(t, got.Completed)
	assert.True(t, got.UpdatedAt.After(added.UpdatedAt))
	assert.Equal(t, added.CreatedAt, got.CreatedAt, "createdAt must not change")
}

func TestStore_Update_ClearDueDate(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	due := time.Now().Add(time.Hour)
	added := s.Add(NewTodo{Title: "x", DueDate: &due})
	got, err := s.Update(added.ID, Patch{ClearDueDate: true})
	assert.NoError(t, err)
	assert.Nil(t, got.DueDate)
}

func TestStore_Update_NotFound(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	completed := true
	_, err := s.Update("nope", Patch{Completed: &completed})
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStore_Update_ValidationError(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	added := s.Add(NewTodo{Title: "x"})
	empty := ""
	_, err := s.Update(added.ID, Patch{Title: &empty})
	var verr *ValidationError
	assert.ErrorAs(t, err, &verr)
	assert.Equal(t, "title", verr.Field)
}

func TestStore_Delete_Success(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	added := s.Add(NewTodo{Title: "x"})
	assert.NoError(t, s.Delete(added.ID))
	_, ok := s.Get(added.ID)
	assert.False(t, ok)
}

func TestStore_Delete_NotFound(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	assert.ErrorIs(t, s.Delete("nope"), ErrNotFound)
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test ./ai/superpowers/todo/backend/todo/...
```

Expected: `Store`, `NewStore`, `Add`, `Get`, `Update`, `Delete` 미정의 컴파일 에러.

- [ ] **Step 3: `store.go` CRUD 부분 구현**

`tutorials-go/ai/superpowers/todo/backend/todo/store.go`:

```go
package todo

import (
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Store is an in-memory, concurrency-safe collection of Todo values.
// All read methods take an RLock; mutations take a Lock.
// Returned values are copies — callers cannot affect stored state by mutation.
type Store struct {
	mu    sync.RWMutex
	todos map[string]Todo
}

// NewStore creates an empty Store ready for use.
func NewStore() *Store {
	return &Store{todos: make(map[string]Todo)}
}

// Add validates input, assigns ID/timestamps, stores, and returns a copy.
// Validation errors panic instead of returning — callers should validate first.
// Add only assigns defaults; it does not call Validate. Use NewTodo.Validate() at the boundary.
func (s *Store) Add(input NewTodo) Todo {
	now := time.Now().UTC()
	priority := input.Priority
	if priority == "" {
		priority = PriorityMedium
	}
	t := Todo{
		ID:        uuid.NewString(),
		Title:     strings.TrimSpace(input.Title),
		Completed: false,
		Priority:  priority,
		DueDate:   input.DueDate,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.mu.Lock()
	s.todos[t.ID] = t
	s.mu.Unlock()
	return t
}

// Get returns the todo with the given id and true, or zero value and false.
func (s *Store) Get(id string) (Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.todos[id]
	return t, ok
}

// Update applies the patch to the todo identified by id.
// Returns ErrNotFound when id is unknown, or *ValidationError when patch is invalid.
// Validates the patch before applying.
func (s *Store) Update(id string, p Patch) (Todo, error) {
	if err := p.Validate(); err != nil {
		return Todo{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.todos[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	if p.Title != nil {
		t.Title = strings.TrimSpace(*p.Title)
	}
	if p.Completed != nil {
		t.Completed = *p.Completed
	}
	if p.Priority != nil {
		t.Priority = *p.Priority
	}
	if p.ClearDueDate {
		t.DueDate = nil
	} else if p.DueDate != nil {
		t.DueDate = p.DueDate
	}
	t.UpdatedAt = time.Now().UTC()
	s.todos[id] = t
	return t, nil
}

// Delete removes the todo by id. Returns ErrNotFound when id is unknown.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return ErrNotFound
	}
	delete(s.todos, id)
	return nil
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test -race ./ai/superpowers/todo/backend/todo/...
```

Expected: Phase 1 + Phase 2.1 모든 테스트 통과.

- [ ] **Step 5: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/backend/todo/store.go ai/superpowers/todo/backend/todo/store_test.go
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: in-memory Store CRUD (Add/Get/Update/Delete)

* sync.RWMutex로 동시성 안전, 값 복사로 외부 mutation 차단
* Update는 Patch 검증 후 부분 갱신, ClearDueDate로 null 명시 처리

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 2.2: `store.go` — List 필터/정렬 + 동시성 검증

**Files:**
- Modify: `tutorials-go/ai/superpowers/todo/backend/todo/store.go`
- Modify: `tutorials-go/ai/superpowers/todo/backend/todo/store_test.go`

- [ ] **Step 1: 실패 테스트 추가 — List 필터/정렬/동시성**

`store_test.go` 끝에 다음 코드 추가:

```go
func TestStore_List_FilterByStatus(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	a := s.Add(NewTodo{Title: "active"})
	b := s.Add(NewTodo{Title: "done"})
	completed := true
	if _, err := s.Update(b.ID, Patch{Completed: &completed}); err != nil {
		t.Fatalf("update: %v", err)
	}
	_ = a

	tests := []struct {
		status   StatusFilter
		wantLen  int
		wantTitle string
	}{
		{StatusAll, 2, ""},
		{StatusActive, 1, "active"},
		{StatusCompleted, 1, "done"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(string(tc.status), func(t *testing.T) {
			t.Parallel()
			got := s.List(Query{Status: tc.status})
			assert.Len(t, got, tc.wantLen)
			if tc.wantTitle != "" && len(got) == 1 {
				assert.Equal(t, tc.wantTitle, got[0].Title)
			}
		})
	}
}

func TestStore_List_SortByPriority(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	low := s.Add(NewTodo{Title: "low", Priority: PriorityLow})
	high := s.Add(NewTodo{Title: "high", Priority: PriorityHigh})
	mid := s.Add(NewTodo{Title: "mid", Priority: PriorityMedium})

	asc := s.List(Query{Sort: SortPriority, Order: OrderAsc})
	assert.Equal(t, []string{low.ID, mid.ID, high.ID}, ids(asc))

	desc := s.List(Query{Sort: SortPriority, Order: OrderDesc})
	assert.Equal(t, []string{high.ID, mid.ID, low.ID}, ids(desc))
}

func TestStore_List_SortByDueDate_NilLast(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	t1 := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	a := s.Add(NewTodo{Title: "a", DueDate: &t2})
	b := s.Add(NewTodo{Title: "b", DueDate: nil})
	c := s.Add(NewTodo{Title: "c", DueDate: &t1})

	asc := s.List(Query{Sort: SortDueDate, Order: OrderAsc})
	assert.Equal(t, []string{c.ID, a.ID, b.ID}, ids(asc), "nil dueDate must be last")

	desc := s.List(Query{Sort: SortDueDate, Order: OrderDesc})
	assert.Equal(t, []string{a.ID, c.ID, b.ID}, ids(desc), "nil dueDate must be last regardless of order")
}

func TestStore_List_SortByCreatedAt_DescDefault(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	first := s.Add(NewTodo{Title: "first"})
	time.Sleep(time.Millisecond)
	second := s.Add(NewTodo{Title: "second"})

	desc := s.List(Query{}) // default sort=createdAt, order=desc
	assert.Equal(t, []string{second.ID, first.ID}, ids(desc))
}

func TestStore_List_TieBreakByCreatedAtAsc(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	a := s.Add(NewTodo{Title: "a", Priority: PriorityHigh})
	time.Sleep(time.Millisecond)
	b := s.Add(NewTodo{Title: "b", Priority: PriorityHigh})

	asc := s.List(Query{Sort: SortPriority, Order: OrderAsc})
	assert.Equal(t, []string{a.ID, b.ID}, ids(asc))
}

func TestStore_ConcurrentAddList(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Add(NewTodo{Title: "t"})
		}(i)
	}
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_ = s.List(Query{})
			}
		}
	}()
	wg.Wait()
	close(done)
	assert.Len(t, s.List(Query{}), 100)
}

func ids(ts []Todo) []string {
	out := make([]string, len(ts))
	for i, t := range ts {
		out[i] = t.ID
	}
	return out
}
```

테스트 파일 상단 import에 `"sync"` 추가:

```go
import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test -race ./ai/superpowers/todo/backend/todo/...
```

Expected: `Query`, `StatusFilter`, `SortKey`, `OrderDir`, `s.List` 미정의 컴파일 에러.

- [ ] **Step 3: `store.go`에 Query 타입 + List 구현 추가**

`store.go` 끝에 다음 코드 추가:

```go
// StatusFilter narrows List results by completion state.
type StatusFilter string

const (
	StatusAll       StatusFilter = "all"
	StatusActive    StatusFilter = "active"
	StatusCompleted StatusFilter = "completed"
)

// IsValid reports whether s is one of the defined StatusFilter constants.
// Empty string is treated as StatusAll by callers and is not considered valid here.
func (s StatusFilter) IsValid() bool {
	switch s {
	case StatusAll, StatusActive, StatusCompleted:
		return true
	default:
		return false
	}
}

// SortKey selects the field used to order List results.
type SortKey string

const (
	SortCreatedAt SortKey = "createdAt"
	SortDueDate   SortKey = "dueDate"
	SortPriority  SortKey = "priority"
)

// IsValid reports whether k is one of the defined SortKey constants.
func (k SortKey) IsValid() bool {
	switch k {
	case SortCreatedAt, SortDueDate, SortPriority:
		return true
	default:
		return false
	}
}

// OrderDir is ascending or descending.
type OrderDir string

const (
	OrderAsc  OrderDir = "asc"
	OrderDesc OrderDir = "desc"
)

// IsValid reports whether o is one of the defined OrderDir constants.
func (o OrderDir) IsValid() bool {
	switch o {
	case OrderAsc, OrderDesc:
		return true
	default:
		return false
	}
}

// Query is the input to Store.List. Zero values mean defaults:
// Status defaults to StatusAll, Sort to SortCreatedAt, Order to OrderDesc.
type Query struct {
	Status StatusFilter
	Sort   SortKey
	Order  OrderDir
}

// List returns todos matching the query, sorted as requested.
// nil dueDate values are placed last regardless of Order.
// Tie-breaks fall back to createdAt ascending for determinism.
func (s *Store) List(q Query) []Todo {
	q = q.withDefaults()
	s.mu.RLock()
	out := make([]Todo, 0, len(s.todos))
	for _, t := range s.todos {
		if !matchesStatus(t, q.Status) {
			continue
		}
		out = append(out, t)
	}
	s.mu.RUnlock()
	sortTodos(out, q.Sort, q.Order)
	return out
}

func (q Query) withDefaults() Query {
	if q.Status == "" {
		q.Status = StatusAll
	}
	if q.Sort == "" {
		q.Sort = SortCreatedAt
	}
	if q.Order == "" {
		q.Order = OrderDesc
	}
	return q
}

func matchesStatus(t Todo, f StatusFilter) bool {
	switch f {
	case StatusActive:
		return !t.Completed
	case StatusCompleted:
		return t.Completed
	default:
		return true
	}
}

// sortTodos sorts ts in place by the given key/order.
// nil dueDate values always sort last regardless of order (deterministic).
// Tie-breaks fall back to createdAt ascending (stable sort).
func sortTodos(ts []Todo, key SortKey, order OrderDir) {
	asc := order == OrderAsc
	sort.SliceStable(ts, func(i, j int) bool {
		a, b := ts[i], ts[j]

		// dueDate 정렬: nil은 항상 마지막 (asc/desc 무관)
		if key == SortDueDate {
			if a.DueDate == nil && b.DueDate == nil {
				return a.CreatedAt.Before(b.CreatedAt)
			}
			if a.DueDate == nil {
				return false
			}
			if b.DueDate == nil {
				return true
			}
		}

		cmp := primaryCompare(a, b, key)
		if cmp == 0 {
			return a.CreatedAt.Before(b.CreatedAt) // tie-break: createdAt asc
		}
		if asc {
			return cmp < 0
		}
		return cmp > 0
	})
}

func primaryCompare(a, b Todo, key SortKey) int {
	switch key {
	case SortPriority:
		return priorityRank(a.Priority) - priorityRank(b.Priority)
	case SortDueDate:
		return compareTime(*a.DueDate, *b.DueDate) // 호출 시점에 둘 다 non-nil 보장됨 (sortTodos에서 분기)
	default: // SortCreatedAt
		return compareTime(a.CreatedAt, b.CreatedAt)
	}
}

func priorityRank(p Priority) int {
	switch p {
	case PriorityLow:
		return 1
	case PriorityMedium:
		return 2
	case PriorityHigh:
		return 3
	default:
		return 0
	}
}

func compareTime(a, b time.Time) int {
	if a.Before(b) {
		return -1
	}
	if a.After(b) {
		return 1
	}
	return 0
}
```

`store.go` 상단 import에 `"sort"` 추가:

```go
import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test -race ./ai/superpowers/todo/backend/todo/...
```

Expected: Phase 2.2 추가 테스트 모두 통과. `-race` 플래그로 동시성 테스트도 통과.

- [ ] **Step 5: vet 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go vet ./ai/superpowers/todo/backend/todo/...
```

Expected: 출력 없음.

- [ ] **Step 6: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/backend/todo/store.go ai/superpowers/todo/backend/todo/store_test.go
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: Store List 필터/정렬 + 동시성 안전 검증

* Query (Status/Sort/Order) 타입 추가
* dueDate nil은 항상 마지막, createdAt asc 안정 tie-break
* 100 goroutine 동시 Add+List race 테스트 통과

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Phase 3: Handler

### Task 3.1: `handler.go` — Handler 구조 + Create + Delete + writeError

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/backend/todo/handler.go`
- Create: `tutorials-go/ai/superpowers/todo/backend/todo/handler_test.go`

- [ ] **Step 1: 실패 테스트 작성 — Create + Delete + 에러 매핑**

`tutorials-go/ai/superpowers/todo/backend/todo/handler_test.go`:

```go
package todo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newTestHandler(t *testing.T) (*Handler, *Store) {
	t.Helper()
	s := NewStore()
	return NewHandler(s), s
}

func newJSONRequest(method, path, body string) (*echo.Echo, echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return e, c, rec
}

func TestHandler_Create_Returns201WithBody(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)
	_, c, rec := newJSONRequest(http.MethodPost, "/api/todos", `{"title":"buy milk"}`)
	assert.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"title":"buy milk"`)
	assert.Contains(t, rec.Body.String(), `"priority":"medium"`)
}

func TestHandler_Create_Returns400OnEmptyTitle(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)
	_, c, rec := newJSONRequest(http.MethodPost, "/api/todos", `{"title":"  "}`)
	assert.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"code":"validation_failed"`)
	assert.Contains(t, rec.Body.String(), `"field":"title"`)
}

func TestHandler_Create_Returns400OnInvalidJSON(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)
	_, c, rec := newJSONRequest(http.MethodPost, "/api/todos", `{not json`)
	assert.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"code":"invalid_json"`)
}

func TestHandler_Delete_Returns204(t *testing.T) {
	t.Parallel()
	h, s := newTestHandler(t)
	added := s.Add(NewTodo{Title: "x"})

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/todos/"+added.ID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(added.ID)

	assert.NoError(t, h.Delete(c))
	assert.Equal(t, http.StatusNoContent, rec.Code)
	_, ok := s.Get(added.ID)
	assert.False(t, ok)
}

func TestHandler_Delete_Returns404OnUnknownID(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/todos/nope", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nope")

	assert.NoError(t, h.Delete(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), `"code":"not_found"`)
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test ./ai/superpowers/todo/backend/todo/...
```

Expected: `Handler`, `NewHandler`, `Create`, `Delete` 미정의 컴파일 에러.

- [ ] **Step 3: `handler.go` 최소 구현**

`tutorials-go/ai/superpowers/todo/backend/todo/handler.go`:

```go
package todo

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler exposes the todo store as JSON HTTP endpoints.
// It is responsible for HTTP I/O concerns only — parsing, validation
// shape, status code mapping. Domain logic lives in Store.
type Handler struct {
	store *Store
}

// NewHandler returns a Handler bound to the given store.
func NewHandler(s *Store) *Handler {
	return &Handler{store: s}
}

// Create handles POST /api/todos.
// Body: { "title": "...", "priority"?: "...", "dueDate"?: "RFC3339" }.
// Returns 201 with the created Todo, or 400 on validation/JSON errors.
func (h *Handler) Create(c echo.Context) error {
	var body createBody
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return c.JSON(http.StatusBadRequest, errBody("invalid_json", "request body is not valid JSON", nil))
	}
	input := NewTodo{
		Title:    body.Title,
		Priority: Priority(body.Priority),
		DueDate:  body.DueDate,
	}
	if err := input.Validate(); err != nil {
		return writeError(c, err)
	}
	created := h.store.Add(input)
	return c.JSON(http.StatusCreated, created)
}

// Delete handles DELETE /api/todos/{id}.
// Returns 204 on success or 404 when the id is unknown.
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.store.Delete(id); err != nil {
		return writeError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

type createBody struct {
	Title    string  `json:"title"`
	Priority string  `json:"priority"`
	DueDate  *time.Time `json:"dueDate"`
}

// errBody builds the standard error envelope { "error": { code, message, details? } }.
func errBody(code, msg string, details map[string]any) echo.Map {
	body := echo.Map{"code": code, "message": msg}
	if details != nil {
		body["details"] = details
	}
	return echo.Map{"error": body}
}

// writeError maps domain errors to JSON error responses with appropriate HTTP status codes.
func writeError(c echo.Context, err error) error {
	var verr *ValidationError
	switch {
	case errors.As(err, &verr):
		details := map[string]any{}
		if verr.Field != "" {
			details["field"] = verr.Field
		}
		if len(details) == 0 {
			details = nil
		}
		return c.JSON(http.StatusBadRequest, errBody("validation_failed", verr.Message, details))
	case errors.Is(err, ErrNotFound):
		return c.JSON(http.StatusNotFound, errBody("not_found", err.Error(), nil))
	default:
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, errBody("internal_error", "internal server error", nil))
	}
}
```

위 코드의 `createBody`에 `time.Time` import 추가가 필요. 파일 상단 import 블록에 `"time"` 추가:

```go
import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test -race ./ai/superpowers/todo/backend/todo/...
```

Expected: Phase 1, 2, 3.1 테스트 모두 통과.

- [ ] **Step 5: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/backend/todo/handler.go ai/superpowers/todo/backend/todo/handler_test.go
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: Handler 구조 + Create/Delete + writeError

* writeError가 ValidationError → 400, ErrNotFound → 404, 그 외 → 500 매핑
* errBody 헬퍼로 표준 에러 envelope { error: {code, message, details?} } 생성

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 3.2: `handler.go` — List + 쿼리 파라미터 검증

**Files:**
- Modify: `tutorials-go/ai/superpowers/todo/backend/todo/handler.go`
- Modify: `tutorials-go/ai/superpowers/todo/backend/todo/handler_test.go`

- [ ] **Step 1: 실패 테스트 추가 — List 정상 + 잘못된 쿼리 파라미터**

`handler_test.go` 끝에 추가:

```go
func TestHandler_List_Defaults(t *testing.T) {
	t.Parallel()
	h, s := newTestHandler(t)
	s.Add(NewTodo{Title: "first"})
	s.Add(NewTodo{Title: "second"})

	_, c, rec := newJSONRequest(http.MethodGet, "/api/todos", "")
	assert.NoError(t, h.List(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	// JSON array of length 2
	body := rec.Body.String()
	assert.Contains(t, body, `"title":"first"`)
	assert.Contains(t, body, `"title":"second"`)
	assert.True(t, strings.HasPrefix(strings.TrimSpace(body), "["))
}

func TestHandler_List_FilterAndSort(t *testing.T) {
	t.Parallel()
	h, s := newTestHandler(t)
	a := s.Add(NewTodo{Title: "active"})
	b := s.Add(NewTodo{Title: "done"})
	completed := true
	if _, err := s.Update(b.ID, Patch{Completed: &completed}); err != nil {
		t.Fatalf("update: %v", err)
	}
	_ = a

	_, c, rec := newJSONRequest(http.MethodGet, "/api/todos?status=active", "")
	assert.NoError(t, h.List(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `"title":"active"`)
	assert.NotContains(t, body, `"title":"done"`)
}

func TestHandler_List_InvalidQueryParam(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		query string
	}{
		{"bad status", "status=invalid"},
		{"bad sort", "sort=garbage"},
		{"bad order", "order=sideways"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h, _ := newTestHandler(t)
			_, c, rec := newJSONRequest(http.MethodGet, "/api/todos?"+tc.query, "")
			assert.NoError(t, h.List(c))
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), `"code":"validation_failed"`)
		})
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test ./ai/superpowers/todo/backend/todo/...
```

Expected: `h.List` 미정의 컴파일 에러.

- [ ] **Step 3: `handler.go`에 List 추가**

`handler.go`에 다음 메서드 추가 (Create 위/아래 어디든 OK):

```go
// List handles GET /api/todos with optional status/sort/order query params.
// Empty params apply defaults (status=all, sort=createdAt, order=desc).
func (h *Handler) List(c echo.Context) error {
	q, err := parseQuery(c)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, h.store.List(q))
}

func parseQuery(c echo.Context) (Query, error) {
	q := Query{}
	if v := c.QueryParam("status"); v != "" {
		s := StatusFilter(v)
		if !s.IsValid() {
			return Query{}, &ValidationError{Field: "status", Message: "status must be one of: all, active, completed"}
		}
		q.Status = s
	}
	if v := c.QueryParam("sort"); v != "" {
		s := SortKey(v)
		if !s.IsValid() {
			return Query{}, &ValidationError{Field: "sort", Message: "sort must be one of: createdAt, dueDate, priority"}
		}
		q.Sort = s
	}
	if v := c.QueryParam("order"); v != "" {
		o := OrderDir(v)
		if !o.IsValid() {
			return Query{}, &ValidationError{Field: "order", Message: "order must be one of: asc, desc"}
		}
		q.Order = o
	}
	return q, nil
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test -race ./ai/superpowers/todo/backend/todo/...
```

Expected: 모든 테스트 통과.

- [ ] **Step 5: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/backend/todo/handler.go ai/superpowers/todo/backend/todo/handler_test.go
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: Handler.List + 쿼리 파라미터 검증

* status/sort/order는 부재 시 기본값, 잘못된 값은 400 validation_failed

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 3.3: `handler.go` — Update (PATCH 의미론)

**Files:**
- Modify: `tutorials-go/ai/superpowers/todo/backend/todo/handler.go`
- Modify: `tutorials-go/ai/superpowers/todo/backend/todo/handler_test.go`

- [ ] **Step 1: 실패 테스트 추가 — PATCH 의미론 (table-driven)**

`handler_test.go` 끝에 추가:

```go
func TestHandler_Update_PatchSemantics(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		body        string
		wantStatus  int
		wantInBody  string
		notInBody   string
		wantField   string
	}{
		{
			name:       "completed only",
			body:       `{"completed":true}`,
			wantStatus: http.StatusOK,
			wantInBody: `"completed":true`,
		},
		{
			name:       "clear duedate",
			body:       `{"dueDate":null}`,
			wantStatus: http.StatusOK,
			wantInBody: `"dueDate":null`,
		},
		{
			name:       "set duedate",
			body:       `{"dueDate":"2026-05-15T18:00:00Z"}`,
			wantStatus: http.StatusOK,
			wantInBody: `"dueDate":"2026-05-15T18:00:00Z"`,
		},
		{
			name:       "title only",
			body:       `{"title":"updated"}`,
			wantStatus: http.StatusOK,
			wantInBody: `"title":"updated"`,
		},
		{
			name:       "priority only",
			body:       `{"priority":"high"}`,
			wantStatus: http.StatusOK,
			wantInBody: `"priority":"high"`,
		},
		{
			name:       "empty body",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
		},
		{
			name:       "title null rejected",
			body:       `{"title":null}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
			wantField:  "title",
		},
		{
			name:       "completed null rejected",
			body:       `{"completed":null}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
			wantField:  "completed",
		},
		{
			name:       "priority null rejected",
			body:       `{"priority":null}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
			wantField:  "priority",
		},
		{
			name:       "invalid priority",
			body:       `{"priority":"urgent"}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
			wantField:  "priority",
		},
		{
			name:       "bad duedate string",
			body:       `{"dueDate":"not-a-date"}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
			wantField:  "dueDate",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h, s := newTestHandler(t)
			added := s.Add(NewTodo{Title: "x"})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPatch, "/api/todos/"+added.ID, strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(added.ID)

			assert.NoError(t, h.Update(c))
			assert.Equal(t, tc.wantStatus, rec.Code, "body: %s", rec.Body.String())
			assert.Contains(t, rec.Body.String(), tc.wantInBody)
			if tc.wantField != "" {
				assert.Contains(t, rec.Body.String(), `"field":"`+tc.wantField+`"`)
			}
		})
	}
}

func TestHandler_Update_NotFound(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/todos/nope", strings.NewReader(`{"completed":true}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nope")

	assert.NoError(t, h.Update(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), `"code":"not_found"`)
}

func TestHandler_Update_InvalidJSON(t *testing.T) {
	t.Parallel()
	h, s := newTestHandler(t)
	added := s.Add(NewTodo{Title: "x"})
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/todos/"+added.ID, strings.NewReader(`{not json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(added.ID)

	assert.NoError(t, h.Update(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"code":"invalid_json"`)
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test ./ai/superpowers/todo/backend/todo/...
```

Expected: `h.Update` 미정의.

- [ ] **Step 3: `handler.go`에 Update + parsePatch 추가**

먼저 `handler.go` 상단 import 블록에 `"bytes"`를 추가:

```go
import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)
```

`handler.go`에 다음 추가:

```go
// Update handles PATCH /api/todos/{id}. The body must be a non-empty JSON object.
// Only dueDate accepts JSON null (clears the value); other fields with null are rejected.
func (h *Handler) Update(c echo.Context) error {
	patch, err := parsePatch(c)
	if err != nil {
		return writeError(c, err)
	}
	id := c.Param("id")
	updated, err := h.store.Update(id, patch)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, updated)
}

func parsePatch(c echo.Context) (Patch, error) {
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(c.Request().Body).Decode(&raw); err != nil {
		return Patch{}, &jsonDecodeError{}
	}

	patch := Patch{}

	if v, ok := raw["title"]; ok {
		if isJSONNull(v) {
			return Patch{}, &ValidationError{Field: "title", Message: "title cannot be null"}
		}
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			return Patch{}, &ValidationError{Field: "title", Message: "title must be a string"}
		}
		patch.Title = &s
	}
	if v, ok := raw["completed"]; ok {
		if isJSONNull(v) {
			return Patch{}, &ValidationError{Field: "completed", Message: "completed cannot be null"}
		}
		var b bool
		if err := json.Unmarshal(v, &b); err != nil {
			return Patch{}, &ValidationError{Field: "completed", Message: "completed must be a boolean"}
		}
		patch.Completed = &b
	}
	if v, ok := raw["priority"]; ok {
		if isJSONNull(v) {
			return Patch{}, &ValidationError{Field: "priority", Message: "priority cannot be null"}
		}
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			return Patch{}, &ValidationError{Field: "priority", Message: "priority must be a string"}
		}
		p := Priority(s)
		patch.Priority = &p
	}
	if v, ok := raw["dueDate"]; ok {
		if isJSONNull(v) {
			patch.ClearDueDate = true
		} else {
			var t time.Time
			if err := json.Unmarshal(v, &t); err != nil {
				return Patch{}, &ValidationError{Field: "dueDate", Message: "dueDate must be RFC3339 timestamp or null"}
			}
			patch.DueDate = &t
		}
	}
	return patch, nil
}

// isJSONNull reports whether the raw message is the literal "null" (modulo whitespace).
// json.RawMessage may carry surrounding whitespace from the original document.
func isJSONNull(b json.RawMessage) bool {
	return string(bytes.TrimSpace(b)) == "null"
}

// jsonDecodeError wraps JSON parse failures so writeError emits invalid_json instead of validation_failed.
type jsonDecodeError struct{}

func (jsonDecodeError) Error() string { return "invalid_json" }
```

`writeError`도 jsonDecodeError를 처리하도록 수정:

```go
func writeError(c echo.Context, err error) error {
	var verr *ValidationError
	switch {
	case errors.As(err, new(*jsonDecodeError)):
		return c.JSON(http.StatusBadRequest, errBody("invalid_json", "request body is not valid JSON", nil))
	case errors.As(err, &verr):
		details := map[string]any{}
		if verr.Field != "" {
			details["field"] = verr.Field
		}
		if len(details) == 0 {
			details = nil
		}
		return c.JSON(http.StatusBadRequest, errBody("validation_failed", verr.Message, details))
	case errors.Is(err, ErrNotFound):
		return c.JSON(http.StatusNotFound, errBody("not_found", err.Error(), nil))
	default:
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, errBody("internal_error", "internal server error", nil))
	}
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test -race ./ai/superpowers/todo/backend/todo/...
```

Expected: 모든 테스트 통과.

- [ ] **Step 5: vet 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go vet ./ai/superpowers/todo/backend/todo/...
```

Expected: 출력 없음.

- [ ] **Step 6: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/backend/todo/handler.go ai/superpowers/todo/backend/todo/handler_test.go
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: Handler.Update + PATCH 의미론 (null/null-rejected/clearDueDate)

* map[string]json.RawMessage로 키 존재 여부 검출 후 필드별 unmarshal
* dueDate만 null 허용(클리어), 나머지 필드 null은 400 validation_failed

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Phase 4: Server 통합

### Task 4: `server.go` + `main.go` 통합 + httptest 통합 테스트

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/backend/server/server.go`
- Create: `tutorials-go/ai/superpowers/todo/backend/server/server_test.go`

- [ ] **Step 1: 통합 테스트 작성 (실패 예상)**

`tutorials-go/ai/superpowers/todo/backend/server/server_test.go`:

```go
package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/tutorials-go/ai/superpowers/todo/backend/server"
	"github.com/kenshin579/tutorials-go/ai/superpowers/todo/backend/todo"
)

func setupServer(t *testing.T) *httptest.Server {
	t.Helper()
	s := todo.NewStore()
	e := server.New(s)
	ts := httptest.NewServer(e)
	t.Cleanup(ts.Close)
	return ts
}

func TestServer_Health(t *testing.T) {
	t.Parallel()
	ts := setupServer(t)

	resp, err := http.Get(ts.URL + "/api/health")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]string
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "ok", body["status"])
}

func TestServer_CRUDLifecycle(t *testing.T) {
	t.Parallel()
	ts := setupServer(t)

	// 1. Create
	create := bytes.NewBufferString(`{"title":"buy milk","priority":"high"}`)
	resp, err := http.Post(ts.URL+"/api/todos", "application/json", create)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created todo.Todo
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "buy milk", created.Title)

	// 2. List
	resp, err = http.Get(ts.URL + "/api/todos")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var list []todo.Todo
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&list))
	assert.Len(t, list, 1)

	// 3. Update (toggle complete)
	patch := bytes.NewBufferString(`{"completed":true}`)
	req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/todos/"+created.ID, patch)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var updated todo.Todo
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&updated))
	assert.True(t, updated.Completed)

	// 4. Delete
	req, _ = http.NewRequest(http.MethodDelete, ts.URL+"/api/todos/"+created.ID, nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// 5. List again — should be empty
	resp, err = http.Get(ts.URL + "/api/todos")
	assert.NoError(t, err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "[]\n", string(body))
}

func TestServer_CORSPreflight(t *testing.T) {
	t.Parallel()
	ts := setupServer(t)

	req, _ := http.NewRequest(http.MethodOptions, ts.URL+"/api/todos", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Origin"), "http://localhost:5173")
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test ./ai/superpowers/todo/backend/server/...
```

Expected: `server.New` 미정의.

- [ ] **Step 3: `server.go` 구현**

`tutorials-go/ai/superpowers/todo/backend/server/server.go`:

```go
// Package server wires the Echo HTTP server with middleware and route registration.
// It is intentionally domain-agnostic — domain logic is in the todo package.
package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/kenshin579/tutorials-go/ai/superpowers/todo/backend/todo"
)

// New constructs a configured *echo.Echo bound to the given store.
// Middleware order: Recover → Logger → CORS. CORS allows the Vite dev/preview origins.
func New(store *todo.Store) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:4173"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
	}))

	h := todo.NewHandler(store)

	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	e.GET("/api/todos", h.List)
	e.POST("/api/todos", h.Create)
	e.PATCH("/api/todos/:id", h.Update)
	e.DELETE("/api/todos/:id", h.Delete)

	return e
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go test -race ./ai/superpowers/todo/backend/...
```

Expected: server + todo 패키지 모든 테스트 통과.

- [ ] **Step 5: 백엔드 부팅 검증**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/backend
go run . &
sleep 1
curl -s http://localhost:8080/api/health
curl -s -X POST -H "Content-Type: application/json" -d '{"title":"hello"}' http://localhost:8080/api/todos
curl -s http://localhost:8080/api/todos
kill %1 2>/dev/null
```

Expected: 헬스체크 `{"status":"ok"}`, POST 결과 + 목록 결과 출력.

- [ ] **Step 6: vet 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
go vet ./ai/superpowers/todo/backend/...
```

Expected: 출력 없음.

- [ ] **Step 7: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/backend/server/server.go ai/superpowers/todo/backend/server/server_test.go
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: Echo 서버 통합 (라우팅, 미들웨어, CORS) + httptest 통합 테스트

* Recover, Logger, CORS(:5173, :4173) 미들웨어
* /api/health + Todo CRUD 엔드포인트 등록
* CRUD lifecycle + CORS preflight 통합 테스트

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Phase 5: 프론트엔드 인프라

### Task 5: `types.ts` + `api.ts` + MSW 셋업

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/types.ts`
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/api.ts`
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/__tests__/mocks/handlers.ts`
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/__tests__/mocks/server.ts`
- Modify: `tutorials-go/ai/superpowers/todo/frontend/src/__tests__/setup.ts`

- [ ] **Step 1: `types.ts` — 백엔드 JSON 모델과 동기화**

`tutorials-go/ai/superpowers/todo/frontend/src/types.ts`:

```ts
// 본 파일은 backend/todo/todo.go와 수동 동기화한다.
// BE 모델 변경 시 같이 수정 (CLAUDE.md 정책).

export type Priority = 'low' | 'medium' | 'high'

export type Status = 'all' | 'active' | 'completed'

export type SortKey = 'createdAt' | 'dueDate' | 'priority'

export type Order = 'asc' | 'desc'

export interface Todo {
  id: string
  title: string
  completed: boolean
  priority: Priority
  dueDate: string | null // RFC3339 또는 null
  createdAt: string
  updatedAt: string
}

export interface NewTodo {
  title: string
  priority?: Priority
  dueDate?: string | null
}

export interface Patch {
  title?: string
  completed?: boolean
  priority?: Priority
  dueDate?: string | null // null이면 명시적 클리어
}

export interface Query {
  status: Status
  sort: SortKey
  order: Order
}
```

- [ ] **Step 2: `api.ts` — 타입 안전 fetch wrapper**

`tutorials-go/ai/superpowers/todo/frontend/src/api.ts`:

```ts
import type { NewTodo, Patch, Query, Todo } from './types'

const BASE = '/api'

export class ApiError extends Error {
  constructor(public code: string, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

async function call<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) },
  })
  if (!res.ok) {
    let code = 'unknown'
    let message = res.statusText
    try {
      const body = (await res.json()) as { error?: { code?: string; message?: string } }
      code = body.error?.code ?? code
      message = body.error?.message ?? message
    } catch {
      // body가 JSON이 아닐 수 있음. statusText 그대로 사용.
    }
    throw new ApiError(code, message)
  }
  if (res.status === 204) return undefined as T
  return (await res.json()) as T
}

export const api = {
  list: (q: Query) =>
    call<Todo[]>(`/todos?${new URLSearchParams(q as unknown as Record<string, string>)}`),
  create: (input: NewTodo) =>
    call<Todo>('/todos', { method: 'POST', body: JSON.stringify(input) }),
  update: (id: string, patch: Patch) =>
    call<Todo>(`/todos/${id}`, { method: 'PATCH', body: JSON.stringify(patch) }),
  remove: (id: string) =>
    call<void>(`/todos/${id}`, { method: 'DELETE' }),
}
```

- [ ] **Step 3: MSW 핸들러 작성 (in-memory store 흉내)**

`tutorials-go/ai/superpowers/todo/frontend/src/__tests__/mocks/handlers.ts`:

```ts
import { http, HttpResponse } from 'msw'
import type { Todo, NewTodo, Patch } from '../../types'

let todos: Todo[] = []

function uid(): string {
  return Math.random().toString(36).slice(2) + Date.now().toString(36)
}

export function resetMockTodos(seed: Todo[] = []) {
  todos = seed.map((t) => ({ ...t }))
}

export const handlers = [
  http.get('/api/health', () => HttpResponse.json({ status: 'ok' })),

  http.get('/api/todos', ({ request }) => {
    const url = new URL(request.url)
    const status = url.searchParams.get('status') ?? 'all'
    let out = [...todos]
    if (status === 'active') out = out.filter((t) => !t.completed)
    if (status === 'completed') out = out.filter((t) => t.completed)
    return HttpResponse.json(out)
  }),

  http.post('/api/todos', async ({ request }) => {
    const body = (await request.json()) as NewTodo
    if (!body.title?.trim()) {
      return HttpResponse.json(
        { error: { code: 'validation_failed', message: 'title is required', details: { field: 'title' } } },
        { status: 400 },
      )
    }
    const now = new Date().toISOString()
    const created: Todo = {
      id: uid(),
      title: body.title.trim(),
      completed: false,
      priority: body.priority ?? 'medium',
      dueDate: body.dueDate ?? null,
      createdAt: now,
      updatedAt: now,
    }
    todos.push(created)
    return HttpResponse.json(created, { status: 201 })
  }),

  http.patch('/api/todos/:id', async ({ params, request }) => {
    const id = params.id as string
    const idx = todos.findIndex((t) => t.id === id)
    if (idx === -1) {
      return HttpResponse.json(
        { error: { code: 'not_found', message: 'todo not found' } },
        { status: 404 },
      )
    }
    const patch = (await request.json()) as Patch
    const t = { ...todos[idx], ...patch, updatedAt: new Date().toISOString() }
    if ('dueDate' in patch && patch.dueDate === null) t.dueDate = null
    todos[idx] = t
    return HttpResponse.json(t)
  }),

  http.delete('/api/todos/:id', ({ params }) => {
    const id = params.id as string
    const before = todos.length
    todos = todos.filter((t) => t.id !== id)
    if (todos.length === before) {
      return HttpResponse.json(
        { error: { code: 'not_found', message: 'todo not found' } },
        { status: 404 },
      )
    }
    return new HttpResponse(null, { status: 204 })
  }),
]
```

- [ ] **Step 4: MSW server 셋업**

`tutorials-go/ai/superpowers/todo/frontend/src/__tests__/mocks/server.ts`:

```ts
import { setupServer } from 'msw/node'
import { handlers } from './handlers'

export const server = setupServer(...handlers)
```

- [ ] **Step 5: setup.ts에 MSW lifecycle 훅 등록**

`tutorials-go/ai/superpowers/todo/frontend/src/__tests__/setup.ts` 전체 교체:

```ts
import { afterAll, afterEach, beforeAll } from 'vitest'
import '@testing-library/react'
import { server } from './mocks/server'
import { resetMockTodos } from './mocks/handlers'

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => {
  server.resetHandlers()
  resetMockTodos([])
})
afterAll(() => server.close())
```

- [ ] **Step 6: 빌드 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm run build
```

Expected: 컴파일 성공.

- [ ] **Step 7: 빈 테스트 실행 (MSW 셋업 자체가 깨지지 않는지)**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run
```

Expected: "No test files found" 또는 비슷한 메시지. 에러 없이 종료.

- [ ] **Step 8: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/types.ts \
        ai/superpowers/todo/frontend/src/api.ts \
        ai/superpowers/todo/frontend/src/__tests__/mocks/handlers.ts \
        ai/superpowers/todo/frontend/src/__tests__/mocks/server.ts \
        ai/superpowers/todo/frontend/src/__tests__/setup.ts
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: types/api 클라이언트 + MSW 핸들러 (in-memory mock store)

* types.ts는 BE 모델과 수기 동기화 (CLAUDE.md 정책)
* api.ts: ApiError throw, 204 처리, query string 직렬화
* MSW handlers는 BE 동작을 흉내내어 FE 통합 테스트의 contract 역할

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Phase 6: useTodos hook

### Task 6: `useTodos.ts` — useReducer + 비동기 액션

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/hooks/useTodos.ts`
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/__tests__/useTodos.test.ts`

- [ ] **Step 1: hook 테스트 작성**

`tutorials-go/ai/superpowers/todo/frontend/src/__tests__/useTodos.test.ts`:

```ts
import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { useTodos } from '../hooks/useTodos'
import { resetMockTodos } from './mocks/handlers'
import type { Query } from '../types'

const defaultQuery: Query = { status: 'all', sort: 'createdAt', order: 'desc' }

describe('useTodos', () => {
  it('초기 fetch 후 todos 채워짐', async () => {
    resetMockTodos([
      {
        id: '1',
        title: 'seeded',
        completed: false,
        priority: 'medium',
        dueDate: null,
        createdAt: '2026-04-01T00:00:00Z',
        updatedAt: '2026-04-01T00:00:00Z',
      },
    ])
    const { result } = renderHook(() => useTodos(defaultQuery))
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(result.current.todos).toHaveLength(1)
    expect(result.current.todos[0].title).toBe('seeded')
  })

  it('create 후 list refetch', async () => {
    const { result } = renderHook(() => useTodos(defaultQuery))
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(result.current.todos).toHaveLength(0)

    await act(async () => {
      await result.current.create({ title: '새 할일' })
    })
    await waitFor(() => expect(result.current.todos).toHaveLength(1))
    expect(result.current.todos[0].title).toBe('새 할일')
  })

  it('update 후 list refetch', async () => {
    const { result } = renderHook(() => useTodos(defaultQuery))
    await waitFor(() => expect(result.current.loading).toBe(false))
    await act(async () => {
      await result.current.create({ title: 'x' })
    })
    await waitFor(() => expect(result.current.todos).toHaveLength(1))
    const id = result.current.todos[0].id

    await act(async () => {
      await result.current.update(id, { completed: true })
    })
    await waitFor(() => expect(result.current.todos[0].completed).toBe(true))
  })

  it('remove 후 목록에서 제거', async () => {
    const { result } = renderHook(() => useTodos(defaultQuery))
    await waitFor(() => expect(result.current.loading).toBe(false))
    await act(async () => {
      await result.current.create({ title: 'x' })
    })
    await waitFor(() => expect(result.current.todos).toHaveLength(1))
    const id = result.current.todos[0].id

    await act(async () => {
      await result.current.remove(id)
    })
    await waitFor(() => expect(result.current.todos).toHaveLength(0))
  })

  it('네트워크 실패 시 error 상태', async () => {
    // BE down 시뮬: 핸들러 일시 비활성화
    const { server } = await import('./mocks/server')
    const { http, HttpResponse } = await import('msw')
    server.use(http.get('/api/todos', () => HttpResponse.error()))

    const { result } = renderHook(() => useTodos(defaultQuery))
    await waitFor(() => expect(result.current.error).not.toBeNull())
    expect(result.current.loading).toBe(false)
  })
})
```

- [ ] **Step 2: 테스트 실행 (실패 예상)**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run
```

Expected: useTodos 미정의로 실패.

- [ ] **Step 3: hook 구현**

`tutorials-go/ai/superpowers/todo/frontend/src/hooks/useTodos.ts`:

```ts
import { useCallback, useEffect, useReducer, useRef } from 'react'
import { api, ApiError } from '../api'
import type { NewTodo, Patch, Query, Todo } from '../types'

type State = {
  todos: Todo[]
  loading: boolean
  error: string | null
}

type Action =
  | { type: 'fetch_start' }
  | { type: 'fetch_success'; todos: Todo[] }
  | { type: 'fetch_error'; error: string }

const initial: State = { todos: [], loading: false, error: null }

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case 'fetch_start':
      return { ...state, loading: true, error: null }
    case 'fetch_success':
      return { todos: action.todos, loading: false, error: null }
    case 'fetch_error':
      return { ...state, loading: false, error: action.error }
  }
}

export function useTodos(query: Query) {
  const [state, dispatch] = useReducer(reducer, initial)
  const queryRef = useRef(query)
  queryRef.current = query

  const refetch = useCallback(async () => {
    dispatch({ type: 'fetch_start' })
    try {
      const todos = await api.list(queryRef.current)
      dispatch({ type: 'fetch_success', todos })
    } catch (e) {
      const msg = e instanceof ApiError ? e.message : e instanceof Error ? e.message : '서버에 연결할 수 없습니다'
      dispatch({ type: 'fetch_error', error: msg })
    }
  }, [])

  useEffect(() => {
    void refetch()
  }, [query.status, query.sort, query.order, refetch])

  const create = useCallback(
    async (input: NewTodo) => {
      await api.create(input)
      await refetch()
    },
    [refetch],
  )

  const update = useCallback(
    async (id: string, patch: Patch) => {
      await api.update(id, patch)
      await refetch()
    },
    [refetch],
  )

  const remove = useCallback(
    async (id: string) => {
      await api.remove(id)
      await refetch()
    },
    [refetch],
  )

  return { ...state, create, update, remove, refetch }
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run
```

Expected: 5개 테스트 모두 통과.

- [ ] **Step 5: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/hooks/useTodos.ts \
        ai/superpowers/todo/frontend/src/__tests__/useTodos.test.ts
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: useTodos hook (useReducer + create/update/remove + 자동 refetch)

* mutation 후 항상 list 재호출 (서버 사이드 정렬/필터 정합성 유지)
* ApiError/Network 에러 모두 fetch_error로 누적

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Phase 7: 컴포넌트

### Task 7.1: `TodoForm.tsx`

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/components/TodoForm.tsx`

- [ ] **Step 1: 구현**

```tsx
import { useState, type FormEvent } from 'react'
import type { NewTodo, Priority } from '../types'

interface Props {
  onCreate: (input: NewTodo) => Promise<void>
}

export function TodoForm({ onCreate }: Props) {
  const [title, setTitle] = useState('')
  const [priority, setPriority] = useState<Priority>('medium')
  const [dueDate, setDueDate] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const submit = async (e: FormEvent) => {
    e.preventDefault()
    if (!title.trim() || submitting) return
    setSubmitting(true)
    try {
      await onCreate({
        title: title.trim(),
        priority,
        dueDate: dueDate ? new Date(dueDate).toISOString() : undefined,
      })
      setTitle('')
      setDueDate('')
      setPriority('medium')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={submit} aria-label="새 할일 추가">
      <input
        aria-label="제목"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="할 일 입력..."
        maxLength={200}
      />
      <select
        aria-label="우선순위"
        value={priority}
        onChange={(e) => setPriority(e.target.value as Priority)}
      >
        <option value="low">낮음</option>
        <option value="medium">보통</option>
        <option value="high">높음</option>
      </select>
      <input
        aria-label="마감일"
        type="datetime-local"
        value={dueDate}
        onChange={(e) => setDueDate(e.target.value)}
      />
      <button type="submit" disabled={!title.trim() || submitting}>추가</button>
    </form>
  )
}
```

- [ ] **Step 2: 빌드 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm run build
```

Expected: 컴파일 성공.

- [ ] **Step 3: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/components/TodoForm.tsx
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: TodoForm 컴포넌트 (title/priority/dueDate 입력)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 7.2: `FilterBar.tsx`

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/components/FilterBar.tsx`

- [ ] **Step 1: 구현**

```tsx
import type { Order, Query, SortKey, Status } from '../types'

interface Props {
  query: Query
  onChange: (next: Query) => void
}

export function FilterBar({ query, onChange }: Props) {
  return (
    <div role="toolbar" aria-label="필터/정렬">
      <fieldset>
        <legend>상태</legend>
        {(['all', 'active', 'completed'] as Status[]).map((s) => (
          <label key={s}>
            <input
              type="radio"
              name="status"
              value={s}
              checked={query.status === s}
              onChange={() => onChange({ ...query, status: s })}
            />
            {s === 'all' ? '전체' : s === 'active' ? '미완료' : '완료'}
          </label>
        ))}
      </fieldset>
      <label>
        정렬
        <select
          value={query.sort}
          onChange={(e) => onChange({ ...query, sort: e.target.value as SortKey })}
        >
          <option value="createdAt">생성일</option>
          <option value="dueDate">마감일</option>
          <option value="priority">우선순위</option>
        </select>
      </label>
      <button
        type="button"
        aria-label="정렬 방향 토글"
        onClick={() => onChange({ ...query, order: query.order === 'asc' ? 'desc' : 'asc' as Order })}
      >
        {query.order === 'asc' ? '↑' : '↓'}
      </button>
    </div>
  )
}
```

- [ ] **Step 2: 빌드 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm run build
```

- [ ] **Step 3: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/components/FilterBar.tsx
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: FilterBar (status 라디오 / sort 셀렉트 / order 토글)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 7.3: `TodoItem.tsx` + 컴포넌트 테스트

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/components/TodoItem.tsx`
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/__tests__/TodoItem.test.tsx`

- [ ] **Step 1: 컴포넌트 테스트 먼저**

`tutorials-go/ai/superpowers/todo/frontend/src/__tests__/TodoItem.test.tsx`:

```tsx
import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { TodoItem } from '../components/TodoItem'
import type { Todo } from '../types'

function makeTodo(overrides: Partial<Todo> = {}): Todo {
  return {
    id: 'id1',
    title: '할일 1',
    completed: false,
    priority: 'high',
    dueDate: null,
    createdAt: '2026-04-30T10:00:00Z',
    updatedAt: '2026-04-30T10:00:00Z',
    ...overrides,
  }
}

describe('TodoItem', () => {
  it('priority 뱃지 표시', () => {
    render(<TodoItem todo={makeTodo({ priority: 'high' })} onUpdate={vi.fn()} onRemove={vi.fn()} />)
    expect(screen.getByText('높음')).toBeDefined()
  })

  it('dueDate 표시 (있을 때만)', () => {
    const { rerender } = render(
      <TodoItem todo={makeTodo({ dueDate: null })} onUpdate={vi.fn()} onRemove={vi.fn()} />,
    )
    expect(screen.queryByText(/마감/)).toBeNull()

    rerender(
      <TodoItem
        todo={makeTodo({ dueDate: '2026-05-15T18:00:00Z' })}
        onUpdate={vi.fn()}
        onRemove={vi.fn()}
      />,
    )
    expect(screen.getByText(/마감/)).toBeDefined()
  })

  it('체크박스 클릭 시 onUpdate 호출', async () => {
    const onUpdate = vi.fn().mockResolvedValue(undefined)
    render(<TodoItem todo={makeTodo()} onUpdate={onUpdate} onRemove={vi.fn()} />)
    const cb = screen.getByRole('checkbox', { name: /완료/ })
    await userEvent.click(cb)
    expect(onUpdate).toHaveBeenCalledWith('id1', { completed: true })
  })

  it('삭제 버튼 클릭 시 onRemove 호출', async () => {
    const onRemove = vi.fn().mockResolvedValue(undefined)
    render(<TodoItem todo={makeTodo()} onUpdate={vi.fn()} onRemove={onRemove} />)
    await userEvent.click(screen.getByRole('button', { name: /삭제/ }))
    expect(onRemove).toHaveBeenCalledWith('id1')
  })

  it('제목 클릭 → 편집 → blur 시 onUpdate 호출', async () => {
    const onUpdate = vi.fn().mockResolvedValue(undefined)
    render(<TodoItem todo={makeTodo()} onUpdate={onUpdate} onRemove={vi.fn()} />)
    const title = screen.getByText('할일 1')
    await userEvent.click(title)
    const input = screen.getByRole('textbox', { name: /제목/ })
    await userEvent.clear(input)
    await userEvent.type(input, '수정됨')
    input.blur()
    expect(onUpdate).toHaveBeenCalledWith('id1', { title: '수정됨' })
  })
})
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run
```

Expected: TodoItem 미정의 실패.

- [ ] **Step 3: TodoItem 구현**

`tutorials-go/ai/superpowers/todo/frontend/src/components/TodoItem.tsx`:

```tsx
import { useState, type KeyboardEvent } from 'react'
import type { Patch, Priority, Todo } from '../types'

interface Props {
  todo: Todo
  onUpdate: (id: string, patch: Patch) => Promise<void>
  onRemove: (id: string) => Promise<void>
}

const priorityLabel: Record<Priority, string> = { low: '낮음', medium: '보통', high: '높음' }

export function TodoItem({ todo, onUpdate, onRemove }: Props) {
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState(todo.title)

  const commit = () => {
    setEditing(false)
    const next = draft.trim()
    if (next && next !== todo.title) {
      void onUpdate(todo.id, { title: next })
    } else {
      setDraft(todo.title)
    }
  }

  const onKey = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.currentTarget.blur()
    } else if (e.key === 'Escape') {
      setDraft(todo.title)
      setEditing(false)
    }
  }

  return (
    <li>
      <input
        type="checkbox"
        aria-label="완료"
        checked={todo.completed}
        onChange={(e) => onUpdate(todo.id, { completed: e.target.checked })}
      />
      {editing ? (
        <input
          aria-label="제목 편집"
          value={draft}
          autoFocus
          onChange={(e) => setDraft(e.target.value)}
          onBlur={commit}
          onKeyDown={onKey}
        />
      ) : (
        <span onClick={() => setEditing(true)}>{todo.title}</span>
      )}
      <span data-priority={todo.priority}>{priorityLabel[todo.priority]}</span>
      {todo.dueDate && <span>마감: {new Date(todo.dueDate).toLocaleString()}</span>}
      <button type="button" aria-label="삭제" onClick={() => onRemove(todo.id)}>×</button>
    </li>
  )
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run
```

Expected: 5개 TodoItem 테스트 통과.

- [ ] **Step 5: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/components/TodoItem.tsx \
        ai/superpowers/todo/frontend/src/__tests__/TodoItem.test.tsx
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: TodoItem (체크박스/인라인 편집/삭제) + 컴포넌트 테스트

* 제목 클릭→input 전환, blur/Enter로 저장, Escape로 취소

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 7.4: `TodoList.tsx`

**Files:**
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/components/TodoList.tsx`

- [ ] **Step 1: 구현**

```tsx
import type { Patch, Todo } from '../types'
import { TodoItem } from './TodoItem'

interface Props {
  todos: Todo[]
  onUpdate: (id: string, patch: Patch) => Promise<void>
  onRemove: (id: string) => Promise<void>
}

export function TodoList({ todos, onUpdate, onRemove }: Props) {
  if (todos.length === 0) {
    return <p>할 일이 없습니다.</p>
  }
  return (
    <ul aria-label="할 일 목록">
      {todos.map((t) => (
        <TodoItem key={t.id} todo={t} onUpdate={onUpdate} onRemove={onRemove} />
      ))}
    </ul>
  )
}
```

- [ ] **Step 2: 빌드 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm run build
```

- [ ] **Step 3: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/components/TodoList.tsx
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: TodoList (빈 상태 처리 + TodoItem 매핑)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Phase 8: App 통합

### Task 8: `App.tsx` 조립 + `App.test.tsx` 통합 시나리오

**Files:**
- Modify: `tutorials-go/ai/superpowers/todo/frontend/src/App.tsx`
- Create: `tutorials-go/ai/superpowers/todo/frontend/src/__tests__/App.test.tsx`

- [ ] **Step 1: App 통합 테스트 먼저**

`tutorials-go/ai/superpowers/todo/frontend/src/__tests__/App.test.tsx`:

```tsx
import { describe, it, expect } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from '../App'
import { server } from './mocks/server'
import { http, HttpResponse } from 'msw'

describe('App 통합 시나리오', () => {
  it('빈 상태 → 추가 → 토글 → 삭제', async () => {
    render(<App />)
    await waitFor(() => expect(screen.getByText('할 일이 없습니다.')).toBeDefined())

    // 추가
    await userEvent.type(screen.getByRole('textbox', { name: '제목' }), '우유 사기')
    await userEvent.click(screen.getByRole('button', { name: '추가' }))
    await waitFor(() => expect(screen.getByText('우유 사기')).toBeDefined())

    // 완료 토글
    await userEvent.click(screen.getByRole('checkbox', { name: /완료/ }))
    await waitFor(() => {
      const cb = screen.getByRole('checkbox', { name: /완료/ }) as HTMLInputElement
      expect(cb.checked).toBe(true)
    })

    // 삭제
    await userEvent.click(screen.getByRole('button', { name: '삭제' }))
    await waitFor(() => expect(screen.getByText('할 일이 없습니다.')).toBeDefined())
  })

  it('네트워크 실패 시 에러 배너 표시', async () => {
    server.use(http.get('/api/todos', () => HttpResponse.error()))
    render(<App />)
    await waitFor(() => expect(screen.getByRole('alert')).toBeDefined())
  })
})
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run App
```

Expected: App이 아직 placeholder라 통합 시나리오 실패.

- [ ] **Step 3: App.tsx 통합 구현**

`tutorials-go/ai/superpowers/todo/frontend/src/App.tsx`:

```tsx
import { useState } from 'react'
import { useTodos } from './hooks/useTodos'
import { TodoForm } from './components/TodoForm'
import { FilterBar } from './components/FilterBar'
import { TodoList } from './components/TodoList'
import type { Query } from './types'

const defaultQuery: Query = { status: 'all', sort: 'createdAt', order: 'desc' }

export default function App() {
  const [query, setQuery] = useState<Query>(defaultQuery)
  const { todos, error, create, update, remove } = useTodos(query)

  return (
    <main>
      <h1>Superpowers Todo</h1>
      {error && <div role="alert">에러: {error}</div>}
      <TodoForm onCreate={create} />
      <FilterBar query={query} onChange={setQuery} />
      <TodoList todos={todos} onUpdate={update} onRemove={remove} />
    </main>
  )
}
```

- [ ] **Step 4: 테스트 통과 확인 (전체)**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run
```

Expected: useTodos, TodoItem, App 테스트 모두 통과.

- [ ] **Step 5: 빌드 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm run build
```

Expected: 컴파일 성공.

- [ ] **Step 6: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/App.tsx \
        ai/superpowers/todo/frontend/src/__tests__/App.test.tsx
git commit -m "$(cat <<'EOF'
[feat/todo-app] feat: App 조립 + 통합 시나리오 테스트 (MSW 풀 라운드트립)

* 빈 상태 → 추가 → 완료 토글 → 삭제 사이클
* /api/todos GET 실패 시 alert role 배너 노출

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Phase 9: 통합 검증

### Task 9.1: README 마무리 + 수동 e2e

**Files:**
- Modify: `tutorials-go/ai/superpowers/todo/README.md`

- [ ] **Step 1: README에 동작 검증 절차 추가**

기존 README의 "## 정책" 섹션 위에 다음을 삽입:

```markdown
## 동작 검증

### Backend 단독 검증

```bash
make dev-be
# 다른 터미널에서:
curl -s http://localhost:8080/api/health
curl -s -X POST -H "Content-Type: application/json" \
  -d '{"title":"우유 사기","priority":"high"}' \
  http://localhost:8080/api/todos | tee /tmp/created.json
ID=$(cat /tmp/created.json | python3 -c "import sys,json;print(json.load(sys.stdin)['id'])")
curl -s -X PATCH -H "Content-Type: application/json" \
  -d '{"completed":true}' http://localhost:8080/api/todos/$ID
curl -s http://localhost:8080/api/todos
curl -s -X DELETE -i http://localhost:8080/api/todos/$ID | head -1
```

### Frontend + Backend 통합 검증

1. 터미널 1: `make dev-be`
2. 터미널 2: `make dev-fe`
3. 브라우저로 http://localhost:5173 접속
4. 시나리오:
   - 입력 → 추가 → 목록에 등장
   - 체크박스 토글 → 완료 상태 변화
   - FilterBar로 필터링/정렬 변경 → 즉시 반영
   - 제목 클릭 → 편집 → blur로 저장
   - 삭제 버튼 → 항목 제거
5. 백엔드 종료 후 새로고침 → "에러" 배너 표시 확인
```

- [ ] **Step 2: 전체 테스트 실행 (회귀 점검)**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo
make test
```

Expected: 백엔드 + 프론트엔드 테스트 모두 통과.

- [ ] **Step 3: 수동 e2e (사람이 직접)**

위 README "Frontend + Backend 통합 검증" 시나리오를 그대로 수행하고 결과 확인. 모두 정상 동작 시 다음 단계.

- [ ] **Step 4: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/README.md
git commit -m "$(cat <<'EOF'
[feat/todo-app] docs: 동작 검증 절차 (curl + 수동 e2e) README 추가

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Phase 10: Code Review

### Task 10: superpowers code-review skill 호출

**Files:**
- 변경 없음 (리뷰 후 발견 시 수정)

- [ ] **Step 1: `superpowers:requesting-code-review` skill 호출**

리뷰 대상: `feat/todo-app` 브랜치의 모든 커밋. 리뷰어에게 다음을 명시적으로 요청:

- 도메인 모델 검증 누락 케이스가 있는가?
- Store 동시성에 미묘한 race가 있는가? (예: Update 중 Get이 부분 갱신본을 보지 않는가?)
- handler의 PATCH null 처리가 모든 필드에서 일관되게 동작하는가?
- 프론트엔드의 mutation→refetch 패턴이 race 가능성이 있는가? (예: 빠르게 두 번 클릭)
- 에러 메시지가 사용자에게 충분히 명확한가?
- Go/TypeScript 코드 스타일이 프로젝트 컨벤션 (`tutorials-go/.claude/rules/*.md`)을 따르는가?

- [ ] **Step 2: 리뷰 코멘트 정리 후 수정 (있을 시)**

각 critical/major 코멘트에 대해 별도 fixup 커밋 생성.

- [ ] **Step 3: 최종 검증**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo
make test
```

Expected: 모든 테스트 통과.

- [ ] **Step 4: PR 생성 (사용자 컨펌 후)**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git push -u origin feat/todo-app
gh pr create --title "feat: superpowers todo app (학습 샘플)" --body "$(cat <<'EOF'
## Summary
- Echo + React + in-memory Todo 웹 애플리케이션 (학습용)
- superpowers plugin skill 사이클(brainstorm → plan → TDD → review)을 풀로 체험
- BE/FE 두 프로세스 분리, Vite proxy로 dev 통합

## Test plan
- [ ] `make test-be` 통과 (race 포함)
- [ ] `make test-fe` 통과
- [ ] BE 단독 curl CRUD 사이클 동작
- [ ] BE+FE 통합 브라우저에서 CRUD 1 사이클 동작
- [ ] BE 종료 시 FE 에러 배너 표시

## 문서
- 설계: `ai/superpowers/todo/docs/superpowers/specs/2026-04-30-todo-app-design.md`
- 구현 plan: `ai/superpowers/todo/docs/superpowers/plans/2026-04-30-todo-app-plan.md`

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)" --assignee kenshin579
```

---

## 완료 정의

- [ ] Pre-flight Task A 완료 (브랜치 + spec/plan 커밋)
- [ ] Phase 0 완료 (스캐폴드 빌드 가능)
- [ ] Phase 1~3 완료 (`go test -race ./ai/superpowers/todo/backend/todo/...` 통과)
- [ ] Phase 4 완료 (httptest 통합 테스트 통과 + curl 헬스체크 OK)
- [ ] Phase 5~8 완료 (`npm test -- --run` 통과 + `npm run build` 성공)
- [ ] Phase 9 완료 (사람이 브라우저에서 CRUD 1 사이클 직접 수행)
- [ ] Phase 10 완료 (code-review 코멘트 0건 critical, PR 생성)

각 Phase 종료 시 `make test-be` (백엔드) 또는 `make test-fe` (프론트엔드) 또는 `make test` (둘 다)로 회귀 점검 후 다음 Phase 진입.
