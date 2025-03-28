--liquibase formatted sql

--changeset your.name:2
INSERT INTO liquibase_quickstart.person (id, name, address1, address2, city) VALUES (1, 'John Doe', '123 Main St', 'Apt 1', 'Beverly Hills');
