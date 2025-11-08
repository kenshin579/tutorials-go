# Helm charts로 설치

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis -f values.yaml -n database
```
