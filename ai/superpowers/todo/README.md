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
