resource "kubernetes_namespace" "argocd" {
  metadata {
    name = "argocd"
  }
}

resource "kubernetes_namespace" "app" {
  metadata {
    name = "app"
  }
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "7.8.28"
  namespace  = kubernetes_namespace.argocd.metadata[0].name

  set {
    name  = "configs.params.server.insecure"
    value = "true"
  }

  # password: password (bcrypt hash)
  set {
    name  = "configs.secret.argocdServerAdminPassword"
    value = "$2a$10$rRyBsGSHK6.uc8fntPwVIuLVHgsAhAX7TcdrqW/RBER3kHIwBiP3C"
  }

  timeout = 600
}
