ALTER TABLE organizations
ADD COLUMN activeseason INT;

ALTER TABLE organizations
ADD CONSTRAINT fk_activeseason
FOREIGN KEY (activeseason)
REFERENCES seasons(seasonid)
ON DELETE SET NULL;
