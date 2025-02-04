CREATE TABLE seasons (
    seasonid SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    orgid INT NOT NULL,
    CONSTRAINT fk_orgid FOREIGN KEY (orgid) REFERENCES organizations(orgid) ON DELETE CASCADE
);

CREATE INDEX idx_seasons_orgid ON seasons(orgid);
