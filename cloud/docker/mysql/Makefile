DB_CONTAINER_NAME := go-mysql
DB_NAME := gomysql
DB_USER_NAME := root
DB_PASSWORD := password
MYSQL_VERSION=8.0.31

mysql-create:
	@docker run --name $(DB_CONTAINER_NAME) \
	-e MYSQL_ROOT_PASSWORD=$(DB_PASSWORD) \
	-d -p 3306:3306 mysql:$(MYSQL_VERSION) --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

mysql-delete:
	@docker stop $(DB_CONTAINER_NAME)
	@docker rm $(DB_CONTAINER_NAME)

