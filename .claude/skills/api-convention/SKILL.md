---
name: api-convention
description: RESTful API 설계 컨벤션
user-invocable: false
---

## API 설계 규칙

### URL 네이밍
- 복수형 명사 사용: `/users`, `/articles`
- 계층 관계: `/users/{id}/articles`
- 동사 금지: `/getUser` → `/users/{id}`

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
