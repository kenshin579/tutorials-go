# go-clean-arch-v2 프로젝트 구조 리팩토링 Todo

## Phase 1: 폴더 구조 생성

- [x] `cmd/` 폴더 생성
- [x] `pkg/config/` 폴더 생성
- [x] `pkg/database/` 폴더 생성
- [x] `pkg/middleware/` 폴더 생성

## Phase 2: 파일 이동

### pkg 관련
- [x] `main.go` → `cmd/main.go`
- [x] `common/config/config.go` → `pkg/config/config.go`
- [x] `common/database/db.go` → `pkg/database/db.go`
- [x] `article/http/middleware/middleware.go` → `pkg/middleware/middleware.go`
- [x] `article/http/middleware/middleware_test.go` → `pkg/middleware/middleware_test.go`

### article 관련
- [x] `article/http/article_handler.go` → `article/handler.go`
- [x] `article/http/article_handler_test.go` → `article/handler_test.go`
- [x] `article/usecase/article_ucase.go` → `article/usecase.go`
- [x] `article/usecase/article_ucase_test.go` → `article/usecase_test.go`
- [x] `article/repository/mysql/mysql_article.go` → `article/repository.go`
- [x] `article/repository/mysql/mysql_article_test.go` → `article/repository_test.go`
- [x] `article/repository/helper.go` → `article/helper.go`

### author 관련
- [x] `author/repository/mysql/mysql_repository.go` → `author/repository.go`
- [x] `author/repository/mysql/mysql_repository_test.go` → `author/repository_test.go`

## Phase 3: 패키지명 변경

### article 패키지
- [ ] `article/handler.go`: `package http` → `package article`
- [ ] `article/handler_test.go`: `package http_test` → `package article_test`
- [ ] `article/usecase.go`: `package usecase` → `package article`
- [ ] `article/usecase_test.go`: `package usecase_test` → `package article_test`
- [ ] `article/repository.go`: `package mysql` → `package article`
- [ ] `article/repository_test.go`: `package mysql_test` → `package article_test`
- [ ] `article/helper.go`: `package repository` → `package article`

### author 패키지
- [ ] `author/repository.go`: `package mysql` → `package author`
- [ ] `author/repository_test.go`: `package mysql_test` → `package author_test`

## Phase 4: Import 경로 수정

### cmd/main.go
- [ ] import 경로 v1 → v2 변경
- [ ] import alias 제거 및 새 경로 적용
- [ ] 함수 참조 변경 (`_articleHttp.` → `article.` 등)

### article 파일들
- [ ] `article/handler.go`: import 경로 수정
- [ ] `article/handler_test.go`: import 경로 및 참조 수정
- [ ] `article/usecase.go`: import 경로 수정
- [ ] `article/usecase_test.go`: import 경로 및 참조 수정
- [ ] `article/repository.go`: import 경로 수정 + 내부 참조 변경 (`repository.` prefix 제거)
- [ ] `article/repository_test.go`: import 경로 및 참조 수정

### author 파일들
- [ ] `author/repository.go`: import 경로 수정
- [ ] `author/repository_test.go`: import 경로 수정

### domain/mocks 파일들
- [ ] `domain/mocks/ArticleRepository.go`: import 경로 v1 → v2
- [ ] `domain/mocks/ArticleUsecase.go`: import 경로 v1 → v2
- [ ] `domain/mocks/AuthorRepository.go`: import 경로 v1 → v2

## Phase 5: 빈 폴더 삭제

- [ ] `common/` 폴더 삭제
- [ ] `article/http/` 폴더 삭제
- [ ] `article/usecase/` 폴더 삭제
- [ ] `article/repository/` 폴더 삭제
- [ ] `author/repository/` 폴더 삭제

## Phase 6: 검증

- [ ] `go mod tidy` 실행
- [ ] `go build ./...` 성공 확인
- [ ] `go test ./...` 성공 확인
- [ ] 빈 폴더 없음 확인 (`find . -type d -empty`)
