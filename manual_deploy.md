kubectl config set-cluster costs-n-tasks-cluster --server="https://146.148.56.200" --insecure-skip-tls-verify=true
kubectl config set-credentials cloud_okteto_com-user --token="eyJhbGciOiJSUzI1NiIsImtpZCI6IlBJWWEyU3FyUm5LcGlhNTRDZGU3U2tDRkpRSDVicThCcC1qbFFnN2JtNjgifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJva3RldG8iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoiOTI5ZjFlNmYtMjZiMS00Yjk0LTljYmEtNWExYTgyM2U0NjIyLXRva2VuLXc3ejU0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6IjkyOWYxZTZmLTI2YjEtNGI5NC05Y2JhLTVhMWE4MjNlNDYyMiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImQwZjU2ZTk0LTFkYmYtNGRlMS05ZjljLTQ3ZGJiNjNlNGU0MyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpva3RldG86OTI5ZjFlNmYtMjZiMS00Yjk0LTljYmEtNWExYTgyM2U0NjIyIn0.p0lTt21winA9pTXou2eGm5ggvlnZs5Wycx2cOyE7pZxnKu0dgt8ZAlWenqkgJij7OHZ9gJlYUzHqnvQS5H5JnB9WXjfaK9krDazxgGZcoIGlLULTJk6Ujyj5uIjd6PATJiF5SCVICqH2Elb5WItxmV5JPON-fL1LwRDFIg9n-ID2Rs0VS0Cj0bVCPP7KhgY9jzfjY6javKpMxz4SBBrfPGCMS7OViM__Nif93UmJ8LDbPLNY0p4glMnaEGxU52Zia_HZW5ZiwmArkKc2w3XCA99bZpvp-pMvH9APAIo7qsguAqL3oouMAN4-xIVCJUMQGQIn2fcP33n4VTGUuJGAKw"
kubectl config set-context default --cluster=costs-n-tasks-cluster --user=cloud_okteto_com-user --namespace="drstarland"
kubectl config use-context default
kubectl get pods
cd k8s && helm upgrade --install --force --debug services services-chart/
kubectl apply -f load-balancer.yaml
sleep 45
kubectl get pods
