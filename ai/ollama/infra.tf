terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.11"
    }
  }
}

# Kubernetes Provider 설정
provider "kubernetes" {
  config_path = "~/.kube/config"
  config_context = "kind-ollama-cluster"
}

# Helm Provider 설정
provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
    config_context = "kind-ollama-cluster"
  }
}

# 네임스페이스 생성
resource "kubernetes_namespace" "ollama" {
  metadata {
    name = "ollama"
  }
}

# Ollama Helm Chart 배포
resource "helm_release" "ollama" {
  name       = "ollama"
  namespace  = kubernetes_namespace.ollama.metadata[0].name
  repository = "https://helm.otwld.com/"
  chart      = "ollama"
  version    = "0.61.1"

  values = [
    yamlencode({
      # 이미지 설정
      image = {
        repository = "ollama/ollama"
        tag        = "latest"
        pullPolicy = "IfNotPresent"
      }

      # 서비스 설정
      service = {
        type     = "NodePort"
        port     = 11434
        nodePort = 30025
      }

      # 리소스 설정 (작은 모델용)
      resources = {
        requests = {
          cpu    = "1"
          memory = "4Gi"
        }
        limits = {
          cpu    = "2"
          memory = "8Gi"
        }
      }

      # 스토리지 설정
      persistentVolume = {
        enabled      = true
        size         = "20Gi"
        storageClass = "standard"
      }

      # 환경 변수
      ollama = {
        gpu = {
          enabled = false  # CPU 모드
        }
      }

      # Probe 설정
      livenessProbe = {
        enabled             = true
        path                = "/"
        initialDelaySeconds = 60
        periodSeconds       = 10
        timeoutSeconds      = 5
        successThreshold    = 1
        failureThreshold    = 3
      }

      readinessProbe = {
        enabled             = true
        path                = "/"
        initialDelaySeconds = 30
        periodSeconds       = 5
        timeoutSeconds      = 3
        successThreshold    = 1
        failureThreshold    = 3
      }
    })
  ]

  wait    = true
  timeout = 600
} 