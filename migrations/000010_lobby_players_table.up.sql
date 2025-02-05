CREATE TABLE lobbyplayers (
    lobbyid INT NOT NULL,
    userid INT NOT NULL,
    gamesplayed INT DEFAULT 0,
    
    CONSTRAINT fk_lobbyid FOREIGN KEY (lobbyid) REFERENCES lobbies(lobbyid) ON DELETE CASCADE,
    CONSTRAINT fk_userid FOREIGN KEY (userid) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE INDEX idx_lobbyid ON lobbyplayers(lobbyid);
CREATE INDEX idx_userid ON lobbyplayers(userid);
