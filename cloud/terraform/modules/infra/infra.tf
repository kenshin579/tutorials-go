# =============================================================================
# 인프라 모듈: ArgoCD + Ingress-NGINX 설치
# - Module: 관련 리소스를 묶어 재사용 가능한 단위로 분리
# - 이 모듈은 main.tf에서 module "infra" { source = "./modules/infra" }로 호출됨
# =============================================================================

# --- Ingress-NGINX ---
# 외부 트래픽을 클러스터 내부 서비스로 라우팅하는 Ingress Controller
resource "kubernetes_namespace" "ingress-nginx" {
  metadata {
    name = var.ingress-nginx_namespace
  }
}

# Helm 차트로 Ingress-NGINX 설치
# - helm_release: Terraform에서 Helm 차트를 선언적으로 관리
# - values 블록: Helm values.yaml을 인라인으로 정의
resource "helm_release" "ingress-nginx" {
  repository = "https://kubernetes.github.io/ingress-nginx"
  chart      = "ingress-nginx"
  version    = "4.12.1"
  namespace  = var.ingress-nginx_namespace
  name       = "ingress-nginx"

  depends_on = [kubernetes_namespace.ingress-nginx]

  values = [
    <<-EOT
    controller:
      replicaCount: 2
      service:
        type: "NodePort"
      resources:
        requests:
          cpu: "100m"
          memory: "256Mi"
    EOT
  ]
}

# --- ArgoCD ---
# GitOps 기반 지속적 배포(CD) 도구
# Git 저장소의 Helm 차트/매니페스트를 감시하여 K8s에 자동 배포한다
resource "kubernetes_namespace" "argocd" {
  metadata {
    name = "argocd"
  }
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "7.8.28"
  namespace  = kubernetes_namespace.argocd.metadata[0].name

  values = [
    <<-EOT
    configs:
      secret:
        argocdServerAdminPassword: ${var.argocd_password}
    server:
      service:
        type: "ClusterIP"
    EOT
  ]
}
