variable "study_namespace" {
  type    = string
}

# argocd 암호는 argocd 명령어로 생성한다
# argocd account bcrypt --password 'password'
variable "argocd_password" {
  type    = string
  default = "$2a$10$UfwTWJDvQ7e.ed6wBDVVxeoUlk9R0HEfXOEu1PqUfxlomAV46CIze"
}

variable "ingress-nginx_namespace" {
  type    = string
  default = "ingress-nginx"
}
