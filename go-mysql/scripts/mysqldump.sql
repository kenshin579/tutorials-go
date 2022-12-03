CREATE DATABASE mysqlbackup default CHARACTER SET UTF8MB4;

use mysqlbackup;

DROP TABLE IF EXISTS articles;
CREATE TABLE articles
(
    id    INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
    title VARCHAR(200),
    FULLTEXT INDEX ngram_idx (title) WITH PARSER ngram
) ENGINE = InnoDB
  CHARACTER SET utf8mb4;

INSERT INTO articles(title)
VALUES ('title1'),
       ('title2'),
       ('title3');

SELECT * FROM articles;
