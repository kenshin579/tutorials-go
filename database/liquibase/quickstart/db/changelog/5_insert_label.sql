--liquibase formatted sql

--comment: labels도 포함해서 그냥 다 update가 되어 있음 <- todo: 확인 필요
--changeset frank.oh:5 labels:insert2
INSERT INTO liquibase_quickstart.person (id, name, address1, address2, city) VALUES (2, 'Name1', 'Address1', 'Apt 1', 'Beverly Hills');

--rollback DELETE FROM liquibase_quickstart.person WHERE id = 2;
