# Utility Box

디버깅 및 관리 도구가 포함된 유틸리티 컨테이너입니다.

## 포함된 도구들

- **기본 도구**: bash, curl, wget, busybox-extras
- **네트워크 도구**: bind-tools, nmap, tcpdump, iproute2, net-tools
- **데이터베이스 클라이언트**: redis-cli, mysql-client
- **컨테이너 도구**: docker-cli
- **Kubernetes 도구**: kubectl
- **Kafka 도구**: kafkactl

## 설치

### Helm을 사용한 설치
```bash
# Helm으로 설치
helm install utility-box ./infra/utility-box

# 또는 namespace를 지정하여 설치
helm install utility-box ./infra/utility-box -n default
```

### Helm 없이 설치
```bash
# Helm template로 YAML 생성 후 적용
helm template utility-box ./infra/utility-box | kubectl apply -f -
```

## 사용 방법

### Pod에 접속하기

```bash
# Pod 이름 확인
kubectl get pods | grep utility-box

# Pod에 접속
kubectl exec -it <pod-name> -- bash
```

### 주요 사용 예시

#### MySQL 접속
```bash
mysql -h mysql-service -u root -p
```

#### Redis 접속
```bash
redis-cli -h redis-service -p 6379
```

#### Kafka 토픽 목록 확인
```bash
kafkactl get topics --brokers=kafka-broker:9092
```

#### DNS 조회
```bash
nslookup kubernetes.default.svc.cluster.local
dig @10.###.###.### kubernetes.default.svc.cluster.local
```

#### 네트워크 디버깅
```bash
# 포트 스캔
nmap -p 80,443 google.com

# 패킷 캡처 (권한 필요)
tcpdump -i eth0 -n port 80
```

## 설정 옵션

### Docker Socket 마운트

Docker 명령어를 사용하려면 Docker socket을 마운트해야 합니다:

```yaml
# values.yaml
dockerSocket:
  enabled: true
  path: /var/run/docker.sock
```

### 리소스 설정

```yaml
resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

## 빌드 및 푸시

```bash
cd infra/utility-box/scripts
make docker-push
``` 
