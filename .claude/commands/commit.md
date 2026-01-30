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
