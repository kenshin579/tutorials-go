# https://www.mysqltutorial.org/mysql-ngram-full-text-parser/
# https://jinhokwon.github.io/devops/mysql/mysql8-docker/

CREATE DATABASE gomysql default CHARACTER SET UTF8MB4;

use gomysql;

DROP TABLE IF EXISTS articles;
CREATE TABLE articles
(
    id    INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
    title VARCHAR(200),
    body  TEXT,
#     FULLTEXT (title, body) WITH PARSER ngram
#     FULLTEXT (title, body)
    FULLTEXT (title)
) ENGINE = InnoDB
  CHARACTER SET utf8mb4;

# title도 full text 가능하도록 fulltext index 추가
# ALTER TABLE articles
#     ADD FULLTEXT INDEX ft_title_index (title) WITH PARSER ngram;

ALTER TABLE articles
    ADD FULLTEXT INDEX ft_title_index (title);


# use the SET NAMES statement sets the character set to utf8mb4.
SET NAMES utf8mb4;

INSERT INTO articles(title, body)
VALUES ('arc-2356332', 'body1');
INSERT INTO articles(title, body)
VALUES ('arc-1256332', 'body2');
INSERT INTO articles(title, body)
VALUES ('arc-2256332', 'body3');
INSERT INTO articles(title, body)
VALUES ('arc-2316332', 'body4');
INSERT INTO articles(title, body)
VALUES ('arc-2326332', 'body5');
INSERT INTO articles(title, body)
VALUES ('arc-2336332', 'body6');
INSERT INTO articles(title, body)
VALUES ('arc-1316332', 'body7');
INSERT INTO articles(title, body)
VALUES ('arc-1326332', 'body8');
INSERT INTO articles(title, body)
VALUES ('arc-1616332', 'body9');
INSERT INTO articles(title, body)
VALUES ('arc-1656332', 'body10');

INSERT INTO articles (title, body)
VALUES ('MySQL Tutorial', 'DBMS stands for DataBase ...'),
       ('How To Use MySQL Well', 'After you went through a ...'),
       ('Optimizing MySQL', 'In this tutorial, we show ...'),
       ('1001 MySQL Tricks', '1. Never run mysqld as root. 2. ...'),
       ('MySQL vs. YourSQL', 'In the following database comparison ...'),
       ('MySQL Security', 'When configured properly, MySQL ...');

# full-text search
# https://gngsn.tistory.com/162
# matching - 언제 표시되는지 잘 모르겠음
# SELECT * FROM articles WHERE MATCH (title) AGAINST ('title' IN NATURAL LANGUAGE MODE);
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('arc-125' IN NATURAL LANGUAGE MODE);

# relevance value을 직접 뽑아보는 방법
SELECT id,
       MATCH(title)
             AGAINST('arc-125' IN NATURAL LANGUAGE MODE) AS score
FROM articles;

# 실제 값 + relevance 갑솓 같이 보고 싶을 때
SELECT id,
       title,
       MATCH(title) AGAINST('arc-125' IN NATURAL LANGUAGE MODE) AS score
FROM articles
WHERE MATCH(title) AGAINST('arc-125' IN NATURAL LANGUAGE MODE);
#
# SELECT id,
#        body,
#        MATCH(title, body) AGAINST('Security implications of running MySQL as root'
#                                   IN NATURAL LANGUAGE MODE) AS score
# FROM articles
# WHERE MATCH(title, body) AGAINST('Security implications of running MySQL as root'
#                                  IN NATURAL LANGUAGE MODE);

# title, body도 같이 검색하기
SELECT *
FROM articles
WHERE MATCH(title, body) AGAINST('database' IN NATURAL LANGUAGE MODE);

# IN NATURAL LANGUAGE MODE은 기본 값이라서 아래 실행하면 위와 같은 결과를 얻음
SELECT *
FROM articles
WHERE MATCH(title, body) AGAINST('database');


# to see how the ngram tokenizes the text
SET GLOBAL innodb_ft_aux_table = "test/articles";

# todo: 실제로 아무것도 보여지지 않음 - 왜?
SELECT *
FROM information_schema.innodb_ft_index_cache
ORDER BY doc_id,
         position;

# ngram token size
SHOW GLOBAL VARIABLES LIKE 'ngram_token_size';


# stopword 조회
SELECT *
FROM INFORMATION_SCHEMA.INNODB_FT_DEFAULT_STOPWORD;

# custom stopword를 만들기

