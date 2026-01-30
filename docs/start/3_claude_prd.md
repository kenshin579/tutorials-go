# Claude Code: Skills vs Commands vs Subagents 블로그 예제 PRD

## 목적

Claude Code의 세 가지 확장 메커니즘(Skills, Commands, Subagents)의 차이점을 실제 예제와 함께 설명하는 블로그를 작성하기 위한 요구사항 정의.

---

## 1. 배경: 세 가지 개념 요약

| 구분 | Custom Commands | Skills | Subagents |
|------|----------------|--------|-----------|
| **정의** | 재사용 가능한 프롬프트 단축키 | 지식/가이드라인을 주입하는 확장 시스템 | 독립된 컨텍스트에서 작업하는 전문 에이전트 |
| **위치** | `.claude/commands/` | `.claude/skills/<name>/SKILL.md` | `.claude/agents/<name>.md` |
| **호출** | `/name` (사용자 수동) | `/name` (수동) + 자동 호출 가능 | Claude가 자동 위임 또는 사용자 요청 |
| **실행 컨텍스트** | 메인 대화 (인라인) | 메인 대화 (인라인), `context: fork` 시 분리 | 별도 컨텍스트 윈도우 (격리) |
| **도구 제한** | `allowed-tools`로 제한 | `allowed-tools`로 제한 | `tools`/`disallowedTools`로 세밀한 제어 |
| **권한 모드** | 없음 | 없음 | `permissionMode` 지원 (acceptEdits, dontAsk 등) |
| **병렬 실행** | 불가 | 불가 | 가능 (여러 subagent 동시 실행) |
| **상태** | 레거시 (Skills로 전환 중) | 현재 권장 방식 | 현재 권장 방식 |

---

## 2. 블로그에 포함할 실제 예제 목록

각 기능별 최소 2개 이상의 실제 사용 가능한 예제를 포함한다.

---

### A. Custom Commands (3개)

#### 예제 A-1: Git Commit 자동화 (`/commit`)

**목적**: Commands의 기본 동작 방식을 보여주는 가장 단순한 예제

**파일**: `.claude/commands/commit.md`

**내용**:
```markdown
## 목적
수정한 파일을 commit하고 원격 저장소에 push한다

## 실행 단계
1. `git status`로 변경 사항 확인
2. 현재 브랜치가 main/master인 경우:
   - `feature/{작업내용요약}` 형식으로 새 브랜치 생성
3. 변경 파일 staging (`git add`)
4. commit 메시지 작성 (conventional commits 형식)
5. 원격 저장소에 push

## 규칙
- main/master 브랜치에는 직접 commit 금지
- commit 메시지 형식: `type: 간단한 설명`
  - type: feat, fix, docs, refactor, test, chore 등
```

**시나리오**: 사용자가 코드를 수정한 후 `/commit` 입력 → Claude가 변경사항을 분석하고 conventional commit 메시지를 자동 생성하여 커밋

**포인트**:
- Commands는 단순한 프롬프트 템플릿
- 사용자가 `/name`으로 직접 호출해야만 동작
- 메인 대화 컨텍스트에서 실행 (별도 격리 없음)

---

#### 예제 A-2: PRD에서 구현 문서 생성 (`/plan-task`)

**목적**: 인자(`$ARGUMENTS`)를 활용하는 Command 예제

**파일**: `.claude/commands/plan-task.md`

**내용**:
```markdown
## 목적
$ARGUMENTS(요구사항 파일)를 기반으로 구현 문서와 todo 체크리스트를 생성한다

## 인자
- $ARGUMENTS: 요구사항 파일 경로 (예: docs/start/1_feature_prd.md)

## 실행 단계
1. 요구사항 파일($ARGUMENTS) 읽기 및 분석
2. 구현 문서 생성 ({순서}_{기능명}_implementation.md)
   - 핵심 구현 사항만 포함
   - 향후 계획, 확장 가능성 등 불필요한 내용 제외
3. todo 체크리스트 생성 ({순서}_{기능명}_todo.md)
   - 단계별로 나눠서 작성
   - 각 항목은 체크박스 형식 (- [ ])

## 파일 생성 규칙
- 요구사항 파일과 같은 디렉토리에 생성
- 파일명 패턴: `{순서}_{기능명}_{문서타입}.md`

## 예시
입력: `/plan-task docs/start/1_github_action_prd.md`

출력:
- `docs/start/1_github_action_implementation.md`
- `docs/start/1_github_action_todo.md`
```

**시나리오**: `/plan-task docs/start/3_claude_prd.md` 입력 → PRD를 분석하여 구현 문서와 todo 체크리스트를 자동 생성

**포인트**:
- `$ARGUMENTS`로 사용자 인자를 받아 동적으로 동작
- 파일 읽기 → 분석 → 파일 생성의 다단계 워크플로우
- 여전히 메인 대화 컨텍스트에서 인라인 실행

---

#### 예제 A-3: 작업 시작 (`/start-task`)

**목적**: 여러 도구(GitHub MCP, Git)를 조합하는 다단계 워크플로우 Command 예제

**파일**: `.claude/commands/start-task.md`

**내용**:
```markdown
## 목적
$ARGUMENTS(PRD 문서)를 기반으로 작업을 시작한다

## 인자
- $ARGUMENTS: PRD 문서 경로 (예: docs/start/3_claude_prd.md)

## 실행 단계
1. GitHub Issue 생성 (GitHub MCP 사용)
   - title: PRD 제목
   - body: PRD 내용 요약
   - assignee: kenshin579
2. master 브랜치에서 최신 코드 pull
3. 새 feature 브랜치 생성: `feature/{issue번호}-{기능명}`
4. todo 파일이 있으면 순서대로 작업 진행
5. 각 단계 완료 시:
   - 변경 사항 commit
   - todo 파일에 완료 체크 표시 (- [x])

## 규칙
- master 브랜치에 직접 commit 금지
- todo 파일의 단계별 순서 준수
- commit 메시지는 conventional commits 형식

## 예시
입력: `/start-task docs/start/3_claude_prd.md`

동작:
1. GitHub Issue #607 생성: "Claude Code Skills vs Commands vs Subagents 블로그 예제"
2. `feature/607-claude-blog-examples` 브랜치 생성
3. todo 파일 순서대로 작업 시작
```

**시나리오**: `/plan-task`로 구현 문서와 todo를 생성한 후 → `/start-task`로 실제 작업을 착수

**포인트**:
- `/plan-task` → `/start-task`로 이어지는 **Command 간 워크플로우** 시연
- GitHub MCP + Git 명령을 조합하는 **복합 Command**
- 여전히 메인 대화 컨텍스트에서 인라인 실행 (Subagent와의 차이)

---

### B. Skills (3개)

#### 예제 B-1: Go 코딩 컨벤션 가이드 (`/go-convention`)

**목적**: Skill의 핵심 기능인 "자동 호출"과 "지식 주입"을 보여주는 예제

**파일**: `.claude/skills/go-convention/SKILL.md`

**내용**:
```yaml
---
name: go-convention
description: Go 코드를 작성하거나 리뷰할 때 프로젝트의 코딩 컨벤션을 적용합니다. Go 파일을 생성, 수정, 리뷰할 때 자동으로 사용됩니다.
user-invocable: true
allowed-tools: Read, Grep, Glob
---

## 이 프로젝트의 Go 코딩 컨벤션

### 패키지 구조
- 테스트 파일은 소스 파일과 같은 디렉토리에 위치: `*_test.go`
- mock 파일은 `mocks/` 하위 디렉토리에 생성
- 각 디렉토리는 독립적인 예제로 구성 (자체 `go.mod` 가능)

### 테스트 작성 규칙
- testify 사용: `github.com/stretchr/testify/assert`
- 테이블 기반 테스트 패턴 적용
- 테스트 함수명: `Test{함수명}_{시나리오}` 형식

### 에러 처리
- 커스텀 에러 타입은 `errors.New()` 또는 `fmt.Errorf("context: %w", err)` 사용
- sentinel 에러는 패키지 레벨 변수로 정의: `var ErrNotFound = errors.New("not found")`

### Import 정렬
- 표준 라이브러리 → 외부 패키지 → 내부 패키지 순서
- 내부 패키지 alias: underscore prefix 사용 (`_articleHttp`)
```

**시나리오**:
1. **자동 호출**: 사용자가 "새로운 Go 파일을 만들어줘"라고 요청 → Claude가 자동으로 이 Skill을 로드하여 컨벤션 적용
2. **수동 호출**: `/go-convention`으로 직접 호출하여 기존 코드가 컨벤션을 따르는지 확인

**포인트**:
- `description`이 있으면 Claude가 상황에 맞게 **자동 호출** (Command와의 핵심 차이)
- `disable-model-invocation: true` 설정 시 자동 호출 방지 가능
- 메인 대화에 지식이 주입되어 이후 작업에 **계속 영향**을 줌
- 보조 파일(`examples.md` 등)을 같은 디렉토리에 추가 가능 (Command 대비 장점)

---

#### 예제 B-2: 코드베이스 분석 with `context: fork` (`/analyze-codebase`)

**목적**: Skill이 Subagent처럼 격리된 컨텍스트에서 실행될 수 있음을 보여주는 예제

**파일**: `.claude/skills/analyze-codebase/SKILL.md`

**내용**:
```yaml
---
name: analyze-codebase
description: 코드베이스 구조를 분석하고 요약합니다
user-invocable: true
context: fork
agent: Explore
---

$ARGUMENTS 경로의 코드베이스를 분석하세요.

## 분석 항목
1. **디렉토리 구조**: 파일과 패키지 구성
2. **핵심 타입/인터페이스**: 주요 struct와 interface 목록
3. **의존성**: 외부 라이브러리와 내부 패키지 의존 관계
4. **테스트 현황**: 테스트 파일 유무, 테스트 패턴 (testify, testcontainers 등)
5. **진입점**: main.go 또는 주요 실행 파일

## 출력 형식
마크다운으로 요약하되, 각 파일은 `파일명:라인번호` 형식으로 참조
```

**시나리오**: `/analyze-codebase golang/concurrency` 실행 → Explore 에이전트가 별도 컨텍스트에서 해당 디렉토리를 분석하고 결과만 메인 대화에 반환

**포인트**:
- `context: fork`를 사용하면 Skill이 **별도 컨텍스트에서 격리 실행**됨
- `agent: Explore`로 사용할 에이전트 타입 지정 (빠른 탐색에 최적화)
- 대량의 파일을 읽어도 메인 컨텍스트 윈도우를 오염시키지 않음
- Skills와 Subagents의 경계가 만나는 지점: **Skill 형식 + Subagent 실행 방식**

---

#### 예제 B-3: Go 프로젝트 레이아웃 (`/go-project-layout`)

**목적**: 동적 컨텍스트(`!` 명령)로 실제 프로젝트 구조를 참조하는 Skill 예제

**파일**: `.claude/skills/go-project-layout/SKILL.md`

**내용**:
```yaml
---
name: go-project-layout
description: Go 프로젝트를 새로 생성하거나 구조를 변경할 때 clean architecture 폴더 구조를 적용합니다.
user-invocable: true
allowed-tools: Read, Grep, Glob
---

## 참조 프로젝트 구조

아래는 이 프로젝트에서 사용하는 clean architecture 레이아웃입니다.
새 Go 프로젝트를 생성할 때 이 구조를 따르세요.

!`tree project-layout/go-clean-arch-v2 -I vendor -L 3`

## 레이어별 역할

### `cmd/`
- 애플리케이션 진입점 (`main.go`)
- DI 컨테이너 설정, 서버 시작

### `domain/`
- 핵심 비즈니스 엔티티 (struct)
- 리포지토리/유스케이스 인터페이스 정의
- 에러 타입 정의 (`errors.go`)
- mock 파일은 `domain/mocks/`에 위치

### `{도메인명}/` (예: `article/`, `author/`)
- `handler.go`: HTTP 핸들러 (Echo 라우팅)
- `usecase.go`: 비즈니스 로직 구현
- `repository.go`: 데이터 접근 구현
- `*_test.go`: 각 파일의 테스트

### `pkg/`
- `config/`: Viper 기반 설정 관리
- `database/`: DB 연결 설정
- `middleware/`: CORS, 인증 등 공통 미들웨어

## 의존성 방향
```
cmd/ → {도메인}/ → domain/
         ↓
        pkg/
```
- domain 패키지는 외부 의존성 없음 (순수 Go)
- 도메인별 패키지는 domain 인터페이스를 구현
- cmd/는 모든 패키지를 조립하는 역할
```

**시나리오**:
1. **자동 호출**: "새 Go 프로젝트를 만들어줘" 요청 시 Claude가 자동으로 이 Skill을 로드하여 폴더 구조 적용
2. **수동 호출**: `/go-project-layout`으로 직접 호출하여 현재 프로젝트 구조가 규칙에 맞는지 확인

**포인트**:
- `!`(동적 컨텍스트 명령)으로 실행 시점에 **실제 프로젝트 구조를 tree로 읽어옴**
- `go-convention`(코딩 스타일)과 분리하여 **관심사 분리** → subagent에 선택적 조합 가능
- `api-developer` subagent에서 `skills: [api-convention, go-project-layout]`으로 조합 가능

---

### C. Subagents (3개)

#### 예제 C-1: 코드 리뷰 전문가 (`code-reviewer`)

**목적**: Subagent의 독립적 실행과 도구 제한을 보여주는 예제

**파일**: `.claude/agents/code-reviewer.md`

**내용**:
```yaml
---
name: code-reviewer
description: 전문 코드 리뷰 전문가. 코드의 품질, 보안 및 유지보수성을 적극적으로 검토합니다. 코드를 작성하거나 수정한 직후에 사용하세요.
tools: Read, Grep, Glob, Bash
model: inherit
---

당신은 높은 수준의 코드 품질과 보안을 보장하는 시니어 코드 리뷰어입니다.

호출될 때:
1. git diff를 실행하여 최근 변경사항을 확인
2. 수정된 파일에 집중
3. 즉시 리뷰 시작

리뷰 체크리스트:
- 코드가 간단하고 읽기 쉬운가
- 함수와 변수가 잘 명명되었는가
- 중복된 코드가 없는가
- 적절한 오류 처리가 되어 있는가
- 노출된 비밀이나 API 키가 없는가
- 입력 검증이 구현되어 있는가
- 좋은 테스트 커버리지가 있는가
- 성능 고려사항이 다뤄졌는가

우선순위별로 정리된 피드백 제공:
- 🔴 중요한 문제 (반드시 수정)
- 🟡 경고 (수정해야 함)
- 🟢 제안 (개선 고려)

문제를 해결하는 방법의 구체적인 코드 예시를 포함하세요.
```

**시나리오**:
1. 사용자가 코드를 작성/수정
2. Claude가 description의 "코드를 작성하거나 수정한 직후에 사용하세요"를 인식
3. 자동으로 code-reviewer subagent에 작업 위임
4. subagent가 별도 컨텍스트에서 `git diff` 분석 후 리뷰 결과 반환

**포인트**:
- `tools: Read, Grep, Glob, Bash` → **Write/Edit 없음** (읽기 전용, 코드 수정 불가)
- `description`에 "적극적으로 사용하세요"가 있어 Claude가 **proactive하게 자동 호출**
- 별도 컨텍스트 윈도우에서 실행 → 메인 대화에 영향 없음
- Skill과의 차이: 별도 시스템 프롬프트 + 도구 제한 + 격리 실행이 **기본 동작**

---

#### 예제 C-2: 디버깅 전문가 (`debugger`)

**목적**: Subagent의 도구 접근 설정과 문제 해결 특화를 보여주는 예제

**파일**: `.claude/agents/debugger.md`

**내용**:
```yaml
---
name: debugger
description: 오류, 테스트 실패 및 예상치 못한 동작을 위한 디버깅 전문가. 문제가 발생할 때 적극적으로 사용하세요.
tools: Read, Edit, Bash, Grep, Glob
---

당신은 근본 원인 분석을 전문으로 하는 전문 디버거입니다.

호출될 때:
1. 오류 메시지와 스택 트레이스 캡처
2. 재현 단계 식별
3. 실패 위치 격리
4. 최소한의 수정 구현
5. `go test`로 수정 사항 검증

디버깅 프로세스:
- 오류 메시지와 로그 분석
- 최근 코드 변경사항 확인 (git diff, git log)
- 가설 형성 및 테스트
- 관련 테스트 파일 확인 (*_test.go)

각 문제에 대해 다음을 제공:
- 근본 원인 설명
- 진단을 뒷받침하는 증거 (파일명:라인번호)
- 구체적인 코드 수정
- 수정 후 테스트 실행 결과

증상이 아닌 근본적인 문제를 해결하는 데 집중하세요.
```

**시나리오**:
1. 테스트 실행 시 실패 발생: `go test ./golang/context/...` → FAIL
2. Claude가 "문제가 발생할 때 적극적으로 사용하세요"를 인식
3. debugger subagent에 자동 위임
4. subagent가 테스트 로그 분석 → 원인 파악 → 코드 수정 → 재테스트

**포인트**:
- `tools: Read, Edit, Bash, Grep, Glob` → code-reviewer와 달리 **Edit 포함** (수정 가능)
- code-reviewer(읽기 전용) vs debugger(수정 가능)로 **도구 제한의 차이** 비교 가능
- 문제 발생 시 자동 위임되는 **reactive 패턴**

---

#### 예제 C-3: 테스트 실행 전문가 (`test-runner`)

**목적**: Bash 실행 중심의 Subagent로, 기존 code-reviewer/debugger와 도구 조합이 다른 패턴을 보여주는 예제

**파일**: `.claude/agents/test-runner.md`

**내용**:
```yaml
---
name: test-runner
description: Go 테스트를 실행하고 결과를 분석하는 전문가. 테스트 실행 요청이나 코드 변경 후 테스트 검증이 필요할 때 사용하세요.
tools: Read, Bash, Grep, Glob
disallowedTools: Write, Edit
model: haiku
---

당신은 Go 프로젝트의 테스트 실행 및 결과 분석 전문가입니다.

호출될 때:
1. 대상 패키지의 테스트 파일 확인 (*_test.go)
2. `go test` 실행 (verbose + coverage)
3. 결과 분석 및 리포트 생성

실행 명령:
- 단일 패키지: `go test -v -cover ./path/to/package/...`
- 특정 테스트: `go test -v -run TestFunctionName ./path/to/package`
- 전체: `go test -v -cover ./...`

리포트 형식:
- 총 테스트 수 / 성공 / 실패 / 스킵
- 실패한 테스트별: 테스트명, 에러 메시지, 실패 위치 (파일명:라인번호)
- 커버리지: 패키지별 커버리지 %
- 실행 시간

실패 분석:
- 실패한 테스트의 expected vs actual 값 비교
- 관련 소스 코드 참조 (파일명:라인번호)
- 가능한 원인 추정 (코드 수정은 하지 않음)

코드를 직접 수정하지 마세요. 분석 결과만 제공합니다.
```

**시나리오**:
1. "golang/context 패키지 테스트를 돌려줘"라고 요청
2. Claude가 test-runner subagent에 위임
3. subagent가 `go test -v -cover` 실행 → 결과 분석 → 리포트 반환

**포인트**:
- `disallowedTools: Write, Edit` → 코드 수정 불가, **분석만 수행** (code-reviewer와 유사하지만 Bash로 실행까지)
- `model: haiku` → 단순 실행+분석 작업에 **저비용 모델** 사용 (Subagent의 모델 선택 기능)
- 3개 Subagent의 도구 조합 비교:

| Subagent | 핵심 도구 | 역할 |
|----------|-----------|------|
| code-reviewer | Read, Grep (Edit 없음) | 읽고 분석 |
| debugger | Read, Edit, Bash | 읽고, 실행하고, 수정 |
| test-runner | Read, Bash (Edit 없음) | 읽고, 실행하고, 분석 |

---

### D. Subagent + Skill 연동 (1개)

#### 예제 D-1: API 개발 전문가 (`api-developer` + `api-convention`)

**목적**: Subagent가 Skill을 preload하여 도메인 지식을 갖춘 전문가로 동작하는 예제

**파일**:
- `.claude/skills/api-convention/SKILL.md` (지식)
- `.claude/agents/api-developer.md` (전문가)

**Skill 내용** (`.claude/skills/api-convention/SKILL.md`):
```yaml
---
name: api-convention
description: RESTful API 설계 컨벤션
user-invocable: false
---

## API 설계 규칙

### URL 네이밍
- 복수형 명사 사용: `/users`, `/articles`
- 계층 관계: `/users/{id}/articles`
- 동사 금지: `/getUser` ❌ → `/users/{id}` ✅

### 응답 형식
- 성공: `{ "data": ..., "meta": { "page": 1, "total": 100 } }`
- 에러: `{ "error": { "code": "NOT_FOUND", "message": "..." } }`

### HTTP 상태 코드
- 200: 성공, 201: 생성, 204: 삭제 성공
- 400: 잘못된 요청, 401: 미인증, 403: 권한 없음, 404: 없음
- 500: 서버 오류

### Echo 프레임워크 패턴
- 핸들러는 `http/` 디렉토리에 위치
- 라우트 그룹: `e.Group("/api/v1")`
- 미들웨어: CORS, JWT 검증은 그룹 레벨에 적용
```

**Subagent 내용** (`.claude/agents/api-developer.md`):
```yaml
---
name: api-developer
description: API 엔드포인트를 구현하는 전문 개발자. Echo 프레임워크 기반의 RESTful API를 설계하고 구현합니다.
tools: Read, Write, Edit, Bash, Grep, Glob
model: sonnet
skills:
  - api-convention
  - go-project-layout
---

당신은 Go Echo 프레임워크 기반의 API 개발 전문가입니다.

호출될 때:
1. 요청된 API의 도메인 모델(struct) 정의
2. 리포지토리 인터페이스 설계
3. 유스케이스(비즈니스 로직) 구현
4. HTTP 핸들러 작성
5. 라우터 등록
6. 테스트 작성

프로젝트 구조 (clean architecture):
- `domain/`: 엔티티 및 인터페이스
- `repository/`: 데이터 접근 구현
- `usecase/`: 비즈니스 로직
- `http/`: 핸들러 및 라우터

preload된 api-convention skill의 규칙을 반드시 준수하세요.
```

**시나리오**:
- "사용자 CRUD API를 추가해줘"라고 요청
- Claude가 api-developer subagent에 위임
- subagent는 시작 시 api-convention skill 내용을 자동 로드
- API 컨벤션을 준수하면서 엔드포인트 구현

**포인트**:
- `skills: [api-convention, go-project-layout]` 필드로 subagent에 도메인 지식 주입
- Skill의 `user-invocable: false`로 직접 호출은 차단하되 subagent에서는 사용 가능
- **Skills = 지식, Subagents = 실행자**라는 역할 분담이 명확

---

## 3. 핵심 비교 포인트 (블로그 본문에 반드시 포함)

### 3.1 언제 무엇을 쓸까?

| 상황 | 추천 |
|------|------|
| 단순한 단축 명령어가 필요할 때 | Command 또는 Skill |
| 인자를 받아 정해진 절차를 실행할 때 | Command (`$ARGUMENTS`) |
| 코딩 컨벤션/가이드라인을 적용할 때 | Skill (자동 호출 활용) |
| 독립적인 작업을 격리 실행할 때 | Subagent |
| 읽기 전용으로 도구를 제한할 때 | Subagent (`tools` 필드) |
| 여러 작업을 병렬 실행할 때 | Subagent (여러 개 동시 실행) |
| 지식 + 실행을 결합할 때 | Subagent + Skill (`skills` 필드) |
| Skill인데 격리 실행이 필요할 때 | Skill (`context: fork`) |

### 3.2 비유로 설명

- **Commands** = 매크로 (반복 작업을 단축키로 실행)
- **Skills** = 참고 문서 (책상 위에 펼쳐놓고 작업하면서 참조)
- **Subagents** = 팀원 (별도의 책상에서 독립적으로 작업 후 결과만 전달)

### 3.3 진화 흐름

```
Commands (레거시) → Skills (현재 권장) → Subagents (고급 위임)
                          ↓
                    context: fork로
                    Subagent 영역까지 확장
```

### 3.4 제약사항

| 제약 | Commands | Skills | Subagents |
|------|----------|--------|-----------|
| 자동 호출 | 불가 | 가능 (`description` 필요) | 가능 (`description` 필요) |
| 컨텍스트 격리 | 불가 | `context: fork`로 가능 | 기본 동작 |
| 다른 agent 생성 | 불가 | 불가 | 불가 (nesting 불가) |
| MCP 도구 사용 | 가능 | 가능 | 백그라운드 실행 시 불가 |
| 병렬 실행 | 불가 | 불가 | 가능 |
| 보조 파일 지원 | 불가 (단일 파일) | 가능 (디렉토리 구조) | 가능 (단일 파일) |

---

## 4. 예제 작성 범위

### 실제 생성할 파일 목록

모든 파일은 이 프로젝트의 `.claude/` 디렉토리에 생성한다.

| # | 파일 경로 | 타입 | 예제 |
|---|-----------|------|------|
| 1 | `.claude/commands/commit.md` | Command | A-1: Git Commit 자동화 |
| 2 | `.claude/commands/plan-task.md` | Command | A-2: PRD → 구현 문서 생성 |
| 3 | `.claude/commands/start-task.md` | Command | A-3: 작업 시작 |
| 4 | `.claude/skills/go-convention/SKILL.md` | Skill | B-1: Go 코딩 컨벤션 |
| 5 | `.claude/skills/analyze-codebase/SKILL.md` | Skill (fork) | B-2: 코드베이스 분석 |
| 6 | `.claude/skills/go-project-layout/SKILL.md` | Skill | B-3: Go 프로젝트 레이아웃 |
| 7 | `.claude/agents/code-reviewer.md` | Subagent | C-1: 코드 리뷰 전문가 |
| 8 | `.claude/agents/debugger.md` | Subagent | C-2: 디버깅 전문가 |
| 9 | `.claude/agents/test-runner.md` | Subagent | C-3: 테스트 실행 전문가 |
| 10 | `.claude/skills/api-convention/SKILL.md` | Skill | D-1: API 컨벤션 (지식) |
| 11 | `.claude/agents/api-developer.md` | Subagent | D-1: API 개발자 (실행자) |

### 작성하지 않는 것

- 블로그 본문 (이 PRD의 범위가 아님)

---

## 5. 검증 방법

각 예제가 의도대로 동작하는지 다음과 같이 검증:

| # | 예제 | 검증 방법 |
|---|------|-----------|
| A-1 | `/commit` | 파일 수정 후 `/commit` 실행 → 브랜치 생성 + 커밋 메시지 자동 생성 확인 |
| A-2 | `/plan-task` | `/plan-task docs/start/3_claude_prd.md` 실행 → implementation.md, todo.md 생성 확인 |
| A-3 | `/start-task` | `/start-task docs/start/3_claude_prd.md` 실행 → GitHub Issue 생성 + feature 브랜치 생성 확인 |
| B-1 | `/go-convention` | Go 파일 작성 요청 시 컨벤션 자동 적용 여부 확인 (testify, 테이블 테스트 등) |
| B-2 | `/analyze-codebase` | `/analyze-codebase golang/context` 실행 → 별도 컨텍스트에서 분석 결과 반환 확인 |
| B-3 | `/go-project-layout` | 새 Go 프로젝트 생성 요청 시 clean architecture 폴더 구조 자동 적용 확인 + `!tree` 동적 컨텍스트 동작 확인 |
| C-1 | `code-reviewer` | 코드 수정 후 리뷰 자동 위임 확인 + Write/Edit 도구 사용 불가 확인 |
| C-2 | `debugger` | 테스트 실패 시 자동 위임 확인 + Edit 도구로 코드 수정 가능 확인 |
| C-3 | `test-runner` | 테스트 실행 요청 시 위임 확인 + haiku 모델 사용 확인 + 코드 수정 없이 리포트만 반환 확인 |
| D-1 | `api-developer` | API 구현 요청 시 api-convention + go-project-layout skill 내용 준수 확인 |
