output "kubeconfig_path" {
  description = "KIND 클러스터의 kubeconfig 파일 경로"
  value       = kind_cluster.local_cluster.kubeconfig_path
}
