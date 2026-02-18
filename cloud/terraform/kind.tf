# =============================================================================
# Kind 클러스터 및 Kubernetes 리소스 정의
# - Kind(Kubernetes IN Docker): 로컬 개발용 K8s 클러스터
# - Resource: Terraform이 생성/관리하는 인프라 단위
# =============================================================================

# Kind 클러스터 리소스
# resource "<provider>_<type>" "<이름>" 형식으로 선언한다
resource "kind_cluster" "local_cluster" {
  name           = var.kind_cluster_name     # Variable 참조: var.<변수명>
  wait_for_ready = true                      # 클러스터가 Ready 상태가 될 때까지 대기
  node_image     = "kindest/node:v1.28.15"   # 사용할 Kubernetes 버전

  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"

    # Control Plane 노드: 클러스터를 관리하는 마스터 노드
    node {
      role = "control-plane"

      # 호스트 포트 매핑: 클러스터 외부에서 NodePort 서비스에 접근 가능하게 한다
      extra_port_mappings {
        container_port = 30080
        host_port      = 30080
        listen_address = "127.0.0.1"
      }

      # 로컬 스토리지 마운트: PV(Persistent Volume) 데이터를 호스트에 저장
      extra_mounts {
        host_path      = "/tmp/kind-storage"
        container_path = "/opt/local-path-provisioner"
      }

      # kubelet 리소스 예약 설정
      kubeadm_config_patches = [<<-EOF
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            system-reserved: cpu=250m,memory=250Mi
      EOF
      ]
    }

    # Worker 노드 2개: 실제 Pod가 스케줄링되는 노드
    node {
      role = "worker"
    }

    node {
      role = "worker"
    }

  }
}

# Namespace 생성
# - Namespace: K8s에서 리소스를 논리적으로 격리하는 단위
resource "kubernetes_namespace" "study" {
  depends_on = [kind_cluster.local_cluster]  # 클러스터 생성 후에 실행

  metadata {
    name = var.study_namespace
  }
}

# kubectl 기본 Namespace 설정
# - null_resource: 실제 인프라가 아닌 보조 작업에 사용
# - local-exec: Terraform이 실행되는 로컬 머신에서 명령어 실행
resource "null_resource" "set_default_namespace" {
  depends_on = [
    kind_cluster.local_cluster,
    kubernetes_namespace.study
  ]

  provisioner "local-exec" {
    command = "kubectl config set-context kind-${var.kind_cluster_name} --namespace=${var.study_namespace}"
  }
}
