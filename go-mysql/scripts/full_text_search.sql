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
    FULLTEXT (title) WITH PARSER ngram
#     FULLTEXT (title, body)
#     FULLTEXT (title)
) ENGINE = InnoDB
  CHARACTER SET utf8mb4;

# title도 full text 가능하도록 fulltext index 추가
ALTER TABLE articles
    ADD FULLTEXT INDEX ft_title_index (title) WITH PARSER ngram;

# ALTER TABLE articles
#     ADD FULLTEXT INDEX ft_title_index (title);


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

########################################## STOPWORDS ##########################################
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


########################################## BOOLEAN SEARCH ##########################################
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('arc-125' IN NATURAL LANGUAGE MODE);


########################################## NGRAM ##########################################
# ngram token size
SHOW GLOBAL VARIABLES LIKE 'ngram_token_size';

# 이거 실행이 안되는 듯함
SET GLOBAL ngram_token_size = 1;

########################################## FULLL INDEX 확인 ##########################################

# to see how the ngram tokenizes the text
SET GLOBAL innodb_ft_aux_table = "gomysql/articles";

# todo: 실제로 아무것도 보여지지 않음 - 왜?
SELECT *
FROM information_schema.innodb_ft_index_cache
ORDER BY doc_id,
         position;


# FULL 인텍스 어떻게 처리가 되는지 확인해보기
SHOW VARIABLES LIKE 'innodb_optimize_fulltext_only';

SET GLOBAL innodb_optimize_fulltext_only = ON;

SET GLOBAL innodb_ft_aux_table = 'gomysql/articles';
SHOW VARIABLES LIKE 'innodb_ft_aux_table';

# todo : 동작을 하지 않음 - 왜?
SELECT word, doc_count, doc_id, position
FROM INFORMATION_SCHEMA.INNODB_FT_INDEX_TABLE
LIMIT 5;

# 이거는 잘 보임
SELECT *
FROM information_schema.innodb_ft_index_cache
ORDER BY doc_id, position;

########################################## STOPWORDS ##########################################
# stopword default table 조회
SELECT *
FROM INFORMATION_SCHEMA.INNODB_FT_DEFAULT_STOPWORD;

# custom stopwords 생성하기
-- Create a new stopword table
CREATE TABLE my_stopwords
(
    value VARCHAR(30)
) ENGINE = INNODB;

# -- Insert stopwords (for simplicity, a single stopword is used in this example)
INSERT INTO my_stopwords(value)
VALUES ('Ishmael');

SELECT *
FROM my_stopwords;

-- Set the innodb_ft_server_stopword_table option to the new stopword table
SET GLOBAL innodb_ft_server_stopword_table = 'gomysql/my_stopwords';

-- Create the full-text index (which rebuilds the table if no FTS_DOC_ID column is defined)
# CREATE FULLTEXT INDEX idx ON articles(title);
SET GLOBAL innodb_ft_aux_table = 'gomysql/articles';
SELECT *
FROM information_schema.innodb_ft_index_cache
ORDER BY doc_id, position;


########################################## NGRAM ##########################################
# https://gywn.net/2017/04/mysql_57-ngram-ft-se/

CREATE TABLE `articles`
(
    `id`   int(10) unsigned NOT NULL AUTO_INCREMENT,
    `body` text,
    PRIMARY KEY (`id`),
    FULLTEXT KEY `ftx` (`body`) WITH PARSER ngram
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

insert into articles (body)
values ('east');
insert into articles (body)
values ('east area');
insert into articles (body)
values ('east job');
insert into articles (body)
values ('eastnation');
insert into articles (body)
values ('eastway, try try');
insert into articles (body)
values ('try try');

SELECT *
FROM articles
WHERE MATCH(body) AGAINST('st' IN BOOLEAN MODE);

# 아래는 검색이 안되어야 하는데, 너무 잘됨
SELECT *
FROM articles
WHERE MATCH(body) AGAINST('ea' IN BOOLEAN MODE);
SELECT *
FROM articles
WHERE MATCH(body) AGAINST('eas' IN BOOLEAN MODE);

CREATE TABLE ngram_stopwords
(
    value VARCHAR(18)
) ENGINE = INNODB;

INSERT INTO ngram_stopwords VALUE ('east');
SELECT *
FROM ngram_stopwords;

SET GLOBAL innodb_ft_server_stopword_table = 'gomysql/ngram_stopwords';

alter table articles
    engine = innodb;

SELECT *
FROM articles
WHERE MATCH(body) AGAINST('ea' IN BOOLEAN MODE);

########################################## NGRAM2 ##########################################
CREATE TABLE IF NOT EXISTS `books`
(
    `id`     INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    `title`  VARCHAR(100)     NOT NULL,
    `author` VARCHAR(350)     NULL     DEFAULT NULL,
    `reg_at` TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FULLTEXT INDEX `ft_idx_title` (`title`) WITH PARSER `ngram`
);


INSERT INTO books(title, author)
VALUES ('철학은 어떻게 삶의 무기가 되는가', '야마구치 슈');
set global innodb_ft_aux_table = 'gomysql/books';

SELECT *
FROM INFORMATION_SCHEMA.INNODB_FT_INDEX_CACHE;

SHOW GLOBAL VARIABLES LIKE 'ngram_token_size';


SELECT *
FROM books;

INSERT INTO books(title, author)
VALUES ('abcdef', 'author1');
INSERT INTO books(title, author)
VALUES ('ab bc de ef', 'author2');

SELECT *
FROM books
WHERE MATCH(title) AGAINST('abc');

# a가 stopword에 있는 값이라서 index로 등록이 안되어 있음
SELECT *
FROM books
WHERE MATCH(title) AGAINST('ab');

# abcdef가 조회가 됨
# SELECT * FROM books WHERE MATCH(title) AGAINST ('arc-12' IN BOOLEAN MODE);

# 1256 -> 12 56을 검색하는 것과 같다
SELECT *
FROM books
WHERE MATCH(title) AGAINST('1256');

SELECT *
FROM books
WHERE MATCH(title) AGAINST('2326');
SELECT *
FROM books
WHERE MATCH(title) AGAINST('23 26 32');

SELECT *
FROM books
WHERE MATCH(title) AGAINST('23 26');


# stopword 적용하기
INSERT INTO books(title, author)
VALUES ('arc-2356332', 'author1');
INSERT INTO books(title, author)
VALUES ('arc-1256332', 'author2');
INSERT INTO books(title, author)
VALUES ('arc-2256332', 'author3');
INSERT INTO books(title, author)
VALUES ('arc-2316332', 'author4');
INSERT INTO books(title, author)
VALUES ('arc-2326332', 'author5');
INSERT INTO books(title, author)
VALUES ('arc-2336332', 'author6');
INSERT INTO books(title, author)
VALUES ('arc-1316332', 'author7');
INSERT INTO books(title, author)
VALUES ('arc-1326332', 'author8');
INSERT INTO books(title, author)
VALUES ('arc-1616332', 'author9');
INSERT INTO books(title, author)
VALUES ('arc-1656332', 'author10');


# 전체 변수 값 보기
SHOW VARIABLES;

# 전체 global 값 보기
SHOW GLOBAL VARIABLES;

# 현재 stop 관련된 옵션 보기
SHOW GLOBAL VARIABLES LIKE '%stop%';

# stopword 설정
SET GLOBAL innodb_ft_server_stopword_table = 'gomysql/ngram_stopwords';

# stopword 설정이후 index 가 잘 생성되는 거 확인
alter table books
    engine = innodb;

########################################## NGRAM3 ##########################################
# https://dev.mysql.com/blog-archive/innodb-full-text-n-gram-parser-ko/
CREATE TABLE articles
(
    FTS_DOC_ID BIGINT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
    title      VARCHAR(100),
    FULLTEXT INDEX ngram_idx (title) WITH PARSER ngram
) Engine = InnoDB
  CHARACTER SET utf8mb4;

INSERT INTO articles (title)
VALUES ('my sql');
INSERT INTO articles (title)
VALUES ('mysql');
INSERT INTO articles (title)
VALUES ('sq');
INSERT INTO articles (title)
VALUES ('sl');

SELECT *
FROM articles;

# full index: my, ql, sl, sq, ys
SET GLOBAL innodb_ft_aux_table = "gomysql/articles";


SELECT *
FROM INFORMATION_SCHEMA.INNODB_FT_INDEX_CACHE;

# NATURAL LANGUAGE MODE에서 검색되는 텍스트는 n-gram의 합집합으로 변환됨
# 예를 들어, ‘sql’는 ‘sq ql’로 변환

SELECT *
FROM articles
WHERE MATCH(title) AGAINST('sql');
# 결과:
# my sql
# mysql
# sq

# BOOLEAN MODE 에서 검색되는 텍스트는 n-gram 구문 검색으로 변환됩니다.
# 예를 들어, ‘sql’은 ‘ “sq ql”‘로 변환됩니다.
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('sql' IN BOOLEAN MODE);
# 결과:
# my sql
# mysql

########################################## NGRAM4 - NGRAM_TOKEN_SIZE=2 ##########################################
DROP TABLE IF EXISTS articles;
CREATE TABLE articles
(
    id    INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
    title VARCHAR(200),
    FULLTEXT INDEX ngram_idx (title) WITH PARSER ngram
) ENGINE = InnoDB
  CHARACTER SET utf8mb4;

# FULLTEXT index 생성
# ab bc -> ab bc
# abc -> ab bc
# ab bc de ef -> ab bc de ef
# abcdef -> ab bc cd de ef
# agh -> ag gh
# abc def -> ab bc de ef
# bcde -> bc cd de
INSERT INTO articles(title)
VALUES ('ab bc'),
       ('abc'),
       ('ab bc de ef'),
       ('abcdef'),
       ('bcde'),
       ('ab'),
       ('bc'),
       ('abcdef'),
       ('abc def');

# INSERT INTO articles(title)
# VALUES ('abc def');

SET GLOBAL innodb_ft_aux_table = "gomysql/articles";

SELECT *
FROM INFORMATION_SCHEMA.INNODB_FT_INDEX_CACHE;

# ngram Parser Term Search
# 1.natural language mode - the search term is converted to a union of ngram terms
# search term: abc -> ab bc 변환됨
# 조회 결과: ab, abc를 포함한 값을 반환한다
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('abc');

# 아래는 같은 결과
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('ab bc');

# 2.boolean mode - search term is converted to an ngram phrase search
# - the + and - operators indicate that a word must be present or absent, respectively, for a match to occur
# search term : abc -> ab bc

SELECT *
FROM articles
WHERE MATCH(title) AGAINST('abc' IN BOOLEAN MODE);

# search term : abc def -> ab bc de ef
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('abc def' IN BOOLEAN MODE);


# search term: +abc -def => +ab bc(는 있고), -de ef(는 없는)
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('+abc -def' IN BOOLEAN MODE);

# boolean vs natural language mode
# search term: abbc -> ab bc (ab bc 순서대로 모두 일치해야 하므로 검색 결과에 나오지 않음)
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('abbc' IN BOOLEAN MODE);

# search term: abbc -> ab bb bc 매칭이 되는 합집합의 결과를 노출한다
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('abbc');

# ngram Parser Wildcard Search
# - prefix term이 ngram_token_size보다 작은 경우
# search term: a* -> a로 시작하는 모든 rows
# - prefix term이 ngram_token_size보다 큰 경우, the prefix term is converted to an ngram phrase and the wildcard operator is ignored
# search term: abc* -> ab bc

# search term: a* -> a로 시작하는 모든 rows
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('a*' IN BOOLEAN  MODE);

# search term: abc* -> ab bc
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('abc*' IN BOOLEAN  MODE);


# ngram Parser Phrase Search <-- 이게 무슨 차이가 있는 건가?
# search term : abc -> ab bc
# search term : abc def -> ab bc de ef
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('abc def');

# index 다시 생성하기
alter table articles
    engine = innodb;

########################################## NGRAM - NGRAM_TOKEN_SIZE=1 ##########################################
# 어떻게 달라지나? 사용성이 더 높아지나?
#
