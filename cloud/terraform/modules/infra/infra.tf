# namespace 을 생성한다
resource "kubernetes_namespace" "ingress-nginx" {
  metadata {
    name = var.ingress-nginx_namespace
  }
}

# helm으로 ingress-nginx를 설치한다
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

# argocd namespace을 생성한다
resource "kubernetes_namespace" "argocd" {
  metadata {
    name = "argocd"
  }
}

# helm으로 argocd를 설치한다
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
