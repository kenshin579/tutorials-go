# tutorials-go

이 저장소는 Go 언어로 작성된 튜토리얼과 예제 코드 모음입니다. 실제 프로젝트에서 자주 쓰이는 주제들을 작은 단위로 나누어 실습 중심으로 정리했습니다. 각 디렉터리는 독립적으로 동작하도록 구성되어 있어 관심 있는 주제만 골라 실행하고 학습할 수 있습니다.

## 구성
- common, golang, readme: Go 기본 문법/패턴, 유틸, 학습 자료
- database: 데이터베이스 연동 예제 및 베스트 프랙티스
- finance: 금융/환율 등 도메인 예제
- go-unit-test: 테스트 코드 작성법과 패턴
- jwt: JWT 인증/인가 관련 예제
- keycloak: Keycloak을 이용한 인증/인가 데모 (backend, frontend 포함)
- message-queue: 메시지 큐(Kafka/RabbitMQ 등) 관련 예제
- project-layout: 표준적인 Go 프로젝트 레이아웃 가이드
- scheduler: 스케줄러/잡 처리 예제
- third-party: 외부 라이브러리 활용 예제
- webhook: 웹훅 처리 예제
- cloud, chatgpt, webhook 등: 클라우드/AI 연계 데모 및 실습 코드

## 요구 사항
- Go 1.21 이상 권장 (go.mod 참고)
- 각 예제별로 추가 의존성이 있을 수 있습니다. 디렉터리 내 README를 참고하세요.

## 빠른 시작
1) 저장소 클론 후 의존성 정리
- git clone https://github.com/<your-org>/tutorials-go
- cd tutorials-go
- go mod download

2) 전체 테스트 실행
- go test ./...

3) 개별 예제 실행
- 각 디렉터리로 이동 후 README 또는 Makefile 안내에 따라 실행
- 예) keycloak/backend: 환경 변수 설정 후 `go run ./...`

## 기여
- PR과 이슈 환영합니다. 버그/제안은 이슈로 남겨주세요.

## 라이선스
- 별도 명시가 없는 경우, 소스 파일 헤더 또는 디렉터리 내 문서를 참고하세요.
