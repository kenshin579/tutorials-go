curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"message":"My first webhook"}' \
    http://localhost:12000/example

kubectl -n argo-events port-forward pod/eventsource-controller-54f89b495b-xrgn7 12000:12000
