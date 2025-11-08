## redisinsight 설치 

helm dep update
helm upgrade --install redisinsight . -f values.yaml -n toolbox
