# =============================================================================
# 인프라 모듈 변수
# - 모듈 외부에서 값을 전달받을 수 있는 입력 변수들
# =============================================================================

variable "study_namespace" {
  description = "실습용 Namespace 이름 (상위 모듈에서 전달)"
  type        = string
}

# ArgoCD 관리자 비밀번호 (bcrypt 해시)
# 생성 방법: argocd account bcrypt --password 'password'
variable "argocd_password" {
  description = "ArgoCD admin 비밀번호 (bcrypt 해시값)"
  type        = string
  default     = "$2a$10$UfwTWJDvQ7e.ed6wBDVVxeoUlk9R0HEfXOEu1PqUfxlomAV46CIze"
}

variable "ingress-nginx_namespace" {
  description = "Ingress-NGINX가 설치될 Namespace"
  type        = string
  default     = "ingress-nginx"
}
