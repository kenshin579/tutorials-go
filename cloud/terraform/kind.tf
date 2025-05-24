resource "kind_cluster" "local_cluster" {
  name = var.kind_cluster_name
  wait_for_ready = true
  node_image = "kindest/node:v1.28.15"

  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"

    node {
      role = "control-plane"
      extra_port_mappings {
        container_port = 30080
        host_port      = 30080
        listen_address = "127.0.0.1"
      }

      extra_mounts {
        host_path = "/tmp/kind-storage"
        container_path = "/opt/local-path-provisioner"
      }

      kubeadm_config_patches = [<<-EOF
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            system-reserved: cpu=250m,memory=250Mi
      EOF
      ]
    }

    node {
      role = "worker"
    }

    node {
      role = "worker"
    }

  }
}

resource "kubernetes_namespace" "study" {
  depends_on = [kind_cluster.local_cluster]

  metadata {
    name = var.study_namespace
  }
}

resource "null_resource" "set_default_namespace" {
  depends_on = [
    kind_cluster.local_cluster,
    kubernetes_namespace.study
  ]

  provisioner "local-exec" {
    command = "kubectl config set-context kind-${var.kind_cluster_name} --namespace=${var.study_namespace}"
  }
}
