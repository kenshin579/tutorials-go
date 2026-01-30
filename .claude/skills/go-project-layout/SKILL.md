---
name: go-project-layout
description: Go 프로젝트를 새로 생성하거나 구조를 변경할 때 clean architecture 폴더 구조를 적용합니다.
user-invocable: true
allowed-tools: Read, Grep, Glob
---

## 참조 프로젝트 구조

아래는 이 프로젝트에서 사용하는 clean architecture 레이아웃입니다.
새 Go 프로젝트를 생성할 때 이 구조를 따르세요.

!`tree project-layout/go-clean-arch-v2 -I vendor -L 3`

## 레이어별 역할

### `cmd/`
- 애플리케이션 진입점 (`main.go`)
- DI 컨테이너 설정, 서버 시작

### `domain/`
- 핵심 비즈니스 엔티티 (struct)
- 리포지토리/유스케이스 인터페이스 정의
- 에러 타입 정의 (`errors.go`)
- mock 파일은 `domain/mocks/`에 위치

### `{도메인명}/` (예: `article/`, `author/`)
- `handler.go`: HTTP 핸들러 (Echo 라우팅)
- `usecase.go`: 비즈니스 로직 구현
- `repository.go`: 데이터 접근 구현
- `*_test.go`: 각 파일의 테스트

### `pkg/`
- `config/`: Viper 기반 설정 관리
- `database/`: DB 연결 설정
- `middleware/`: CORS, 인증 등 공통 미들웨어

## 의존성 방향
```
cmd/ → {도메인}/ → domain/
         ↓
        pkg/
```
- domain 패키지는 외부 의존성 없음 (순수 Go)
- 도메인별 패키지는 domain 인터페이스를 구현
- cmd/는 모든 패키지를 조립하는 역할
