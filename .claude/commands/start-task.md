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
