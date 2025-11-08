## redis cluster 배포 방법

- "redis cluster install" action 실행
- 수동 배포

```bash
cd {CHART_HOME}/charts/redis
helm dep update
helm upgrade --install redis-cluster --set redis-cluster.global.redis.password="password",redis-cluster.password="password" . -n database
```
