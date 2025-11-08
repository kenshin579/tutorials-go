terraform {
  required_providers {
    kind = {
      source  = "tehcyx/kind"
      version = "~> 0.2.1"
    }
  }
}

provider "kind" {}

resource "kind_cluster" "ollama_cluster" {
  name            = "ollama-cluster"
  wait_for_ready  = true
  node_image      = "kindest/node:v1.27.3"
  
  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"

    node {
      role = "control-plane"
      
      extra_port_mappings {
        container_port = 30025
        host_port      = 30025
        listen_address = "127.0.0.1"
      }
    }

    node {
      role = "worker"
    }
  }
} 