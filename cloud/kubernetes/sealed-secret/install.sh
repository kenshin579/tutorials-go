
# secret 생성하기
kubectl -n default create secret generic my-secret \
  --from-literal=username=admin \
  --from-literal=password=hello-world \
  --dry-run=client \
  -o yaml > my-secret.yaml

# sealed secret로 변환하기
kubeseal --format=yaml < my-secret.yaml > my-sealed-secret.yaml

# sealed-secret를 클러스터에 등록하기


# 백업하기
kubectl get secret -n kube-system -l sealedsecrets.bitnami.com/sealed-secrets-key -o yaml >master.key

#복구하기
kubectl apply -f master.key
kubectl delete pod -n kube-system -l name=sealed-secrets-controller
