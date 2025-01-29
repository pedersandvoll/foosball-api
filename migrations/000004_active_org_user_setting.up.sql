ALTER TABLE users
ADD COLUMN activeorg INT;

ALTER TABLE users
ADD CONSTRAINT fk_activeorg
FOREIGN KEY (activeorg)
REFERENCES organizations(orgi)
ON DELETE SET NULL;
