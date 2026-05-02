# 3-abac frontend

React 19 + TypeScript + Vite + Tailwind v4. 백엔드(`localhost:8082`)와 함께 동작하며 dev 서버가 `/api`와 `/auth` 요청을 백엔드로 proxy한다. 1·2편과 동시에 띄울 수 있도록 `:3002`를 사용한다.

## 실행

```bash
npm install
npm run dev    # http://localhost:3002
npm run build  # 프로덕션 번들
```

## 화면

- `/login` — 시드 계정 로그인
- `/pages` — ABAC read 정책을 통과한 페이지만 표시. 페이지마다 `confidentiality` + 부서 뱃지
- `/pages/:id` — 페이지 상세 + **Decision 카드** (read/edit 각각의 reason과 policy 식별자를 사용자에게 표시)

## ABAC 결과 표시 (`Decision` 카드)

1·2편과의 가장 큰 UX 차이는 페이지 상세 화면이다. 서버가 내려준 `Decision{allowed, reason, policy}`를 그대로 카드로 보여준다.

```
✓ 읽기 권한  policy: internal
   same department, internal page

✗ 편집 권한  policy: department-match
   different department
```

권한이 없을 때 단순히 "거부"가 아니라 **"왜 거부됐는지"** 를 사용자에게 명시한다 — ABAC 평가 결과가 풍부하기 때문에 가능한 UX다.

## 주의

권한 평가는 항상 **서버**가 한다. 클라이언트는 서버가 내려준 decision을 표시할 뿐이며, 권한 없는 사용자가 직접 API를 호출해도 서버가 정책 평가로 막는다.
