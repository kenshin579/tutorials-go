# https://www.mysqltutorial.org/mysql-ngram-full-text-parser/
# https://jinhokwon.github.io/devops/mysql/mysql8-docker/

CREATE DATABASE gomysql default CHARACTER SET UTF8MB4;

DROP TABLE IF EXISTS posts;

use gomysql;
CREATE TABLE posts
(
    id    INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255),
    body  TEXT,
    FULLTEXT (title, body) WITH PARSER NGRAM
) ENGINE = INNODB
  CHARACTER SET utf8mb4;


# use the SET NAMES statement sets the character set to utf8mb4.
SET NAMES utf8mb4;

INSERT INTO posts(title, body)
VALUES ('title-2356332', 'body1');
INSERT INTO posts(title, body)
VALUES ('title-1256332', 'body2');
INSERT INTO posts(title, body)
VALUES ('title-2256332', 'body3');
INSERT INTO posts(title, body)
VALUES ('title-2316332', 'body4');
INSERT INTO posts(title, body)
VALUES ('title-2326332', 'body5');
INSERT INTO posts(title, body)
VALUES ('title-2336332', 'body6');
INSERT INTO posts(title, body)
VALUES ('title-1316332', 'body7');
INSERT INTO posts(title, body)
VALUES ('title-1326332', 'body8');
INSERT INTO posts(title, body)
VALUES ('title-1616332', 'body9');
INSERT INTO posts(title, body)
VALUES ('title-1656332', 'body10');

# to see how the ngram tokenizes the text
SET GLOBAL innodb_ft_aux_table = "test/posts";

SELECT *
FROM information_schema.innodb_ft_index_cache
ORDER BY doc_id,
         position;

ALTER TABLE posts ADD FULLTEXT INDEX ft_title_index (title) WITH PARSER ngram;

# ngram token size
SHOW GLOBAL VARIABLES LIKE 'ngram_token_size';

# https://gngsn.tistory.com/162
# matching - 언제 표시되는지 잘 모르겠음
# SELECT * FROM posts WHERE MATCH (title) AGAINST ('title' IN NATURAL LANGUAGE MODE);
SELECT * FROM posts WHERE MATCH (title) AGAINST ('titl' IN NATURAL LANGUAGE MODE);

# CREATE TABLE IF NOT EXISTS `books`
# (
#     `id`     INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
#     `title`  VARCHAR(100)     NOT NULL,
#     `author` VARCHAR(350)     NULL     DEFAULT NULL,
#     `reg_at` TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
#     PRIMARY KEY (`id`),
#     FULLTEXT INDEX `ft_idx_title` (`title`) WITH PARSER `ngram`
# );
#
# INSERT INTO books(title, author) VALUES('철학은 어떻게 삶의 무기가 되는가','야마구치 슈');

# set global innodb_ft_aux_table = 'test/books';
# SELECT * FROM INFORMATION_SCHEMA.INNODB_FT_INDEX_CACHE;

# stopword 조회
SELECT * FROM INFORMATION_SCHEMA.INNODB_FT_DEFAULT_STOPWORD;
