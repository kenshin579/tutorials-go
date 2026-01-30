# Claude Code: Skills vs Commands vs Subagents 작업 체크리스트

## Phase 1: 디렉토리 생성

- [x] `.claude/commands/` 디렉토리 생성
- [x] `.claude/skills/go-convention/` 디렉토리 생성
- [x] `.claude/skills/analyze-codebase/` 디렉토리 생성
- [x] `.claude/skills/go-project-layout/` 디렉토리 생성
- [x] `.claude/skills/api-convention/` 디렉토리 생성
- [x] `.claude/agents/` 디렉토리 생성

## Phase 2: Commands 작성 (3개)

- [x] A-1: `.claude/commands/commit.md` 작성
- [x] A-2: `.claude/commands/plan-task.md` 작성
- [x] A-3: `.claude/commands/start-task.md` 작성

## Phase 3: Skills 작성 (4개)

- [x] B-1: `.claude/skills/go-convention/SKILL.md` 작성
- [x] B-2: `.claude/skills/analyze-codebase/SKILL.md` 작성 (`context: fork`, `agent: Explore`)
- [x] B-3: `.claude/skills/go-project-layout/SKILL.md` 작성 (`!tree` 동적 컨텍스트 포함)
- [x] D-1(지식): `.claude/skills/api-convention/SKILL.md` 작성 (`user-invocable: false`)

## Phase 4: Subagents 작성 (4개)

- [x] C-1: `.claude/agents/code-reviewer.md` 작성 (tools: Read, Grep, Glob, Bash)
- [x] C-2: `.claude/agents/debugger.md` 작성 (tools: Read, Edit, Bash, Grep, Glob)
- [x] C-3: `.claude/agents/test-runner.md` 작성 (disallowedTools: Write, Edit / model: haiku)
- [x] D-1(실행자): `.claude/agents/api-developer.md` 작성 (skills: api-convention, go-project-layout)

## Phase 5: 검증

- [ ] A-1: `/commit` 실행 → 브랜치 생성 + 커밋 메시지 자동 생성 확인
- [ ] A-2: `/plan-task` 실행 → implementation.md, todo.md 생성 확인
- [ ] A-3: `/start-task` 실행 → GitHub Issue 생성 + feature 브랜치 생성 확인
- [ ] B-1: Go 파일 작성 요청 시 `/go-convention` 자동 호출 확인
- [ ] B-2: `/analyze-codebase golang/context` → 별도 컨텍스트에서 분석 결과 반환 확인
- [ ] B-3: 새 Go 프로젝트 생성 요청 시 `/go-project-layout` 자동 적용 확인
- [ ] C-1: 코드 수정 후 `code-reviewer` 자동 위임 + Write/Edit 불가 확인
- [ ] C-2: 테스트 실패 시 `debugger` 자동 위임 + Edit으로 수정 가능 확인
- [ ] C-3: 테스트 실행 요청 시 `test-runner` 위임 + haiku 모델 + 코드 수정 없이 리포트만 반환 확인
- [ ] D-1: API 구현 요청 시 `api-developer`가 api-convention + go-project-layout 준수 확인
