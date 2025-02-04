ALTER TABLE organizations
DROP CONSTRAINT fk_activeseason;

ALTER TABLE organizations
DROP COLUMN activeseason;
