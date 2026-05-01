# 1편 — ACL (Access Control List)

페이지마다 사용자에게 read/edit 권한을 직접 부여하는 가장 단순한 권한 모델.

## 구성

- [`backend/`](./backend/) — Go + Echo + GORM + SQLite (자세한 내용은 backend/README.md)
- [`frontend/`](./frontend/) — React 19 + TS + Vite + Tailwind v4 (자세한 내용은 frontend/README.md)

## 빠른 시작

```bash
# Backend
cd backend && go run main.go &

# Frontend (다른 터미널)
cd frontend && npm install && npm run dev
```

브라우저: http://localhost:3000

## 시드 시나리오

| 사용자 | 비밀번호 | 권한 요약 |
|---|---|---|
| alice@example.com | password | Engineering Roadmap·Public Onboarding Guide owner / Q4 Marketing Plan read |
| bob@example.com | password | Engineering Roadmap edit / Q4·Onboarding read |
| carol@example.com | password | Q4 Marketing Plan owner / Engineering Roadmap·Onboarding read |
| dave@example.com | password | Public Onboarding Guide read만 |

## 핵심 도메인 로직

ACL 평가의 본질은 [`backend/domain/acl_check.go`](./backend/domain/acl_check.go)의 `EvaluateACL` 함수에 있다 (~30줄, 5가지 평가 규칙). 블로그 글은 이 함수를 중심으로 ACL 모델을 설명한다.

## 관련 블로그 글

- 시리즈: 웹 권한 모델 비교
- 1편: ACL — 페이지 단위 공유 (작성 예정)
