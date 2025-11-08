# mysql 설치

mysql와 같은 서버는 app 단위로 두지 않고 개인용으로 사용할 거라서 database namespace에 설치하고 공용으로 사용한다

## storage PVC 설치

```bash
> kubectl apply -f mysql-storage.yaml -n database
```
## mysql 서버 설치

```bash
> helm install mysql bitnami/mysql -f values.yaml -n database

```
