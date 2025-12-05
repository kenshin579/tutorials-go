output "cluster_name" {
  description = "Kind cluster name"
  value       = kind_cluster.default.name
}

output "cluster_endpoint" {
  description = "Kubernetes API endpoint"
  value       = kind_cluster.default.endpoint
}

output "kubeconfig" {
  description = "Kubeconfig for the cluster"
  value       = kind_cluster.default.kubeconfig
  sensitive   = true
}
