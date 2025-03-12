--liquibase formatted sql
--changeset your.name:4
create table employee
(
    id       int primary key auto_increment not null,
    name     varchar(50) not null,
    address1 varchar(50),
    address2 varchar(50),
    city     varchar(30)
);

INSERT INTO liquibase_quickstart.employee (id, name, address1, address2, city)
VALUES (1, 'John Doe', '123 Main St', 'Apt 1', 'Beverly Hills');

INSERT INTO liquibase_quickstart.employee (id, name, address1, address2, city)
VALUES (2, 'John Doe2', '123 Main St', 'Apt 1', 'Beverly Hills');

ALTER TABLE liquibase_quickstart.employee ADD COLUMN email VARCHAR(50);

UPDATE liquibase_quickstart.employee SET email = 'user1@naver.com' WHERE id = 1;
UPDATE liquibase_quickstart.employee SET email = 'user2@naver.com' WHERE id = 2;


--comment rollback 실행하면 아래 전체가 한번에 실행이 된다. rollback 하는 단위는 changeset 단위로 되는 듯하다
--rollback UPDATE liquibase_quickstart.employee SET email = NULL WHERE id = 2;
--rollback UPDATE liquibase_quickstart.employee SET email = NULL WHERE id = 1;
--rollback ALTER TABLE liquibase_quickstart.employee DROP COLUMN email;
--rollback DELETE FROM liquibase_quickstart.employee WHERE id = 2;
--rollback DELETE FROM liquibase_quickstart.employee WHERE id = 1;
--rollback DROP TABLE employee;
