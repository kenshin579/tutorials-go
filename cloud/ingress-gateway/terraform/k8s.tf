resource "kind_cluster" "default" {
  name           = var.cluster_name
  wait_for_ready = true

  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"

    node {
      role  = "control-plane"
      image = "kindest/node:${var.kubernetes_version}"

      kubeadm_config_patches = [
        <<-EOF
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
        EOF
      ]

      extra_port_mappings {
        container_port = 80
        host_port      = 80
        protocol       = "TCP"
      }

      extra_port_mappings {
        container_port = 443
        host_port      = 443
        protocol       = "TCP"
      }

      extra_port_mappings {
        container_port = 30080
        host_port      = 30080
        protocol       = "TCP"
      }
    }

    node {
      role  = "worker"
      image = "kindest/node:${var.kubernetes_version}"
    }
  }
}
