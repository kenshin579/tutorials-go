# =============================================================================
# 출력값 정의 (Output)
# - terraform apply 완료 후 중요한 정보를 출력한다
# - 다른 모듈에서 참조하거나, 사용자에게 정보를 전달할 때 사용
# - terraform output 명령어로 언제든 다시 확인 가능
# =============================================================================

output "kubeconfig_path" {
  description = "Kind 클러스터의 kubeconfig 파일 경로"
  value       = kind_cluster.local_cluster.kubeconfig_path
}
