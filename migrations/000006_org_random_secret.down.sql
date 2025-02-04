DROP TRIGGER IF EXISTS set_orgsecret ON organizations;

DROP FUNCTION IF EXISTS generate_unique_orgsecret;

ALTER TABLE organizations DROP CONSTRAINT IF EXISTS unique_orgsecret;
