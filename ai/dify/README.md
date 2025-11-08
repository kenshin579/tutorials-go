# helm 설치

## dify-helm 이용
todo: kind 클러스터에서 dify 구동이 다 안된다. 
- minikube로는 다 실행되는 거 확인함

helm repo add dify https://borispolonsky.github.io/dify-helm
helm repo update
helm install dify dify/dify

https://github.com/BorisPolonsky/dify-helm
