---
name: go-convention
description: Go 코드를 작성하거나 리뷰할 때 프로젝트의 코딩 컨벤션을 적용합니다. Go 파일을 생성, 수정, 리뷰할 때 자동으로 사용됩니다.
user-invocable: true
allowed-tools: Read, Grep, Glob
---

## 이 프로젝트의 Go 코딩 컨벤션

### 패키지 구조
- 테스트 파일은 소스 파일과 같은 디렉토리에 위치: `*_test.go`
- mock 파일은 `mocks/` 하위 디렉토리에 생성
- 각 디렉토리는 독립적인 예제로 구성 (자체 `go.mod` 가능)

### 테스트 작성 규칙
- testify 사용: `github.com/stretchr/testify/assert`
- 테이블 기반 테스트 패턴 적용
- 테스트 함수명: `Test{함수명}_{시나리오}` 형식

### 에러 처리
- 커스텀 에러 타입은 `errors.New()` 또는 `fmt.Errorf("context: %w", err)` 사용
- sentinel 에러는 패키지 레벨 변수로 정의: `var ErrNotFound = errors.New("not found")`

### Import 정렬
- 표준 라이브러리 → 외부 패키지 → 내부 패키지 순서
- 내부 패키지 alias: underscore prefix 사용 (`_articleHttp`)
