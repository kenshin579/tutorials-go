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
