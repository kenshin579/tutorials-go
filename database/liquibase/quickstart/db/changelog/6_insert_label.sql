--liquibase formatted sql

--changeset frank.oh:6 context:dev labels:insert3
INSERT INTO liquibase_quickstart.person (id, name, address1, address2, city) VALUES (3, 'Name2', 'Address2', 'Apt 2', 'Beverly Hills');

--rollback DELETE FROM liquibase_quickstart.person WHERE id = 3;
