# Grafana Pyroscope - Go Continuous Profiling 예제

Grafana Pyroscope를 활용한 Go 애플리케이션 Continuous Profiling 예제 코드.

## 예제 구성

| 디렉토리 | 설명 |
|----------|------|
| `basic/` | Pyroscope SDK 기본 연동 (CPU, 메모리, 뮤텍스 부하 생성) |
| `http-server/` | Echo HTTP 서버 + Pyroscope + Profiling Labels |

## 실행 방법

### Docker Compose (권장)

```bash
# 전체 환경 실행 (Pyroscope + Grafana + App)
docker compose up -d

# 로그 확인
docker compose logs -f app
```

### 접속 URL

| 서비스 | URL | 설명 |
|--------|-----|------|
| Pyroscope | http://localhost:4040 | Pyroscope UI |
| Grafana | http://localhost:3000 | Grafana 대시보드 (admin/admin) |
| App (http-server) | http://localhost:8080 | Echo HTTP 서버 |

### http-server 부하 생성

```bash
# 빠른 응답 (기준선)
curl http://localhost:8080/fast

# CPU 부하
curl http://localhost:8080/slow

# 메모리 부하
curl http://localhost:8080/memory
```

### 로컬 실행 (개발용)

```bash
# Pyroscope 서버만 실행
docker run -d -p 4040:4040 grafana/pyroscope:latest

# basic 예제 실행
cd basic && go run .

# http-server 예제 실행
cd http-server && go run .
```

## Grafana에서 확인

1. http://localhost:3000 접속 (admin/admin)
2. Explore → Pyroscope 데이터소스 선택
3. Application 선택 후 Flame Graph 확인

## 참고

- [Grafana Pyroscope 공식 문서](https://grafana.com/docs/pyroscope/latest/)
- [Pyroscope Go SDK](https://grafana.com/docs/pyroscope/latest/configure-client/language-sdks/go_push/)
- [pyroscope-go GitHub](https://github.com/grafana/pyroscope-go)
