# =============================================================================
# Terraform 설정 및 Provider 구성
# - Terraform 블록: 필요한 Provider와 버전을 선언
# - Provider 블록: 각 Provider의 연결 설정을 정의
# =============================================================================

# Terraform에서 사용할 Provider 선언
# Provider는 Terraform이 외부 리소스를 관리하기 위한 플러그인이다
terraform {
  required_providers {
    # Kind: 로컬 Kubernetes 클러스터를 Docker 컨테이너로 생성
    kind = {
      source  = "tehcyx/kind"
      version = "0.8"
    }
    # Kubernetes: K8s 리소스(Namespace 등)를 직접 관리
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.36"
    }
    # Helm: Helm 차트를 통해 K8s 애플리케이션 배포
    helm = {
      source  = "hashicorp/helm"
      version = "3.0.0-pre2"
    }
    # Null: 로컬 명령어 실행 등 보조 작업에 사용
    null = {
      source = "hashicorp/null"
    }
  }
}

# Provider 설정
# - 각 Provider가 어떤 클러스터/환경에 연결할지 정의한다
provider "null" {}

provider "kind" {}

# Kind 클러스터 생성 후 자동으로 kubeconfig를 참조
provider "kubernetes" {
  config_path = kind_cluster.local_cluster.kubeconfig_path
}

provider "helm" {
  kubernetes = {
    config_path = kind_cluster.local_cluster.kubeconfig_path
  }
}

# 인프라 모듈 호출
# - Module: 관련 리소스를 묶어 재사용 가능한 단위로 분리한 것
# - 여기서는 ArgoCD, Ingress-NGINX 설치를 별도 모듈로 분리했다
module "infra" {
  source          = "./modules/infra"
  study_namespace = kubernetes_namespace.study.metadata[0].name

  depends_on = [
    kubernetes_namespace.study
  ]
}
