# API Server

<!--
  ✅ 좋은 스펙: 하나의 관심사, "and" 없이 한 문장으로 설명 가능
  "HTTP API 서버가 /api/hello 엔드포인트에서 JSON 응답을 반환한다"

  ❌ 나쁜 스펙 (과도한 스코프):
  "HTTP 서버, 데이터베이스 연동, 인증 시스템, 로깅, 모니터링을 모두 구현한다"
  → 여러 관심사가 섞여 있음. 각각 별도 스펙 파일로 분리해야 함
-->

## 요구사항

HTTP API 서버가 `/api/hello` 엔드포인트에서 JSON 응답을 반환한다.

## 세부 사항

- `GET /api/hello`로 접근 가능
- 응답 형식: JSON
- 응답 필드: `status`, `message`, `time`
- HTTP 상태 코드: 200
- Content-Type: `application/json`

## 수락 기준

- [ ] GET /api/hello가 200 상태코드를 반환한다
- [ ] 응답이 올바른 JSON 형식이다
- [ ] status 필드가 "ok"이다
- [ ] time 필드가 RFC3339 형식이다
