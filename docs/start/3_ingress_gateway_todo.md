# Ingress vs Gateway API 구현 체크리스트

## Phase 1: 기본 인프라 구성

### Terraform
- [x] `terraform/main.tf` - Provider 설정 (kind, kubernetes, helm)
- [x] `terraform/variables.tf` - 변수 정의
- [x] `terraform/outputs.tf` - 출력 정의 (kubeconfig, endpoint)
- [x] `terraform/k8s.tf` - Kind 클러스터 설정 (control-plane + worker)
- [x] `terraform/modules/infra/infra.tf` - ArgoCD 설치

### 자동화
- [x] `Makefile` 작성 (tf-init, tf-apply, tf-destroy, deploy-*, test-*)

---

## Phase 2: 샘플 애플리케이션

### echo-server Helm Chart
- [x] `charts/echo-server/Chart.yaml`
- [x] `charts/echo-server/values.yaml`
- [x] `charts/echo-server/templates/_helpers.tpl`
- [x] `charts/echo-server/templates/deployment.yaml`
- [x] `charts/echo-server/templates/service.yaml`

### ArgoCD Application
- [x] `bootstrap/apps.yaml` - echo-server Application 정의

---

## Phase 3: Ingress 구성

### NGINX Ingress Controller
- [x] `charts/ingress/nginx-ingress/Chart.yaml`
- [x] `charts/ingress/nginx-ingress/values.yaml`

### Ingress 리소스
- [x] `charts/ingress/ingress-routes/Chart.yaml`
- [x] `charts/ingress/ingress-routes/values.yaml`
- [x] `charts/ingress/ingress-routes/templates/ingress.yaml`

### ArgoCD ApplicationSet
- [x] `bootstrap/infra-ingress.yaml`

---

## Phase 4: Gateway API 구성

### Gateway API CRDs
- [ ] `charts/gateway/gateway-api-crds/Chart.yaml`
- [ ] `charts/gateway/gateway-api-crds/values.yaml`
- [ ] `charts/gateway/gateway-api-crds/templates/crds.yaml`

### NGINX Gateway Fabric
- [ ] `charts/gateway/nginx-gateway/Chart.yaml`
- [ ] `charts/gateway/nginx-gateway/values.yaml`

### Gateway/HTTPRoute 리소스
- [ ] `charts/gateway/gateway-routes/Chart.yaml`
- [ ] `charts/gateway/gateway-routes/values.yaml`
- [ ] `charts/gateway/gateway-routes/templates/gateway.yaml`
- [ ] `charts/gateway/gateway-routes/templates/httproutes.yaml`

### ArgoCD ApplicationSet
- [ ] `bootstrap/infra-gateway.yaml`

---

## Phase 5: 테스트 및 검증

### Ingress 방식 테스트
- [ ] Kind 클러스터 생성 확인
- [ ] ArgoCD 설치 확인
- [ ] NGINX Ingress Controller 배포 확인
- [ ] echo-server 배포 확인
- [ ] Ingress 리소스 생성 확인
- [ ] `curl -H "Host: echo.local" http://localhost/ping` 테스트

### Gateway API 방식 테스트
- [ ] Gateway API CRDs 설치 확인
- [ ] NGINX Gateway Fabric 배포 확인
- [ ] Gateway 리소스 생성 확인
- [ ] HTTPRoute 리소스 생성 확인
- [ ] `curl -H "Host: echo.local" http://localhost/ping` 테스트

---

## Phase 6: 문서화

- [ ] `README.md` 작성 (프로젝트 설명, 실행 방법)
- [ ] `docs/comparison.md` 작성 (Ingress vs Gateway 비교표)
