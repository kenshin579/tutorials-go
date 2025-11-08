# Docker Registry 사용 가이드

## HTTP 레지스트리 설정

도커 레지스트리(localhost:7001)를 HTTP로 사용하기 위해서는 로컬 환경의 Docker Desktop 설정에 다음과 같은 설정을 추가해야 합니다:

1. Docker Desktop 실행
2. 우측 상단의 설정(설정 톱니바퀴 아이콘) 클릭
3. 'Docker Engine' 메뉴 선택
4. JSON 설정에 다음 내용 추가:

```json
{
  "insecure-registries": ["localhost:7001"]
}
```

5. 'Apply & Restart' 버튼을 클릭하여 Docker 재시작

이 설정 없이 `make docker-push` 명령을 실행하면 "http: server gave HTTP response to HTTPS client" 에러가 발생합니다.
HTTPS 대신 HTTP 프로토콜로 레지스트리에 접근할 수 있도록 위 설정을 반드시 추가해주세요.
