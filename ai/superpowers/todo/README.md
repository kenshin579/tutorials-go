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
make test-fe         # frontend 단위 테스트만
make test-e2e        # Playwright e2e (BE+FE 자동 기동)
```

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

### 자동 e2e 검증 (Playwright)

```bash
make test-e2e          # 헤드리스로 9개 시나리오 실행
make test-e2e:ui       # Playwright UI 모드 (디버깅용)
# 또는 직접:
cd frontend && npm run test:e2e
```

`webServer` 옵션이 BE+FE를 자동 기동/정리합니다. 이미 `:8080` 또는 `:5173`이 떠있으면 재사용합니다.

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

- `.app-header` — 제목 + 상태 카운트 ("N개 진행 중 · M개 완료")
- `.todo-form` — 흰 카드 안에 input/select/button inline
- `.filter-bar__segments` — pill segmented control (hidden radio + label, 카운트는 aria-hidden)
- `.todo-list` — 카드 컨테이너, `.todo-list--empty`는 dashed border
- `.todo-item--completed` — line-through 적용
- `.todo-item__priority--{low,medium,high}` — priority별 색
- `.error-banner` — dismissible 빨강 배너

### 다크모드

`color-scheme: light`만 선언되어 있고 다크 토큰은 미정의. 추후 별도 작업.

## 정책

- **데이터 영속성 없음**: 서버 재시작 시 모든 todo 손실 (in-memory).
- **시작 순서**: BE 먼저 띄운 후 FE. BE 부재 시 FE는 에러 배너 표시.
- **동시 PATCH**: last-write-wins (버전 필드 미도입).
- **타입 동기화**: `frontend/src/types.ts`는 백엔드 JSON 모델과 수기 동기화. BE 변경 시 같이 수정.

## 문서

- 설계: `docs/superpowers/specs/2026-04-30-todo-app-design.md`
- 구현 plan: `docs/superpowers/plans/2026-04-30-todo-app-plan.md`
