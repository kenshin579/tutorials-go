variable "cluster_name" {
  description = "Kind cluster name"
  type        = string
  default     = "ingress-gateway-demo"
}

variable "kubernetes_version" {
  description = "Kubernetes version for Kind cluster"
  type        = string
  default     = "v1.28.15"
}

variable "argocd_version" {
  description = "ArgoCD Helm chart version"
  type        = string
  default     = "7.8.28"
}
