--liquibase formatted sql

--changeset your.name:3
ALTER TABLE liquibase_quickstart.person ADD COLUMN zip VARCHAR(10);

UPDATE liquibase_quickstart.person SET zip = '90210' WHERE id = 1;

--rollback ALTER TABLE liquibase_quickstart.person DROP COLUMN zip;

