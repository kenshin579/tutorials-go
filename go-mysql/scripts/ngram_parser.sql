CREATE DATABASE gomysql default CHARACTER SET UTF8MB4;

use gomysql;

########################################## 변수 ##########################################

# 전체 변수 값 보기
SHOW VARIABLES;

# 전체 global 값 보기
SHOW GLOBAL VARIABLES;

# 현재 stop 관련된 옵션 보기
SHOW GLOBAL VARIABLES LIKE '%stop%';

# ngram token size
SHOW GLOBAL VARIABLES LIKE 'ngram_token_size';

########################################## STOPWORDS ##########################################
# default STOPWORD 목록 보기
SELECT *
FROM INFORMATION_SCHEMA.INNODB_FT_DEFAULT_STOPWORD;

# ngram stopword table 생성하기
CREATE TABLE ngram_stopwords
(
    value VARCHAR(30)
) ENGINE = INNODB;

# INSERT INTO ngram_stopwords VALUE ('a');

SELECT *
FROM ngram_stopwords;

SET GLOBAL innodb_ft_server_stopword_table = 'gomysql/ngram_stopwords';

########################################## NGRAM PARSER - NGRAM_TOKEN_SIZE=2 ##########################################
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

# full-index 생성되는 거 확인 해보기
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
WHERE MATCH(title) AGAINST('a*' IN BOOLEAN MODE);

# search term: abc* -> ab bc
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('abc*' IN BOOLEAN MODE);


# ngram Parser Phrase Search <-- 이게 무슨 차이가 있는 건가?
# search term : abc -> ab bc
# search term : abc def -> ab bc de ef
SELECT *
FROM articles
WHERE MATCH(title) AGAINST('abc def');

# index 다시 생성하기
alter table articles
    engine = innodb;

# 참고
# https://gywn.net/2017/04/mysql_57-ngram-ft-se/
# https://www.mysqltutorial.org/mysql-ngram-full-text-parser/
# https://jinhokwon.github.io/devops/mysql/mysql8-docker/
