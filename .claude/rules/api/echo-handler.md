---
paths:
  - "**/handler*.go"
  - "**/route*.go"
---

# Echo API 핸들러 규칙

- Echo v4 프레임워크 사용
- 핸들러 시그니처: func(c echo.Context) error
- 에러 응답: echo.NewHTTPError(status, message)
- 입력 바인딩: c.Bind(&req)로 구조체에 바인딩
- 미들웨어 체인에서 인증/로깅 처리
