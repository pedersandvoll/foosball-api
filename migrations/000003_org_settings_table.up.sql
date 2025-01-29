CREATE TABLE organizationSettings (
    orgid INT NOT NULL UNIQUE,
    maxgamesperseason INT DEFAULT 1000,
    maxlobbies INT DEFAULT 1,
    orgowner INT NOT NULL, 
    CONSTRAINT fk_orgid FOREIGN KEY (orgid) REFERENCES organizations(orgid) ON DELETE CASCADE,
    CONSTRAINT fk_orgowner FOREIGN KEY (orgowner) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE INDEX idx_orgsettings_orgid ON organizationSettings(orgid);
CREATE INDEX idx_orgsettings_orgowner ON organizationSettings(orgowner);
