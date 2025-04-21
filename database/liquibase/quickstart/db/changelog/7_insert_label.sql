--liquibase formatted sql

--changeset frank.oh:7 context:"@beta and @post"
INSERT INTO liquibase_quickstart.person (id, name, address1, address2, city) VALUES (4, 'Name3', 'Address3', 'Apt 3', 'Beverly Hills');

--rollback DELETE FROM liquibase_quickstart.person WHERE id = 3;
