# Sealed Secrets 설치 가이드

Charts에 Sealed Secrets 을 사용하기 위해 Sealed Secrets CLI과 Sealed Secret Controller를 설치하는 방법에 대해서 기술합니다. 

## 1. kubeseal CLI 설치

`kubeseal` 명령어는 암호화된 시크릿을 생성할 때 사용하고 아래 명령어로 설치합니다. 

### Homebrew

```bash
# 맥에서 brew로 설치
> brew install kubeseal
```

참고

- [Kubeseal Installation](https://github.com/bitnami-labs/sealed-secrets?tab=readme-ov-file#kubeseal)


## 2. Sealed Secrets Controller 설치

infra-charts에 있는 버전으로 설치합니다. 

```bash
# infra-charts에 있는 버전으로 설치

cd sealed-secrets
> helm dep update
> helm upgrade --install sealed-secrets . -n kube-system -f values.yaml
```

설치가 잘 되었는지 아래 명령어로 확인할 수 있습니다.  

```bash
> kubectl get pods -n kube-system -l app.kubernetes.io/name=sealed-secrets
```


## 3. Controller에 의해서 Sealed Secrets이 잘 동작하는지 확인 (테스트 용도)

Sealed Secrets Controller 설치하면 추가 설정은 필요가 없습니다. 
잘 동작을 하는지 테스트 용도로 아래 실행해서 확인할 수 있습니다. 대부분 잘 동작을 할거라서, 그냥 스터디 차원에서 정리해둡니다. 

### 3.1 Secret 생성하기

```bash
> kubectl create -n frank secret generic mysecret --from-literal hello=world --dry-run=client -oyaml > secret.yaml
```

위 명령어를 실행하면 아래 yaml이 생성이 됩니다. 

```yaml
apiVersion: v1
data:
  hello: d29ybGQ=
kind: Secret
metadata:
  creationTimestamp: null
  name: mysecret
  namespace: frank

```

> secret 값은 base64로 무조건 인코딩을 해야 한다. 그리고 echo에 -n 옵션을 추가히자 않으면 newline이 추가된다. 
>
> ```bash
> echo -n 'your-value' | base64
> ```


### 3.2 Sealed Secret 생성하기

kubeseal 명령어로 sealed sealed 값을 생성합니다. 

```bash
> kubeseal --controller-name=sealed-secrets-controller --controller-namespace=kube-system < secret.yaml > sealedsecret.yaml
```

### 3.3 Sealed Secret 적용후 확인

Sealed Secret을 Kubernetes 클러스터에 적용합니다. 

```bash
> kubectl apply -f sealedsecret.yaml
```

Sealed Secret 값이 잘 생성되었는지 아래 명령어로 확인합니다. 

```bash
kubectl get secret mysecret -o yaml
```
