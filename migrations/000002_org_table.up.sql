CREATE TABLE organizations (
    orgid SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    orgsecret TEXT NOT NULL,
    orgowner INT NOT NULL,
    createdate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_orgowner FOREIGN KEY (orgowner) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE INDEX idx_orgowner ON organizations(orgowner);
