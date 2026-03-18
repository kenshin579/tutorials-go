# Go 코드 스타일

- gofmt/goimports 적용 필수
- 에러는 즉시 처리 (if err != nil 패턴)
- 패키지 export 함수에 GoDoc 주석 작성
- 변수명은 짧고 관용적으로 (err, ctx, req, resp)
- 구조체 필드 태그는 json, validate 순서로 작성
