terraform {
  required_providers {
    kind = {
      source  = "tehcyx/kind"
      version = "0.8"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.36"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "3.0.0-pre2"
    }
    null = {
      source = "hashicorp/null"
    }
  }
}

provider "null" {}

provider "kind" {}

provider "kubernetes" {
  config_path = kind_cluster.local_cluster.kubeconfig_path
}

provider "helm" {
  kubernetes = {
    config_path = kind_cluster.local_cluster.kubeconfig_path
  }
}

module "infra" {
  source = "./modules/infra"
  study_namespace = kubernetes_namespace.study.metadata[0].name

  depends_on = [
    kubernetes_namespace.study
  ]
}
