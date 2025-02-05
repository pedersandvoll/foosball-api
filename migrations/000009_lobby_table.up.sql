CREATE TYPE lobby_status AS ENUM ('not_in_game', 'in_game');

CREATE TABLE lobbies (
    lobbyid SERIAL PRIMARY KEY,
    orgid INT NOT NULL,
    seasonid INT NOT NULL,
    createdby INT NOT NULL,
    status lobby_status NOT NULL DEFAULT 'not_in_game',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_orgid FOREIGN KEY (orgid) REFERENCES organizations(orgid) ON DELETE CASCADE,
    CONSTRAINT fk_seasonid FOREIGN KEY (seasonid) REFERENCES seasons(seasonid) ON DELETE CASCADE,
    CONSTRAINT fk_createdby FOREIGN KEY (createdby) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE INDEX idx_orgid ON organizations(orgid);
CREATE INDEX idx_seasonid ON seasons(seasonid);
CREATE INDEX idx_createdby ON users(userid);
