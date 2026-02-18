# =============================================================================
# 변수 정의 (Variable)
# - 설정값을 외부에서 주입할 수 있게 해준다
# - default 값이 있으면 생략 가능, 없으면 apply 시 입력 필요
# - terraform apply -var="study_namespace=dev" 처럼 오버라이드 가능
# =============================================================================

variable "study_namespace" {
  description = "실습용 Kubernetes Namespace 이름"
  type        = string
  default     = "study"
}

variable "kind_cluster_name" {
  description = "Kind 클러스터 이름 (kubectl context에도 사용됨)"
  type        = string
  default     = "terraform-study-cluster"
}
