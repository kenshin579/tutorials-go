# Claude Code: Skills vs Commands vs Subagents 구현 문서

## 개요

이 프로젝트의 `.claude/` 디렉토리에 Commands 3개, Skills 4개, Subagents 4개 (총 11개 파일)를 생성한다.

---

## 1. 디렉토리 구조

```
.claude/
├── commands/
│   ├── commit.md              # A-1
│   ├── plan-task.md           # A-2
│   └── start-task.md          # A-3
├── skills/
│   ├── go-convention/
│   │   └── SKILL.md           # B-1
│   ├── analyze-codebase/
│   │   └── SKILL.md           # B-2
│   ├── go-project-layout/
│   │   └── SKILL.md           # B-3
│   └── api-convention/
│   │   └── SKILL.md           # D-1 (지식)
└── agents/
    ├── code-reviewer.md       # C-1
    ├── debugger.md            # C-2
    ├── test-runner.md         # C-3
    └── api-developer.md       # D-1 (실행자)
```

---

## 2. Commands 구현

### A-1: `.claude/commands/commit.md`

- 프롬프트 템플릿 형식
- `git status` → 브랜치 확인 → staging → commit → push 절차
- conventional commits 형식 (`feat:`, `fix:`, `docs:` 등)
- main/master 직접 commit 금지 규칙 포함

### A-2: `.claude/commands/plan-task.md`

- `$ARGUMENTS`로 PRD 파일 경로를 인자로 받음
- PRD 분석 → `{순서}_{기능명}_implementation.md` + `{순서}_{기능명}_todo.md` 생성
- 요구사항 파일과 같은 디렉토리에 출력

### A-3: `.claude/commands/start-task.md`

- `$ARGUMENTS`로 PRD 문서 경로를 인자로 받음
- GitHub Issue 생성 (GitHub MCP 사용) → master에서 pull → feature 브랜치 생성
- todo 파일이 있으면 순서대로 작업 진행, 완료 시 체크 표시

---

## 3. Skills 구현

### B-1: `.claude/skills/go-convention/SKILL.md`

- frontmatter: `name`, `description`(자동 호출용), `user-invocable: true`, `allowed-tools`
- 본문: 패키지 구조, 테스트 작성 규칙(testify, 테이블 테스트), 에러 처리, import 정렬

### B-2: `.claude/skills/analyze-codebase/SKILL.md`

- frontmatter: `context: fork`, `agent: Explore` (격리 실행)
- `$ARGUMENTS`로 분석 대상 경로를 받음
- 분석 항목: 디렉토리 구조, 핵심 타입, 의존성, 테스트 현황, 진입점

### B-3: `.claude/skills/go-project-layout/SKILL.md`

- frontmatter: `name`, `description`(자동 호출용), `user-invocable: true`, `allowed-tools`
- `!tree` 동적 컨텍스트로 `project-layout/go-clean-arch-v2` 실제 구조를 런타임에 읽어옴
- 레이어별 역할: cmd/, domain/, {도메인명}/, pkg/
- 의존성 방향 다이어그램 포함

### D-1 (지식): `.claude/skills/api-convention/SKILL.md`

- frontmatter: `user-invocable: false` (subagent 전용, 직접 호출 불가)
- URL 네이밍, 응답 형식, HTTP 상태 코드, Echo 프레임워크 패턴

---

## 4. Subagents 구현

### C-1: `.claude/agents/code-reviewer.md`

- frontmatter: `tools: Read, Grep, Glob, Bash` (Write/Edit 없음 → 읽기 전용)
- `model: inherit`
- git diff 기반 리뷰, 체크리스트 기반 피드백, 우선순위별 분류

### C-2: `.claude/agents/debugger.md`

- frontmatter: `tools: Read, Edit, Bash, Grep, Glob` (Edit 포함 → 수정 가능)
- 오류 분석 → 원인 파악 → 최소 수정 → go test 검증 흐름
- code-reviewer와 Edit 도구 유무로 역할 차이 강조

### C-3: `.claude/agents/test-runner.md`

- frontmatter: `tools: Read, Bash, Grep, Glob`, `disallowedTools: Write, Edit`
- `model: haiku` (저비용 모델)
- go test 실행 → 결과 분석 → 리포트 생성 (코드 수정 없음)
- `disallowedTools` 필드 사용법 시연

### D-1 (실행자): `.claude/agents/api-developer.md`

- frontmatter: `skills: [api-convention, go-project-layout]` (Skill preload)
- `tools: Read, Write, Edit, Bash, Grep, Glob`, `model: sonnet`
- clean architecture 기반 API 구현 절차: 도메인 모델 → 리포지토리 → 유스케이스 → 핸들러 → 테스트

---

## 5. 각 예제의 차별화 포인트

| 예제 | 핵심 시연 포인트 |
|------|-----------------|
| A-1 commit | Command 기본: 단순 프롬프트 템플릿 |
| A-2 plan-task | Command + `$ARGUMENTS`: 인자 전달 |
| A-3 start-task | Command + MCP 연동: 복합 워크플로우 |
| B-1 go-convention | Skill 기본: `description`으로 자동 호출 + 지식 주입 |
| B-2 analyze-codebase | Skill + `context: fork`: 격리 실행 |
| B-3 go-project-layout | Skill + `!` 동적 컨텍스트: 런타임 데이터 참조 |
| C-1 code-reviewer | Subagent 기본: 도구 제한 (읽기 전용) |
| C-2 debugger | Subagent + Edit: 수정 가능한 도구 조합 |
| C-3 test-runner | Subagent + `disallowedTools` + `model: haiku`: 세밀한 제어 |
| D-1 api-developer | Subagent + Skill preload: 지식과 실행의 결합 |
