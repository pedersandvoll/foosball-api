ALTER TABLE organizations ADD CONSTRAINT unique_orgsecret UNIQUE (orgsecret);

CREATE OR REPLACE FUNCTION generate_unique_orgsecret() 
RETURNS TRIGGER AS $$
DECLARE
    new_secret TEXT;
BEGIN
    LOOP
        new_secret := FLOOR(1000 + (RANDOM() * 9000))::TEXT;
        
        IF NOT EXISTS (SELECT 1 FROM organizations WHERE orgsecret = new_secret) THEN
            EXIT;
        END IF;
    END LOOP;
    
    NEW.orgsecret := new_secret;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_orgsecret
BEFORE INSERT ON organizations
FOR EACH ROW
EXECUTE FUNCTION generate_unique_orgsecret();
