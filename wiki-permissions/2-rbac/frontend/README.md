# 2-rbac frontend

React 19 + TypeScript + Vite + Tailwind v4. 백엔드(`localhost:8081`)와 함께 동작하며 dev 서버가 `/api`와 `/auth` 요청을 백엔드로 proxy한다.

1편(`:3000`)과 동시에 띄울 수 있도록 `:3001`을 사용한다.

## 실행

```bash
npm install
npm run dev    # http://localhost:3001
npm run build  # 프로덕션 번들
```

## 화면

- `/login` — 시드 계정 로그인 (alice/bob/carol/dave @example.com, 비밀번호 모두 `password`)
- `/pages` — 페이지 목록 (모든 사용자가 모든 페이지를 봄). `pages:create` 권한이 있으면 새 페이지 버튼 표시
- `/pages/:id` — 페이지 상세. `pages:edit`/`pages:delete` 권한별로 버튼 사전 게이팅
- `/users` — admin 전용. 사용자 목록 + role 부여/회수

## 권한 사전 게이팅 (PermissionGate)

1편(ACL)과의 핵심 차이점이다. 1편에서는 편집 버튼이 모두에게 보이고 서버가 403으로 거부했지만, RBAC에서는 login 응답에 사용자 권한이 함께 오므로 클라이언트가 사전 게이팅 가능하다.

```tsx
<PermissionGate permission="pages:edit">
  <button>편집</button>
</PermissionGate>
```

권한이 없으면 `PermissionGate`는 `null`을 반환해 버튼 자체가 렌더되지 않는다. 단, 보안은 여전히 서버가 책임진다 — 프론트엔드 게이팅은 UX 보조일 뿐이다.
