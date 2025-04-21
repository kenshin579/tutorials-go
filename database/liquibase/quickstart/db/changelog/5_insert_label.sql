--liquibase formatted sql

--changeset frank.oh:5 context:"@dev and @pre"
INSERT INTO liquibase_quickstart.person (id, name, address1, address2, city) VALUES (2, 'Name1', 'Address1', 'Apt 1', 'Beverly Hills');

--rollback DELETE FROM liquibase_quickstart.person WHERE id = 2;
