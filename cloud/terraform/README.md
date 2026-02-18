# Terraform 실습: Kind 클러스터 + ArgoCD + Helm

Terraform으로 로컬 Kubernetes 클러스터를 생성하고, ArgoCD를 설치한 뒤, Helm 차트를 ArgoCD로 배포하는 전체 흐름을 실습합니다.

## 아키텍처

```
┌─ Terraform ──────→ Kind 클러스터 + ArgoCD 설치
├─ ArgoCD ─────────→ ApplicationSet으로 앱 등록
└─ Helm Charts ────→ 각 앱의 K8s 리소스 정의
```

## 디렉토리 구조

```
cloud/terraform/
├── main.tf                          # Provider 설정 + 모듈 호출
├── kind.tf                          # Kind 클러스터 + Namespace
├── variables.tf                     # 변수 정의
├── outputs.tf                       # 출력 정의
├── Makefile                         # 자동화 명령어
├── modules/
│   └── infra/
│       ├── infra.tf                 # ArgoCD + Ingress-NGINX Helm 배포
│       └── variables.tf             # 모듈 변수
├── bootstrap/
│   └── sample-apps.yaml             # ArgoCD ApplicationSet 예제
└── charts/
    └── sample-nginx/                # 샘플 Helm 차트
        ├── Chart.yaml
        ├── values.yaml
        └── templates/
            ├── deployment.yaml
            └── service.yaml
```

## 사전 준비

- [Terraform](https://developer.hashicorp.com/terraform/install) (또는 `brew install terraform`)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (`brew install kind`)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) (`brew install kubectl`)
- [ArgoCD CLI](https://argo-cd.readthedocs.io/en/stable/cli_installation/) (`brew install argocd`) (선택)

## 실습 순서

### Step 1: Terraform으로 클러스터 + ArgoCD 설치

```bash
# 초기화 (Provider 다운로드)
make tf-init

# 리소스 생성 (Kind 클러스터 + ArgoCD)
make tf-install
```

### Step 2: ArgoCD에 앱 등록

```bash
# ArgoCD ApplicationSet 적용
kubectl apply -f bootstrap/sample-apps.yaml
```

### Step 3: 배포 확인

```bash
# ArgoCD 포트포워딩
kubectl port-forward svc/argocd-server -n argocd 8080:443

# ArgoCD UI 접속: https://localhost:8080
# ID: admin / PW: password

# 배포된 앱 확인
kubectl get pods -n study
```

### Step 4: 앱 설정 변경 (ArgoCD 자동 반영)

```bash
# values.yaml에서 replicas 변경
# → Git push 후 ArgoCD가 자동으로 반영
```

### Step 5: 리소스 정리

```bash
make tf-destroy
```

## Makefile 명령어

| 명령어 | 설명 |
|--------|------|
| `make tf-init` | Terraform 초기화 (Provider 다운로드) |
| `make tf-install` | 전체 리소스 생성 |
| `make tf-validate` | 설정 파일 검증 |
| `make tf-destroy` | 전체 리소스 삭제 |
| `make tf-clean` | Terraform 캐시/상태 파일 정리 |
| `make kind-delete` | Kind 클러스터 강제 삭제 |

## 관련 블로그

- [Terraform 완벽 가이드: 기본 개념부터 GitOps 실전까지](https://blog-v2.advenoh.pe.kr) (작성 예정)
