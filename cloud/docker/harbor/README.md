# Install
> > helm repo add harbor https://helm.goharbor.io
> helm upgrade --install harbor -f values.yaml harbor/harbor --create-namespace -n toolbox

# TODO
- 설치까지는 잘 되지만, harbor 사이트에 로그인이 안되는 이슈가 있음
  - NodePort로 설정하는 경우에는 로그인이 안되는 이슈가 있다고는 검색이 되고 원인 분석이 더 필요한 상황  
