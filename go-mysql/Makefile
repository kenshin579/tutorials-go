DB_CONTAINER_NAME := mysql-test
DB_NAME := mysql-test
DB_USER_NAME := root
DB_PASSWORD := mysql12
MYSQL_VERSION=8.0.28

mysql-create:
	@docker run --name $(DB_CONTAINER_NAME) \
	-e MYSQL_ROOT_PASSWORD=$(DB_PASSWORD) \
	-d -p 3306:3306 mysql:$(MYSQL_VERSION) --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

mysql-delete:
	@docker stop $(DB_CONTAINER_NAME)
	@docker rm $(DB_CONTAINER_NAME)
