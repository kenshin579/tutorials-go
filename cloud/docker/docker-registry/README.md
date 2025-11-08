# Helm charts로 설치

# todo: 실제로 테스트는 해보지 못함 (bitnami에는 docker-registry가 없음)
> helm repo add bitnami https://charts.bitnami.com/bitnami
> helm repo update
> helm install my-registry bitnami/docker-registry -f values.yaml


# push sample docker image to private registry

```bash
docker run busybox echo "Hello, World"

docker tag busybox localhost:7001/helloworld
docker push localhost:7001/helloworld
docker run localhost:7001/helloworld echo "Hello, World"
```
