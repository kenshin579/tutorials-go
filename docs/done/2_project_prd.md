# go-clean-arch-v2 프로젝트 구조 리팩토링 PRD

## 1. 개요

`project-layout/go-clean-arch-v2` 프로젝트의 폴더 구조를 표준 Go 프로젝트 레이아웃에 맞게 리팩토링한다.

## 2. 현재 구조

```
go-clean-arch-v2/
├── article/
│   ├── http/
│   │   ├── article_handler.go          (package: http)
│   │   ├── article_handler_test.go     (package: http_test)
│   │   └── middleware/
│   │       ├── middleware.go           (package: middleware)
│   │       └── middleware_test.go
│   ├── repository/
│   │   ├── helper.go                   (package: repository)
│   │   └── mysql/
│   │       ├── mysql_article.go        (package: mysql)
│   │       └── mysql_article_test.go   (package: mysql_test)
│   └── usecase/
│       ├── article_ucase.go            (package: usecase)
│       └── article_ucase_test.go       (package: usecase_test)
├── author/
│   └── repository/
│       └── mysql/
│           ├── mysql_repository.go     (package: mysql)
│           └── mysql_repository_test.go (package: mysql_test)
├── common/
│   ├── config/
│   │   └── config.go                   (package: config)
│   └── database/
│       └── db.go                       (package: database)
├── domain/
│   ├── article.go                      (package: domain)
│   ├── author.go
│   ├── errors.go
│   └── mocks/
│       ├── ArticleRepository.go        (package: mocks)
│       ├── ArticleUsecase.go
│       └── AuthorRepository.go
└── main.go                             (package: main)
```

## 3. 목표 구조

```
go-clean-arch-v2/
├── cmd/
│   └── main.go                         (package: main)
├── pkg/
│   ├── config/
│   │   └── config.go                   (package: config)
│   ├── database/
│   │   └── db.go                       (package: database)
│   └── middleware/
│       ├── middleware.go               (package: middleware)
│       └── middleware_test.go          (package: middleware_test)
├── article/
│   ├── handler.go                      (package: article)
│   ├── handler_test.go                 (package: article_test)
│   ├── usecase.go                      (package: article)
│   ├── usecase_test.go                 (package: article_test)
│   ├── repository.go                   (package: article)
│   ├── repository_test.go              (package: article_test)
│   └── helper.go                       (package: article)
├── author/
│   ├── repository.go                   (package: author)
│   └── repository_test.go              (package: author_test)
└── domain/
    ├── article.go                      (package: domain)
    ├── author.go
    ├── errors.go
    └── mocks/
        ├── ArticleRepository.go        (package: mocks)
        ├── ArticleUsecase.go
        └── AuthorRepository.go
```

## 4. 작업 목록

### 4.1 폴더 생성

| 작업 ID | 작업 내용 |
|---------|----------|
| F-001 | `cmd/` 폴더 생성 |
| F-002 | `pkg/config/` 폴더 생성 |
| F-003 | `pkg/database/` 폴더 생성 |
| F-004 | `pkg/middleware/` 폴더 생성 |

### 4.2 파일 이동 및 이름 변경

| 작업 ID | 원본 경로 | 대상 경로 |
|---------|----------|----------|
| M-001 | `main.go` | `cmd/main.go` |
| M-002 | `common/config/config.go` | `pkg/config/config.go` |
| M-003 | `common/database/db.go` | `pkg/database/db.go` |
| M-004 | `article/http/middleware/middleware.go` | `pkg/middleware/middleware.go` |
| M-005 | `article/http/middleware/middleware_test.go` | `pkg/middleware/middleware_test.go` |
| M-006 | `article/http/article_handler.go` | `article/handler.go` |
| M-007 | `article/http/article_handler_test.go` | `article/handler_test.go` |
| M-008 | `article/usecase/article_ucase.go` | `article/usecase.go` |
| M-009 | `article/usecase/article_ucase_test.go` | `article/usecase_test.go` |
| M-010 | `article/repository/mysql/mysql_article.go` | `article/repository.go` |
| M-011 | `article/repository/mysql/mysql_article_test.go` | `article/repository_test.go` |
| M-012 | `article/repository/helper.go` | `article/helper.go` |
| M-013 | `author/repository/mysql/mysql_repository.go` | `author/repository.go` |
| M-014 | `author/repository/mysql/mysql_repository_test.go` | `author/repository_test.go` |

### 4.3 삭제할 폴더

| 작업 ID | 폴더 경로 | 설명 |
|---------|----------|------|
| D-001 | `common/` | 파일 이동 후 빈 폴더 삭제 |
| D-002 | `article/http/` | 파일 이동 후 빈 폴더 삭제 |
| D-003 | `article/usecase/` | 파일 이동 후 빈 폴더 삭제 |
| D-004 | `article/repository/` | 파일 이동 후 빈 폴더 삭제 |
| D-005 | `author/repository/` | 파일 이동 후 빈 폴더 삭제 |

### 4.4 패키지명 변경

| 작업 ID | 파일 | 변경 전 | 변경 후 |
|---------|-----|---------|---------|
| P-001 | `article/handler.go` | `package http` | `package article` |
| P-002 | `article/handler_test.go` | `package http_test` | `package article_test` |
| P-003 | `article/usecase.go` | `package usecase` | `package article` |
| P-004 | `article/usecase_test.go` | `package usecase_test` | `package article_test` |
| P-005 | `article/repository.go` | `package mysql` | `package article` |
| P-006 | `article/repository_test.go` | `package mysql_test` | `package article_test` |
| P-007 | `article/helper.go` | `package repository` | `package article` |
| P-008 | `author/repository.go` | `package mysql` | `package author` |
| P-009 | `author/repository_test.go` | `package mysql_test` | `package author_test` |

### 4.5 Import 경로 변경

**Base Path 변경**: `go-clean-arch-v1` → `go-clean-arch-v2`

| 작업 ID | 변경 전 Import | 변경 후 Import |
|---------|---------------|---------------|
| I-001 | `.../go-clean-arch-v1/common/config` | `.../go-clean-arch-v2/pkg/config` |
| I-002 | `.../go-clean-arch-v1/common/database` | `.../go-clean-arch-v2/pkg/database` |
| I-003 | `.../go-clean-arch-v1/article/http/middleware` | `.../go-clean-arch-v2/pkg/middleware` |
| I-004 | `.../go-clean-arch-v1/article/http` | `.../go-clean-arch-v2/article` |
| I-005 | `.../go-clean-arch-v1/article/usecase` | `.../go-clean-arch-v2/article` |
| I-006 | `.../go-clean-arch-v1/article/repository/mysql` | `.../go-clean-arch-v2/article` |
| I-007 | `.../go-clean-arch-v1/article/repository` | `.../go-clean-arch-v2/article` |
| I-008 | `.../go-clean-arch-v1/author/repository/mysql` | `.../go-clean-arch-v2/author` |
| I-009 | `.../go-clean-arch-v1/domain` | `.../go-clean-arch-v2/domain` |
| I-010 | `.../go-clean-arch-v1/domain/mocks` | `.../go-clean-arch-v2/domain/mocks` |

### 4.6 함수/타입 참조 변경

패키지 통합으로 인해 alias가 필요없어지거나 변경되는 항목:

| 작업 ID | 파일 | 변경 내용 |
|---------|-----|----------|
| R-001 | `cmd/main.go` | `_articleHttp.NewArticleHandler` → `article.NewArticleHandler` |
| R-002 | `cmd/main.go` | `_articleRepo.NewMysqlArticleRepository` → `article.NewMysqlArticleRepository` |
| R-003 | `cmd/main.go` | `_articleUcase.NewArticleUsecase` → `article.NewArticleUsecase` |
| R-004 | `cmd/main.go` | `_authorRepo.NewMysqlAuthorRepository` → `author.NewMysqlAuthorRepository` |
| R-005 | `cmd/main.go` | `_articleHttpMiddleware.InitMiddleware` → `middleware.InitMiddleware` |
| R-006 | `article/handler_test.go` | `articleHttp.ArticleHandler` → `article.ArticleHandler` (내부 참조) |
| R-007 | `article/usecase_test.go` | `ucase.NewArticleUsecase` → `article.NewArticleUsecase` (내부 참조, 또는 그냥 `NewArticleUsecase`) |
| R-008 | `article/repository_test.go` | `articleMysqlRepo.NewMysqlArticleRepository` → `article.NewMysqlArticleRepository` (내부 참조) |
| R-009 | `article/repository.go` | `repository.DecodeCursor` → `DecodeCursor` (같은 패키지) |
| R-010 | `article/repository.go` | `repository.EncodeCursor` → `EncodeCursor` (같은 패키지) |

## 5. 영향받는 파일 목록

| 파일 | 수정 유형 |
|-----|----------|
| `cmd/main.go` | 이동 + import 수정 + 함수 참조 수정 |
| `pkg/config/config.go` | 이동만 |
| `pkg/database/db.go` | 이동만 |
| `pkg/middleware/middleware.go` | 이동만 |
| `pkg/middleware/middleware_test.go` | 이동만 |
| `article/handler.go` | 이동 + 패키지명 변경 + import 수정 |
| `article/handler_test.go` | 이동 + 패키지명 변경 + import 수정 + 참조 수정 |
| `article/usecase.go` | 이동 + 패키지명 변경 + import 수정 |
| `article/usecase_test.go` | 이동 + 패키지명 변경 + import 수정 + 참조 수정 |
| `article/repository.go` | 이동 + 패키지명 변경 + import 수정 + 참조 수정 |
| `article/repository_test.go` | 이동 + 패키지명 변경 + import 수정 + 참조 수정 |
| `article/helper.go` | 이동 + 패키지명 변경 |
| `author/repository.go` | 이동 + 패키지명 변경 + import 수정 |
| `author/repository_test.go` | 이동 + 패키지명 변경 + import 수정 |
| `domain/mocks/ArticleRepository.go` | import 수정 |
| `domain/mocks/ArticleUsecase.go` | import 수정 |
| `domain/mocks/AuthorRepository.go` | import 수정 |

## 6. 검증 항목

| 검증 ID | 검증 내용 | 명령어 |
|---------|----------|--------|
| V-001 | Go 모듈 의존성 정리 | `go mod tidy` |
| V-002 | 빌드 성공 확인 | `go build ./...` |
| V-003 | 테스트 통과 확인 | `go test ./...` |
| V-004 | 불필요한 빈 폴더 없음 | `find . -type d -empty` |

## 7. 작업 순서 (권장)

1. **폴더 생성** (F-001 ~ F-004)
2. **파일 이동 및 이름 변경** (M-001 ~ M-014)
3. **패키지명 변경** (P-001 ~ P-009)
4. **Import 경로 변경** (I-001 ~ I-010)
5. **함수/타입 참조 변경** (R-001 ~ R-010)
6. **빈 폴더 삭제** (D-001 ~ D-005)
7. **검증** (V-001 ~ V-004)

## 8. 주의사항

- 모든 import 경로가 `go-clean-arch-v1`을 참조하고 있으므로 `go-clean-arch-v2`로 변경 필요
- `article` 패키지로 통합 시 동일 패키지 내 함수 호출은 패키지 prefix 불필요
- test 파일은 `_test` suffix 패키지 사용 가능 (external test package)
- mock 파일들은 자동 생성된 코드이므로 import 수정 후 재생성 권장
