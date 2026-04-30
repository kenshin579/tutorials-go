# Frontend 테마 리디자인 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 현행 11줄 짜리 minimal CSS를 Minimalist+ × Indigo × Pretendard 방향으로 리디자인하고, 헤더 상태 카운트 + segmented control 필터 + dismissible 에러 배너 등 시각/기능 폴리시를 적용한다.

**Architecture:** 단일 `src/index.css`에 디자인 토큰(`:root` CSS 변수)을 모두 정의한 뒤 컴포넌트별 BEM 클래스로 스타일링. JSX는 클래스 부여 + 소량의 신규 mark-up(헤더 카운트, 에러 dismiss, segmented control radio+label) 추가. 외부 의존성은 Pretendard CDN `<link>` 1줄만 추가, 기존 21개 테스트(unit 12 + e2e 9)는 모두 호환되도록 라벨/role을 유지한다.

**Tech Stack:**
- React 19 + TypeScript + Vite (기존)
- Vanilla CSS + `:root` 변수 토큰 (외부 라이브러리 X)
- Pretendard 1.3.9 (CDN `<link>`)
- Spec 참조: `docs/superpowers/specs/2026-04-30-fe-theme-redesign-design.md`

**Scope reference**: 모든 작업은 `/Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/`. git 명령은 `tutorials-go/` 루트에서 실행.

**커밋 규칙 (`tutorials-go/CLAUDE.md` 준수):**
- 형식: `[feat/fe-theme] <type>: <간결한 한국어 설명>`
- 타입: feat, fix, docs, style, refactor, test, chore
- 모든 커밋은 `feat/fe-theme` 브랜치에 적재

---

## Pre-flight

### Task A: 피처 브랜치 생성 + spec 커밋 + .gitignore

**Files:**
- Existing: `ai/superpowers/todo/docs/superpowers/specs/2026-04-30-fe-theme-redesign-design.md` (작성 완료)
- Modify: `ai/superpowers/todo/.gitignore` (`.superpowers/` 추가)

- [ ] **Step 1: 현재 브랜치 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git status --short
git branch --show-current
```

Expected: master, untracked `ai/superpowers/todo/.superpowers/`, `ai/superpowers/todo/docs/superpowers/specs/2026-04-30-fe-theme-redesign-design.md`, modified `ai/superpowers/todo/.gitignore` (만약 미수정이면 그대로).

- [ ] **Step 2: 피처 브랜치 생성**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git checkout -b feat/fe-theme
```

Expected: `Switched to a new branch 'feat/fe-theme'`.

- [ ] **Step 3: `.gitignore`에 `.superpowers/` 추가**

`/Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/.gitignore` 파일을 읽고, 마지막 줄 다음에 다음 항목을 추가:

```
# brainstorm visual companion (mockup files, persisted but not source)
.superpowers/
```

- [ ] **Step 4: spec + .gitignore 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/docs/superpowers/specs/2026-04-30-fe-theme-redesign-design.md \
        ai/superpowers/todo/.gitignore
git commit -m "$(cat <<'EOF'
[feat/fe-theme] docs: FE 테마 리디자인 spec + brainstorm 파일 gitignore

* Minimalist+ × Indigo × Pretendard 방향성
* 디자인 토큰 + 컴포넌트별 사양 + Phase 분해

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git log --oneline -1
git status --short
```

Expected: 커밋 성공, working tree clean (단, `.superpowers/`는 ignored이므로 status에 표시 X).

---

## Phase 0: 디자인 토큰 + Pretendard 셋업

### Task 0: `index.html` CDN + `index.css` 전체 토큰/베이스 작성

**Files:**
- Modify: `ai/superpowers/todo/frontend/index.html` (Pretendard `<link>` 추가)
- Replace: `ai/superpowers/todo/frontend/src/index.css` (11줄 → 토큰 + 베이스만, 컴포넌트 CSS는 다음 phase에서 추가)

이 phase에서는 토큰과 base만 둔다. 각 컴포넌트별 CSS는 그 컴포넌트의 phase에서 추가해 commit이 보기 좋게 분리되도록 한다.

- [ ] **Step 1: `index.html`에 Pretendard CDN 추가**

기존 `ai/superpowers/todo/frontend/index.html`의 `<head>`를 다음 형태로 만든다 (`<title>` 다음 줄에 `<link>` 한 줄 추가):

```html
<!doctype html>
<html lang="ko">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Superpowers Todo</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/orioncactus/pretendard@v1.3.9/dist/web/static/pretendard.min.css">
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

- [ ] **Step 2: `src/index.css` 전체 교체 — 토큰 + base만**

`/Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend/src/index.css` 전체 내용을 다음으로 교체:

```css
:root {
  /* Surface (배경, 카드, 보더) */
  --color-bg:           #fafafa;
  --color-surface:      #ffffff;
  --color-border:       #e4e4e7;
  --color-border-soft:  #f4f4f5;

  /* Text */
  --color-text:         #18181b;
  --color-text-strong:  #0a0a0a;
  --color-text-muted:   #71717a;
  --color-text-subtle:  #a1a1aa;

  /* Accent */
  --color-accent:       #6366f1;
  --color-accent-hover: #4f46e5;
  --color-accent-soft:  #eef2ff;
  --color-accent-on-soft:#4338ca;
  --color-accent-ring:  rgba(99,102,241,0.18);

  /* Priority */
  --color-priority-high-bg:   #fef2f2;
  --color-priority-high-fg:   #b91c1c;
  --color-priority-medium-bg: var(--color-accent-soft);
  --color-priority-medium-fg: var(--color-accent-on-soft);
  --color-priority-low-bg:    var(--color-border-soft);
  --color-priority-low-fg:    var(--color-text-muted);

  /* Error */
  --color-error-bg:     #fef2f2;
  --color-error-border: #fecaca;
  --color-error-fg:     #991b1b;

  /* Typography */
  --font-sans: 'Pretendard Variable', Pretendard, system-ui, -apple-system,
               'Apple SD Gothic Neo', 'Malgun Gothic', sans-serif;
  --text-xs:   11px;
  --text-sm:   12px;
  --text-base: 14px;
  --text-lg:   18px;
  --text-xl:   22px;
  --leading-tight:  1.25;
  --leading-normal: 1.5;
  --leading-loose:  1.7;
  --tracking-tight: -0.02em;

  /* Spacing / radius / shadow */
  --space-1:  4px;
  --space-2:  8px;
  --space-3:  12px;
  --space-4:  16px;
  --space-5:  20px;
  --space-6:  24px;
  --space-8:  32px;
  --radius-sm: 5px;
  --radius-md: 8px;
  --radius-lg: 10px;
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.06);
  --shadow-md: 0 1px 3px rgba(0,0,0,0.04);

  color-scheme: light;
}

* { box-sizing: border-box; }

body {
  margin: 0;
  padding: var(--space-6);
  max-width: 580px;
  margin-inline: auto;
  background: var(--color-bg);
  color: var(--color-text);
  font-family: var(--font-sans);
  font-size: var(--text-base);
  line-height: var(--leading-normal);
}

button { font-family: inherit; cursor: pointer; }
input, select { font-family: inherit; }

:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 2px;
  border-radius: var(--radius-sm);
}

/* utility */
.visually-hidden {
  position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px;
  overflow: hidden; clip: rect(0,0,0,0); white-space: nowrap; border: 0;
}
```

- [ ] **Step 3: 빌드 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm run build 2>&1 | tail -8
```

Expected: vite build 성공. 폰트 변경되어 dev에서 글꼴이 바뀌었어야 함 (시각 변화).

- [ ] **Step 4: 기존 테스트 회귀 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run 2>&1 | tail -15
```

Expected: 12개 unit 테스트 모두 PASS. CSS 변경은 마크업/aria에 영향 없으므로 안전.

- [ ] **Step 5: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/index.html \
        ai/superpowers/todo/frontend/src/index.css
git commit -m "$(cat <<'EOF'
[feat/fe-theme] feat: 디자인 토큰(:root 변수) + Pretendard 도입

* index.html에 Pretendard 1.3.9 CDN <link> 추가
* :root에 색/타이포/스페이싱/라디우스/섀도우 토큰 정의
* base body 스타일 + :focus-visible + .visually-hidden 유틸리티

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git log --oneline -1
git status --short
```

---

## Phase 1: App + Header (상태 카운트)

### Task 1: 헤더 마크업 + 카운트 표시 + 신규 테스트

**Files:**
- Modify: `ai/superpowers/todo/frontend/src/App.tsx`
- Modify: `ai/superpowers/todo/frontend/src/__tests__/App.test.tsx` (새 테스트 추가)
- Modify: `ai/superpowers/todo/frontend/src/index.css` (`.app-*` 클래스 추가)

- [ ] **Step 1: App 통합 테스트에 카운트 검증 케이스 추가 (실패 예상)**

`ai/superpowers/todo/frontend/src/__tests__/App.test.tsx`의 마지막 `})` (describe 닫힘) 직전에 다음 테스트를 추가:

```tsx
  it('헤더에 active/completed 카운트를 표시한다', async () => {
    render(<App />)
    // 빈 상태에서는 카운트 표시 생략
    await waitFor(() => expect(screen.getByText('할 일이 없습니다.')).toBeDefined())
    expect(screen.queryByText(/진행 중/)).toBeNull()

    // 1개 추가
    await userEvent.type(screen.getByRole('textbox', { name: '제목' }), 'A')
    await userEvent.click(screen.getByRole('button', { name: '추가' }))
    await waitFor(() => expect(screen.getByText(/1개 진행 중/)).toBeDefined())

    // 1개 추가 후 완료 토글
    await userEvent.type(screen.getByRole('textbox', { name: '제목' }), 'B')
    await userEvent.click(screen.getByRole('button', { name: '추가' }))
    await waitFor(() => expect(screen.getByText(/2개 진행 중/)).toBeDefined())
    const checkboxes = screen.getAllByRole('checkbox', { name: /완료/ })
    await userEvent.click(checkboxes[0])
    await waitFor(() => expect(screen.getByText(/1개 진행 중 · 1개 완료/)).toBeDefined())
  })
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run App 2>&1 | tail -15
```

Expected: 새 테스트가 FAIL — `getByText(/진행 중/)` 못 찾음.

- [ ] **Step 3: `App.tsx` 수정 — 헤더 + 카운트 + 클래스**

`ai/superpowers/todo/frontend/src/App.tsx` 전체 교체:

```tsx
import { useMemo, useState } from 'react'
import { useTodos } from './hooks/useTodos'
import { TodoForm } from './components/TodoForm'
import { FilterBar } from './components/FilterBar'
import { TodoList } from './components/TodoList'
import type { Query } from './types'

const defaultQuery: Query = { status: 'all', sort: 'createdAt', order: 'desc' }

export default function App() {
  const [query, setQuery] = useState<Query>(defaultQuery)
  const { todos, error, create, update, remove } = useTodos(query)

  const counts = useMemo(() => {
    const active = todos.filter((t) => !t.completed).length
    const completed = todos.length - active
    return { all: todos.length, active, completed }
  }, [todos])

  const statusText = (() => {
    if (counts.all === 0) return null
    if (counts.completed === 0) return `${counts.active}개 진행 중`
    if (counts.active === 0) return `${counts.completed}개 완료`
    return `${counts.active}개 진행 중 · ${counts.completed}개 완료`
  })()

  return (
    <main className="app">
      <header className="app-header">
        <h1 className="app-title">Todo</h1>
        {statusText && <span className="app-status">{statusText}</span>}
      </header>
      {error && <div role="alert">에러: {error}</div>}
      <TodoForm onCreate={create} />
      <FilterBar query={query} onChange={setQuery} />
      <TodoList todos={todos} onUpdate={update} onRemove={remove} />
    </main>
  )
}
```

- [ ] **Step 4: `index.css`에 `.app-*` 클래스 추가**

`src/index.css` 끝에 다음 추가:

```css
/* App / Header */
.app { /* container — body가 이미 max-width 처리 */ }

.app-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: var(--space-5);
}

.app-title {
  margin: 0;
  font-size: var(--text-xl);
  font-weight: 700;
  letter-spacing: var(--tracking-tight);
  color: var(--color-text-strong);
}

.app-status {
  font-size: var(--text-sm);
  color: var(--color-text-muted);
}
```

- [ ] **Step 5: 테스트 통과 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run 2>&1 | tail -15
```

Expected: 13개 unit 테스트 PASS (기존 12 + 신규 1).

- [ ] **Step 6: 빌드**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm run build 2>&1 | tail -5
```

Expected: 빌드 성공.

- [ ] **Step 7: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/App.tsx \
        ai/superpowers/todo/frontend/src/__tests__/App.test.tsx \
        ai/superpowers/todo/frontend/src/index.css
git commit -m "$(cat <<'EOF'
[feat/fe-theme] feat: 헤더에 active/completed 상태 카운트 표시

* useMemo로 todos 분류, "N개 진행 중 · M개 완료" 텍스트 생성
* 카운트 0이면 해당 부분 생략, 전체 0이면 status 자체 비표시
* .app-header / .app-title / .app-status 스타일 추가
* App 통합 테스트에 카운트 검증 시나리오 추가

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git log --oneline -1
```

---

## Phase 2: 에러 배너 (dismissible)

### Task 2: 닫기 버튼 + 자동 reset + 신규 테스트

**Files:**
- Modify: `ai/superpowers/todo/frontend/src/App.tsx`
- Modify: `ai/superpowers/todo/frontend/src/__tests__/App.test.tsx`
- Modify: `ai/superpowers/todo/frontend/src/index.css` (`.error-banner*` 클래스 추가)

- [ ] **Step 1: 에러 배너 dismiss 테스트 추가 (실패 예상)**

`ai/superpowers/todo/frontend/src/__tests__/App.test.tsx`의 describe 끝에 추가:

```tsx
  it('에러 배너 닫기 버튼 클릭 시 배너 숨김', async () => {
    server.use(http.get('/api/todos', () => HttpResponse.error()))
    render(<App />)
    await waitFor(() => expect(screen.getByRole('alert')).toBeDefined())

    await userEvent.click(screen.getByRole('button', { name: '에러 닫기' }))
    expect(screen.queryByRole('alert')).toBeNull()
  })
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run App 2>&1 | tail -15
```

Expected: 새 테스트 FAIL — `'에러 닫기'` 버튼 없음.

- [ ] **Step 3: `App.tsx` 수정 — dismissedError state + 자동 reset + 닫기 버튼**

상단 import에 `useEffect` 추가:

```tsx
import { useEffect, useMemo, useState } from 'react'
```

`App` 함수 내부에서 `error` 사용 부분을 다음으로 교체. 기존:

```tsx
  const counts = useMemo(...)
  const statusText = (() => {...})()
```

뒤에 다음 추가:

```tsx
  const [dismissedError, setDismissedError] = useState(false)
  useEffect(() => {
    setDismissedError(false)
  }, [error])
  const showError = error !== null && !dismissedError
```

그리고 `return` 문 안의 `{error && <div role="alert">에러: {error}</div>}`를 다음으로 교체:

```tsx
      {showError && (
        <div className="error-banner" role="alert">
          <span className="error-banner__label">에러:</span>
          <span className="error-banner__message">{error}</span>
          <button
            type="button"
            className="error-banner__close"
            aria-label="에러 닫기"
            onClick={() => setDismissedError(true)}
          >
            ×
          </button>
        </div>
      )}
```

- [ ] **Step 4: `index.css`에 `.error-banner*` 추가**

`src/index.css` 끝에 추가:

```css
/* Error banner */
.error-banner {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: 10px 14px;
  margin-bottom: var(--space-4);
  background: var(--color-error-bg);
  border: 1px solid var(--color-error-border);
  border-radius: var(--radius-md);
  color: var(--color-error-fg);
  font-size: var(--text-base);
}

.error-banner__label { font-weight: 500; }
.error-banner__message { flex: 1; }

.error-banner__close {
  background: transparent;
  border: 0;
  color: inherit;
  padding: 0;
  font-size: 16px;
  line-height: 1;
}
```

- [ ] **Step 5: 테스트 통과 확인**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run 2>&1 | tail -15
```

Expected: 14개 unit 테스트 PASS (기존 13 + 신규 1).

- [ ] **Step 6: 빌드**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm run build 2>&1 | tail -5
```

- [ ] **Step 7: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/App.tsx \
        ai/superpowers/todo/frontend/src/__tests__/App.test.tsx \
        ai/superpowers/todo/frontend/src/index.css
git commit -m "$(cat <<'EOF'
[feat/fe-theme] feat: 에러 배너에 닫기 버튼 + 자동 reset 추가

* dismissedError 로컬 state, useEffect([error])로 새 에러 발생 시 reset
* aria-label="에러 닫기" 버튼, 빨강 톤 배너 스타일
* App 테스트에 dismiss 시나리오 추가

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git log --oneline -1
```

---

## Phase 3: TodoForm

### Task 3: 폼 카드 + inline 입력 스타일

**Files:**
- Modify: `ai/superpowers/todo/frontend/src/components/TodoForm.tsx` (className 부여)
- Modify: `ai/superpowers/todo/frontend/src/index.css` (`.todo-form*` 클래스 추가)

- [ ] **Step 1: `TodoForm.tsx` 수정 — className만 부여, 동작은 동일**

`ai/superpowers/todo/frontend/src/components/TodoForm.tsx` 전체 교체:

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
    <form className="todo-form" onSubmit={submit} aria-label="새 할일 추가">
      <input
        className="todo-form__input"
        aria-label="제목"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="할 일을 입력하세요…"
        maxLength={200}
      />
      <select
        className="todo-form__priority"
        aria-label="우선순위"
        value={priority}
        onChange={(e) => setPriority(e.target.value as Priority)}
      >
        <option value="low">낮음</option>
        <option value="medium">보통</option>
        <option value="high">높음</option>
      </select>
      <input
        className="todo-form__due"
        aria-label="마감일"
        type="datetime-local"
        value={dueDate}
        onChange={(e) => setDueDate(e.target.value)}
      />
      <button
        className="todo-form__submit"
        type="submit"
        disabled={!title.trim() || submitting}
      >
        추가
      </button>
    </form>
  )
}
```

- [ ] **Step 2: `index.css`에 `.todo-form*` 추가**

```css
/* TodoForm */
.todo-form {
  display: flex;
  gap: var(--space-2);
  margin-bottom: var(--space-4);
  padding: var(--space-3);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
}

.todo-form__input,
.todo-form__priority,
.todo-form__due {
  padding: 9px 12px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-base);
  background: var(--color-surface);
  color: var(--color-text);
  outline: none;
}

.todo-form__input { flex: 1; min-width: 0; }

.todo-form__input:focus,
.todo-form__priority:focus,
.todo-form__due:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px var(--color-accent-ring);
}

.todo-form__submit {
  background: var(--color-accent);
  color: white;
  border: 0;
  border-radius: var(--radius-md);
  padding: 9px var(--space-4);
  font-size: var(--text-base);
  font-weight: 500;
  transition: background-color 0.15s;
}

.todo-form__submit:hover:not(:disabled) { background: var(--color-accent-hover); }
.todo-form__submit:disabled { background: var(--color-text-subtle); cursor: not-allowed; }
```

- [ ] **Step 3: 테스트 + 빌드 통과 확인 (회귀)**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run 2>&1 | tail -10
npm run build 2>&1 | tail -5
```

Expected: 14개 테스트 PASS, 빌드 성공.

- [ ] **Step 4: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/components/TodoForm.tsx \
        ai/superpowers/todo/frontend/src/index.css
git commit -m "$(cat <<'EOF'
[feat/fe-theme] feat: TodoForm을 흰색 카드 + inline 4-요소 레이아웃으로 스타일

* todo-form 카드 컨테이너, todo-form__input/priority/due/submit 클래스
* focus 시 accent ring, submit 버튼 hover/disabled 스타일

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git log --oneline -1
```

---

## Phase 4: FilterBar (segmented control + 카운트)

### Task 4: hidden radio + label 패턴 + count prop

**Files:**
- Modify: `ai/superpowers/todo/frontend/src/components/FilterBar.tsx` (mark-up 변경 + counts prop 추가)
- Modify: `ai/superpowers/todo/frontend/src/App.tsx` (counts prop 전달)
- Modify: `ai/superpowers/todo/frontend/src/index.css` (`.filter-bar*` 클래스 추가)

- [ ] **Step 1: `FilterBar.tsx` 전체 교체 — segmented control + counts**

`ai/superpowers/todo/frontend/src/components/FilterBar.tsx`:

```tsx
import type { Order, Query, SortKey, Status } from '../types'

interface Props {
  query: Query
  counts: { all: number; active: number; completed: number }
  onChange: (next: Query) => void
}

const segments: Array<{ value: Status; label: string }> = [
  { value: 'all', label: '전체' },
  { value: 'active', label: '미완료' },
  { value: 'completed', label: '완료' },
]

export function FilterBar({ query, counts, onChange }: Props) {
  return (
    <div className="filter-bar" role="toolbar" aria-label="필터/정렬">
      <fieldset className="filter-bar__segments">
        <legend className="visually-hidden">상태</legend>
        {segments.map((s) => (
          <label key={s.value} className="filter-bar__segment">
            <input
              type="radio"
              name="status"
              value={s.value}
              checked={query.status === s.value}
              onChange={() => onChange({ ...query, status: s.value })}
            />
            <span className="filter-bar__segment-label">
              {s.label}
              <span className="filter-bar__count">{counts[s.value]}</span>
            </span>
          </label>
        ))}
      </fieldset>

      <div className="filter-bar__sort">
        <select
          className="filter-bar__sort-select"
          aria-label="정렬"
          value={query.sort}
          onChange={(e) => onChange({ ...query, sort: e.target.value as SortKey })}
        >
          <option value="createdAt">생성일</option>
          <option value="dueDate">마감일</option>
          <option value="priority">우선순위</option>
        </select>
        <button
          type="button"
          className="filter-bar__order"
          aria-label="정렬 방향 토글"
          onClick={() =>
            onChange({ ...query, order: (query.order === 'asc' ? 'desc' : 'asc') as Order })
          }
        >
          {query.order === 'asc' ? '↑' : '↓'}
        </button>
      </div>
    </div>
  )
}
```

NOTE: `aria-label`은 input에 바로 붙이지 않고 `<label>`로 wrap한다. 그래서 e2e/unit 테스트의 `getByLabel('전체')` 등이 정상 매칭된다. `getByLabel`은 label 안의 input을 찾을 때 visible text를 사용하므로 `<span>` 내부 "전체" 텍스트가 매칭된다. 단 카운트 숫자(예: "전체 3")가 함께 있어 일부 테스트의 `getByLabel('전체')`가 부분 매칭으로 실패할 수 있다 — 이 점은 Phase 6 회귀 검증에서 직접 다룬다.

- [ ] **Step 2: `App.tsx`에 counts prop 전달**

기존 `<FilterBar query={query} onChange={setQuery} />`를 다음으로 교체:

```tsx
      <FilterBar query={query} counts={counts} onChange={setQuery} />
```

(Phase 1에서 이미 `counts` 변수가 useMemo로 선언되어 있음.)

- [ ] **Step 3: `index.css`에 `.filter-bar*` 추가**

```css
/* FilterBar */
.filter-bar {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

.filter-bar__segments {
  display: flex;
  background: var(--color-border-soft);
  padding: 3px;
  border-radius: var(--radius-md);
  border: 0;
  margin: 0;
}

.filter-bar__segment {
  position: relative;
  cursor: pointer;
}

.filter-bar__segment input[type="radio"] {
  position: absolute;
  opacity: 0;
  pointer-events: none;
}

.filter-bar__segment-label {
  display: inline-block;
  padding: 6px 14px;
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-muted);
  border-radius: 6px;
}

.filter-bar__segment input:checked + .filter-bar__segment-label {
  background: var(--color-surface);
  color: var(--color-text-strong);
  box-shadow: var(--shadow-sm);
}

.filter-bar__segment input:focus-visible + .filter-bar__segment-label {
  outline: 2px solid var(--color-accent);
  outline-offset: 1px;
}

.filter-bar__count {
  margin-left: 4px;
  font-size: 11px;
  opacity: 0.7;
}

.filter-bar__sort {
  margin-left: auto;
  display: flex;
  gap: var(--space-1);
  align-items: center;
}

.filter-bar__sort-select {
  padding: 5px 8px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-size: var(--text-sm);
  background: var(--color-surface);
  color: var(--color-text);
}

.filter-bar__order {
  border: 1px solid var(--color-border);
  background: var(--color-surface);
  color: var(--color-text-muted);
  padding: 5px 9px;
  border-radius: 6px;
  font-size: var(--text-sm);
}

.filter-bar__order:hover { background: var(--color-border-soft); }
```

- [ ] **Step 4: 회귀 테스트**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run 2>&1 | tail -25
```

Expected behavior + 가능한 회귀:

- 14개 테스트 중 일부가 fail할 수 있음 — 특히 `getByLabel('전체')` 류. 이유: label의 visible text가 "전체 0" (count 포함)이 되어 정확 매칭이 안 됨.

만약 fail이 발생하면, 다음 fixup 적용:

테스트 파일 `useTodos.test.ts`, `App.test.tsx`에서 `getByLabel('전체')` 등이 있으면 `getByRole('radio', { name: /전체/ })` 형태로 부분 매칭 사용. 단 e2e (playwright `getByLabel('전체').check()`)도 같은 이슈 — `getByRole('radio', { name: /전체/ }).check()` 또는 클릭 위치를 segment-label로 변경.

대안 (더 안전): segmented control의 `<span class="filter-bar__count">`를 별도 요소로 두지 말고 `aria-hidden="true"` 부여하여 접근성 트리에서 제외시키면 `getByLabel('전체')`가 다시 정확 매칭됨. JSX에서 `<span className="filter-bar__count" aria-hidden="true">{counts[s.value]}</span>`로 수정. Phase 5/6에서 검증 후 채택.

**일단 `aria-hidden="true"`를 카운트 span에 추가**하여 회귀 가능성을 줄인다. Step 1 JSX에서 다음 라인을 수정:

```tsx
<span className="filter-bar__count" aria-hidden="true">{counts[s.value]}</span>
```

수정 후 다시 테스트 실행:

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run 2>&1 | tail -15
```

Expected: 14개 테스트 모두 PASS.

- [ ] **Step 5: 빌드**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm run build 2>&1 | tail -5
```

- [ ] **Step 6: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/components/FilterBar.tsx \
        ai/superpowers/todo/frontend/src/App.tsx \
        ai/superpowers/todo/frontend/src/index.css
git commit -m "$(cat <<'EOF'
[feat/fe-theme] feat: FilterBar segmented control + 카운트 뱃지

* hidden radio + label 패턴으로 시각은 segmented, 마크업은 라디오 (a11y 유지)
* counts prop으로 전체/미완료/완료 N개 inline 표시 (count는 aria-hidden)
* 정렬 셀렉트 + 방향 토글 그룹을 우측으로 정렬

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git log --oneline -1
```

---

## Phase 5: TodoList + TodoItem

### Task 5: 카드 컨테이너, priority 색 뱃지, 완료 line-through, hover, 마감일 subtext

**Files:**
- Modify: `ai/superpowers/todo/frontend/src/components/TodoList.tsx` (className만)
- Modify: `ai/superpowers/todo/frontend/src/components/TodoItem.tsx` (className + 조건부 modifier + 마감일 subtext 분리)
- Modify: `ai/superpowers/todo/frontend/src/index.css` (`.todo-list*`, `.todo-item*` 추가)

- [ ] **Step 1: `TodoList.tsx` 수정 — className만**

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
    return <p className="todo-list todo-list--empty">할 일이 없습니다.</p>
  }
  return (
    <ul className="todo-list" aria-label="할 일 목록">
      {todos.map((t) => (
        <TodoItem key={t.id} todo={t} onUpdate={onUpdate} onRemove={onRemove} />
      ))}
    </ul>
  )
}
```

- [ ] **Step 2: `TodoItem.tsx` 수정 — className + modifier + 마감일 분리**

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

  const itemClass = ['todo-item', todo.completed ? 'todo-item--completed' : '']
    .filter(Boolean)
    .join(' ')

  return (
    <li className={itemClass}>
      <input
        type="checkbox"
        className="todo-item__checkbox"
        aria-label="완료"
        checked={todo.completed}
        onChange={(e) => onUpdate(todo.id, { completed: e.target.checked })}
      />
      <div className="todo-item__main">
        {editing ? (
          <input
            className="todo-item__title-edit"
            aria-label="제목 편집"
            value={draft}
            autoFocus
            onChange={(e) => setDraft(e.target.value)}
            onBlur={commit}
            onKeyDown={onKey}
          />
        ) : (
          <span className="todo-item__title" onClick={() => setEditing(true)}>
            {todo.title}
          </span>
        )}
        {todo.dueDate && (
          <span className="todo-item__due">
            마감 {new Date(todo.dueDate).toLocaleString()}
          </span>
        )}
      </div>
      <span
        className={`todo-item__priority todo-item__priority--${todo.priority}`}
        data-priority={todo.priority}
      >
        {priorityLabel[todo.priority]}
      </span>
      <button
        type="button"
        className="todo-item__delete"
        aria-label="삭제"
        onClick={() => onRemove(todo.id)}
      >
        ×
      </button>
    </li>
  )
}
```

NOTE: 기존 테스트 `TodoItem.test.tsx`의 `screen.queryByText(/마감/)` 기대 — 마감일 표시 텍스트가 "마감 ..." 형식이라 `/마감/` regex 매칭됨. 안전.

- [ ] **Step 3: `index.css`에 `.todo-list*`, `.todo-item*` 추가**

```css
/* TodoList */
.todo-list {
  list-style: none;
  margin: 0;
  padding: 0;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.todo-list--empty {
  text-align: center;
  padding: var(--space-8) var(--space-4);
  border-style: dashed;
  color: var(--color-text-subtle);
  font-size: var(--text-base);
}

/* TodoItem */
.todo-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border-soft);
}

.todo-item:last-child { border-bottom: 0; }

.todo-item:hover { background: var(--color-border-soft); }

.todo-item__checkbox {
  width: 16px;
  height: 16px;
  accent-color: var(--color-accent);
  flex-shrink: 0;
  cursor: pointer;
}

.todo-item__main {
  flex: 1;
  min-width: 0;
}

.todo-item__title {
  display: block;
  font-size: var(--text-base);
  color: var(--color-text);
  cursor: text;
  word-break: break-word;
}

.todo-item--completed .todo-item__title {
  color: var(--color-text-subtle);
  text-decoration: line-through;
}

.todo-item__title-edit {
  width: 100%;
  border: 1px solid var(--color-accent);
  border-radius: var(--radius-sm);
  padding: 4px 6px;
  font-size: var(--text-base);
  outline: none;
  box-shadow: 0 0 0 3px var(--color-accent-ring);
}

.todo-item__due {
  display: block;
  font-size: 11px;
  color: var(--color-text-subtle);
  margin-top: 2px;
}

.todo-item__priority {
  font-size: var(--text-xs);
  padding: 3px 8px;
  border-radius: var(--radius-sm);
  font-weight: 500;
  flex-shrink: 0;
}

.todo-item__priority--high   { background: var(--color-priority-high-bg);   color: var(--color-priority-high-fg);   }
.todo-item__priority--medium { background: var(--color-priority-medium-bg); color: var(--color-priority-medium-fg); }
.todo-item__priority--low    { background: var(--color-priority-low-bg);    color: var(--color-priority-low-fg);    }

.todo-item__delete {
  background: transparent;
  border: 0;
  color: var(--color-text-subtle);
  padding: 4px 6px;
  font-size: 16px;
  line-height: 1;
}

.todo-item__delete:hover { color: var(--color-priority-high-fg); }
```

- [ ] **Step 4: 회귀 테스트 + 빌드**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/frontend
npm test -- --run 2>&1 | tail -20
npm run build 2>&1 | tail -5
```

Expected: 14개 테스트 모두 PASS, 빌드 성공.

- [ ] **Step 5: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/src/components/TodoList.tsx \
        ai/superpowers/todo/frontend/src/components/TodoItem.tsx \
        ai/superpowers/todo/frontend/src/index.css
git commit -m "$(cat <<'EOF'
[feat/fe-theme] feat: TodoList 카드 + TodoItem 스타일 (priority/완료/hover)

* 카드 컨테이너 + 항목 구분선
* todo-item--completed 시 line-through + 회색
* priority 색 뱃지 (high=빨강 / medium=인디고 / low=그레이)
* 마감일을 제목 아래 작은 글씨 subtext로 분리, 행 hover 효과

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git log --oneline -1
```

---

## Phase 6: 회귀 검증 (e2e 포함)

### Task 6: 전체 unit + e2e + 수동 시각 검증

**Files:** 변경 없음 (검증만).

- [ ] **Step 1: 전체 unit 테스트**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo
make test-fe 2>&1 | tail -15
```

Expected: 14개 테스트 모두 PASS.

- [ ] **Step 2: 전체 e2e 테스트**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo
make test-e2e 2>&1 | tail -20
```

Expected: 9개 e2e 시나리오 모두 PASS. 만약 fail이 있으면, 가장 흔한 원인:

- `getByLabel('전체')` / `getByLabel('미완료')` / `getByLabel('완료')`가 실패 → Phase 4 Step 4의 `aria-hidden="true"` 적용이 안 된 경우. 적용 확인.
- 스크롤/타이밍 이슈 (segmented control이 더 작거나 위치가 바뀌어서) → spec 셀렉터에 `.click()` 대신 `.check()` 또는 `getByRole('radio')`로 우회. 단 우선 `aria-hidden` 적용으로 해결되는지 확인.

만약 fail이 발생해 spec 수정이 필요하면, 별도 commit으로 e2e spec만 패치 (테마 동작이 아닌 selector 변경).

- [ ] **Step 3: backend 회귀 (안전 보장)**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo
make test-be 2>&1 | tail -10
```

Expected: 백엔드 변경이 전혀 없으니 그대로 통과. 단 회귀 차원에서 한 번 실행.

- [ ] **Step 4: 수동 시각 검증 (사람이 직접)**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo
make dev-be &
make dev-fe
```

브라우저에서 http://localhost:5173 접속 후 확인:

1. Pretendard 폰트 적용됐는지 (한글 톤이 OS 기본 폰트와 다름)
2. 헤더에 Todo + 우측 카운트
3. 폼이 흰색 카드 안에 inline
4. 필터바가 segmented control + 카운트 뱃지
5. 항목이 카드 컨테이너 + 구분선
6. priority 색 (high=빨강 soft, medium=인디고 soft, low=그레이)
7. 완료 항목 line-through + 회색
8. 마감일이 작은 글씨 subtext로
9. 에러 배너에 닫기 버튼
10. 빈 상태 dashed border 박스

스크린샷 1장 캡처 후 PR 코멘트에 첨부 권장.

- [ ] **Step 5: 검증 결과 OK면 별도 commit 없이 다음 phase로**

만약 e2e fail로 spec 수정한 경우:

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/frontend/e2e/todo.spec.ts
git commit -m "$(cat <<'EOF'
[feat/fe-theme] test: segmented control 변경에 맞춰 e2e 셀렉터 조정

* getByLabel('전체') 등을 getByRole('radio', { name: /전체/ })로 변경
* 동작 변경 없음, selector만 갱신

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

만약 수정 불필요하면 이 step은 skip.

---

## Phase 7: README 테마 섹션

### Task 7: README에 테마/디자인 토큰 설명 추가

**Files:**
- Modify: `ai/superpowers/todo/README.md`

- [ ] **Step 1: README에 "테마" 섹션 추가**

기존 README의 "## 정책" 섹션 위에 새 섹션 INSERT.

찾을 문자열:

```
## 정책
```

다음으로 교체:

```
## 테마

Minimalist+ × Indigo × Pretendard 방향. 디자인 토큰은 `frontend/src/index.css`의 `:root`에 모두 정의되어 있어, 색/타이포/스페이싱을 한 곳에서 일괄 변경할 수 있다.

### 핵심 토큰

| 영역 | 변수 | 기본값 |
|---|---|---|
| 강조색 | `--color-accent` | `#6366f1` (indigo-500) |
| 본문색 | `--color-text` | `#18181b` |
| 배경 | `--color-bg` | `#fafafa` |
| 폰트 | `--font-sans` | `Pretendard Variable` (CDN) → fallback system-ui |
| 라디우스 | `--radius-md` / `--radius-lg` | `8px` / `10px` |
| priority high | `--color-priority-high-bg/fg` | `#fef2f2` / `#b91c1c` |
| priority medium | `--color-priority-medium-bg/fg` | `#eef2ff` / `#4338ca` |
| priority low | `--color-priority-low-bg/fg` | `#f4f4f5` / `#71717a` |

### 컴포넌트 클래스 (BEM 변형)

- `.app-header` — 제목 + 상태 카운트
- `.todo-form` — 흰 카드 안에 input/select/button inline
- `.filter-bar__segments` — pill segmented control (hidden radio + label)
- `.todo-list` — 카드 컨테이너, `.todo-list--empty`는 dashed border
- `.todo-item--completed` — line-through 적용
- `.todo-item__priority--{low,medium,high}` — priority별 색
- `.error-banner` — dismissible 빨강 배너

### 다크모드

`color-scheme: light`만 선언되어 있고 다크 토큰은 미정의. 추후 별도 작업.

## 정책
```

NOTE: 위 inserted block의 triple-backtick 펜스는 markdown 표 외부의 시각 표현 — 그대로 LITERAL 작성됨.

- [ ] **Step 2: 인코딩 확인**

```bash
file -I /Users/user/src/workspace_blog3/tutorials-go/ai/superpowers/todo/README.md
```

Expected: charset=utf-8.

- [ ] **Step 3: 커밋**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git add ai/superpowers/todo/README.md
git commit -m "$(cat <<'EOF'
[feat/fe-theme] docs: README에 테마/디자인 토큰 섹션 추가

* 핵심 토큰 표 + 컴포넌트 BEM 클래스 목록
* 다크모드는 미적용 상태 명시 (향후 작업)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git log --oneline -1
git status --short
```

Expected: 커밋 성공, working tree clean.

---

## Phase 8: PR 생성

### Task 8: 브랜치 push + PR 생성

**Files:** 변경 없음 (git/gh 작업).

- [ ] **Step 1: 브랜치 push (사용자 컨펌 후)**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
git push -u origin feat/fe-theme
```

Expected: 새 브랜치 push, PR 생성 URL 출력.

- [ ] **Step 2: PR 생성**

```bash
cd /Users/user/src/workspace_blog3/tutorials-go
gh pr create --title "feat: FE 테마 리디자인 (Minimalist+ × Indigo × Pretendard)" --body "$(cat <<'EOF'
## Summary

`ai/superpowers/todo/frontend/`의 시각 테마를 리디자인합니다. 11줄짜리 minimal CSS를 디자인 토큰 기반 시스템으로 확장하고, 헤더 상태 카운트 / segmented control 필터 / dismissible 에러 배너 등 폴리시를 추가합니다.

## 주요 변경

- **디자인 토큰**: `:root` CSS 변수에 색/타이포/스페이싱/라디우스 모두 정의 (`index.css`)
- **폰트**: Pretendard 1.3.9 CDN 도입 (한국어 가독성 ↑)
- **헤더**: "N개 진행 중 · M개 완료" 상태 카운트
- **TodoForm**: 흰색 카드 + inline 4-요소 (input/priority/due/추가)
- **FilterBar**: segmented control (hidden radio + label) + 카운트 뱃지 + 정렬 셀렉트/방향 토글
- **TodoList/Item**: 카드 컨테이너, 구분선, priority 색 뱃지(high=빨강 / medium=인디고 / low=그레이), 완료 항목 line-through, 마감일 작은 글씨 subtext, 행 hover
- **에러 배너**: 빨강 톤 dismissible 박스, 새 에러 발생 시 자동 reset
- **빈 상태**: dashed border 박스

## 의도적 제외 (YAGNI)

- 다크모드 (별도 작업)
- 외부 컴포넌트 라이브러리, CSS Modules
- 애니메이션 (CSS hover transition 0.15s 외)
- 반응형 모바일 레이아웃

## 회귀 영향

- 기존 12개 unit + 9개 e2e 테스트는 마크업의 aria-label/role 모두 유지하여 호환
- 신규 unit 테스트 2개 추가 (헤더 카운트 / 에러 dismiss)

## Test plan

- [ ] `make test-fe` — 14개 unit (기존 12 + 신규 2) PASS
- [ ] `make test-e2e` — 9개 e2e PASS
- [ ] `make test-be` — 백엔드 회귀 OK
- [ ] 브라우저로 http://localhost:5173 접속, 시각 검증 (Pretendard 적용 / 카드 / segmented / priority 색 / 완료 line-through / 에러 배너 dismiss)

## 문서

- 설계: `ai/superpowers/todo/docs/superpowers/specs/2026-04-30-fe-theme-redesign-design.md`
- 구현 plan: `ai/superpowers/todo/docs/superpowers/plans/2026-04-30-fe-theme-redesign-plan.md`

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)" --assignee kenshin579
```

Expected: PR URL 출력. 사용자에게 전달.

---

## 완료 정의

- [ ] Pre-flight Task A 완료 (브랜치 + spec/.gitignore 커밋)
- [ ] Phase 0 완료 (토큰 + Pretendard, 회귀 OK)
- [ ] Phase 1 완료 (헤더 + 카운트, 신규 테스트 PASS)
- [ ] Phase 2 완료 (에러 dismiss, 신규 테스트 PASS)
- [ ] Phase 3 완료 (TodoForm 카드)
- [ ] Phase 4 완료 (FilterBar segmented + count)
- [ ] Phase 5 완료 (TodoList/Item)
- [ ] Phase 6 완료 (전체 회귀 + 수동 시각 검증)
- [ ] Phase 7 완료 (README 테마 섹션)
- [ ] Phase 8 완료 (PR 생성, kenshin579 assignee)

각 phase 종료 시 `make test-fe` 또는 `make test` 회귀 점검 후 진입.
