# Grafana Pyroscope - Go Continuous Profiling 예제

Grafana Pyroscope를 활용한 Go 애플리케이션 Continuous Profiling 예제 코드.

## 예제 구성

| 디렉토리 | 설명 | 수집 방식 |
|----------|------|----------|
| `basic/` | Pyroscope SDK 기본 연동 (CPU, 메모리, 뮤텍스 부하 생성) | Push |
| `http-server/` | Echo HTTP 서버 + Pyroscope SDK + Profiling Labels | Push |
| `pull-server/` | pprof 엔드포인트만 노출하는 HTTP 서버 | Pull |
| `alloy/` | Grafana Alloy 설정 (Pull 모드 스크래핑) | Pull |

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
| App (http-server) | http://localhost:7080 | Echo HTTP 서버 (Push 모드) |
| App (pull-server) | http://localhost:6060 | pprof 서버 (Pull 모드) |

### http-server 부하 생성 (Push 모드)

```bash
# 빠른 응답 (기준선)
curl http://localhost:7080/fast

# CPU 부하
curl http://localhost:7080/slow

# 메모리 부하
curl http://localhost:7080/memory
```

### pull-server 부하 생성 (Pull 모드)

```bash
# 빠른 응답
curl http://localhost:6060/fast

# CPU 부하
curl http://localhost:6060/slow

# 메모리 부하
curl http://localhost:6060/memory

# pprof 엔드포인트 직접 확인
curl http://localhost:6060/debug/pprof/
```

> Pull 모드에서는 Alloy가 15초마다 pprof 엔드포인트를 스크래핑한다.
> Grafana에서 `pull.golang.app` 애플리케이션으로 조회할 수 있다.

### 로컬 실행 (개발용)

> Docker Compose 없이 앱만 로컬에서 직접 실행할 때 사용한다.
> Pyroscope 서버가 필요하므로 먼저 별도로 띄워야 한다.

```bash
# 1. Pyroscope 서버만 실행 (Docker Compose를 사용했다면 생략)
docker run -d -p 4040:4040 grafana/pyroscope:latest

# 2. basic 예제 실행
cd basic && go run .

# 3. http-server 예제 실행
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
