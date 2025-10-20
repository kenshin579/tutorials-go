-- 사용자 생성
CREATE USER testuser WITH PASSWORD 'test1234';

-- 데이터베이스 생성
CREATE DATABASE testdb OWNER testuser;

-- testuser에게 권한 부여
GRANT ALL PRIVILEGES ON DATABASE testdb TO testuser;

-- drop table users;

CREATE TABLE users
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100),
    age  INT
);

INSERT INTO users (name, age)
VALUES ('Alice', 25),
       ('Bob', 30),
       ('Charlie', 28);
