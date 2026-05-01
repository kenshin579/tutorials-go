# 1-acl frontend

React 19 + TypeScript + Vite + Tailwind v4. 백엔드(`localhost:8080`)와 함께 동작하며 dev 서버가 `/api`와 `/auth` 요청을 백엔드로 proxy한다.

## 실행

```bash
npm install
npm run dev    # http://localhost:3000
npm run build  # 프로덕션 번들
```

## 화면

- `/login` — 시드 계정 로그인 (alice / bob / carol / dave @example.com, 비밀번호 모두 `password`)
- `/pages` — 본인이 access 가능한 페이지 목록
- `/pages/:id` — 페이지 상세 (편집·저장; owner는 공유 관리 모달 진입)

## 주의

ACL 권한 검증은 모두 서버에서 수행된다. 프론트엔드의 게이팅(공유 버튼 표시/숨김 등)은 UX 개선용이며, 본 1편에서는 편집 버튼을 사전에 숨기지 않는다 — 권한 없는 사용자가 클릭해도 서버가 403을 응답한다. (사전 게이팅이 왜 ACL 모델만으로는 어려운지는 글에서 별도 다룬다.)
