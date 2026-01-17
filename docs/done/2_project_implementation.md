# go-clean-arch-v2 프로젝트 구조 리팩토링 구현 문서

## 1. 폴더 구조 변경

### 1.1 새 폴더 생성

```bash
cd project-layout/go-clean-arch-v2
mkdir -p cmd pkg/config pkg/database pkg/middleware
```

### 1.2 파일 이동

```bash
# main.go 이동
mv main.go cmd/main.go

# common -> pkg 이동
mv common/config/config.go pkg/config/
mv common/database/db.go pkg/database/

# middleware 이동
mv article/http/middleware/middleware.go pkg/middleware/
mv article/http/middleware/middleware_test.go pkg/middleware/

# article 파일 평탄화
mv article/http/article_handler.go article/handler.go
mv article/http/article_handler_test.go article/handler_test.go
mv article/usecase/article_ucase.go article/usecase.go
mv article/usecase/article_ucase_test.go article/usecase_test.go
mv article/repository/mysql/mysql_article.go article/repository.go
mv article/repository/mysql/mysql_article_test.go article/repository_test.go
mv article/repository/helper.go article/helper.go

# author 파일 평탄화
mv author/repository/mysql/mysql_repository.go author/repository.go
mv author/repository/mysql/mysql_repository_test.go author/repository_test.go
```

### 1.3 빈 폴더 삭제

```bash
rm -rf common article/http article/usecase article/repository author/repository
```

## 2. 코드 수정

### 2.1 패키지명 변경

| 파일 | 변경 전 | 변경 후 |
|------|---------|---------|
| `article/handler.go` | `package http` | `package article` |
| `article/handler_test.go` | `package http_test` | `package article_test` |
| `article/usecase.go` | `package usecase` | `package article` |
| `article/usecase_test.go` | `package usecase_test` | `package article_test` |
| `article/repository.go` | `package mysql` | `package article` |
| `article/repository_test.go` | `package mysql_test` | `package article_test` |
| `article/helper.go` | `package repository` | `package article` |
| `author/repository.go` | `package mysql` | `package author` |
| `author/repository_test.go` | `package mysql_test` | `package author_test` |

### 2.2 Import 경로 변경 매핑

```
github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/common/config
→ github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/pkg/config

github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/common/database
→ github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/pkg/database

github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/http/middleware
→ github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/pkg/middleware

github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/http
→ github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/article

github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/usecase
→ github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/article

github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/repository/mysql
→ github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/article

github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/repository
→ github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/article

github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/author/repository/mysql
→ github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/author

github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/domain
→ github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/domain

github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/domain/mocks
→ github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/domain/mocks
```

### 2.3 cmd/main.go 수정 내용

```go
// 변경 전
import (
    _articleHttp "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/http"
    _articleRepo "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/repository/mysql"
    _articleUcase "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/usecase"
    _authorRepo "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/author/repository/mysql"
    "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/common/config"
    "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/common/database"
    _articleHttpMiddleware "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/http/middleware"
)

// 변경 후
import (
    "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/article"
    "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/author"
    "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/pkg/config"
    "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/pkg/database"
    "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/pkg/middleware"
)

// 함수 참조 변경
_articleHttp.NewArticleHandler → article.NewArticleHandler
_articleRepo.NewMysqlArticleRepository → article.NewMysqlArticleRepository
_articleUcase.NewArticleUsecase → article.NewArticleUsecase
_authorRepo.NewMysqlAuthorRepository → author.NewMysqlAuthorRepository
_articleHttpMiddleware.InitMiddleware → middleware.InitMiddleware
```

### 2.4 article/repository.go 수정 내용

같은 패키지 통합으로 인한 참조 변경:
```go
// 변경 전
repository.DecodeCursor(cursor)
repository.EncodeCursor(res[len(res)-1].CreatedAt)

// 변경 후
DecodeCursor(cursor)
EncodeCursor(res[len(res)-1].CreatedAt)
```

### 2.5 테스트 파일 수정 내용

**article/handler_test.go**:
```go
// 변경 전
articleHttp "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/http"
handler := articleHttp.ArticleHandler{...}

// 변경 후 (외부 테스트 패키지)
"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/article"
handler := article.ArticleHandler{...}
```

**article/usecase_test.go**:
```go
// 변경 전
ucase "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/usecase"
u := ucase.NewArticleUsecase(...)

// 변경 후 (외부 테스트 패키지)
"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/article"
u := article.NewArticleUsecase(...)
```

**article/repository_test.go**:
```go
// 변경 전
"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/repository"
articleMysqlRepo "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/repository/mysql"

// 변경 후 (외부 테스트 패키지)
"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/article"
// repository.EncodeCursor → article.EncodeCursor
// articleMysqlRepo.NewMysqlArticleRepository → article.NewMysqlArticleRepository
```

### 2.6 Mock 파일 수정

**domain/mocks/*.go**:
```go
// 변경 전
domain "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/domain"

// 변경 후
domain "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/domain"
```

## 3. 검증

```bash
cd project-layout/go-clean-arch-v2

# 의존성 정리
go mod tidy

# 빌드 확인
go build ./...

# 테스트 실행
go test ./...

# 빈 폴더 확인
find . -type d -empty
```
