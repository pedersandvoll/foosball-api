ALTER TABLE lobbyplayers ADD COLUMN playerid SERIAL PRIMARY KEY;

ALTER TABLE lobbyplayers ADD CONSTRAINT unique_lobby_user UNIQUE (lobbyid, userid);

CREATE TYPE game_status AS ENUM ('pending', 'in_progress', 'completed', 'canceled');

CREATE TABLE games (
    gameid SERIAL PRIMARY KEY,
    lobbyid INT NOT NULL,
    team1_player1 INT NOT NULL,
    team1_player2 INT NOT NULL,
    team2_player1 INT NOT NULL,
    team2_player2 INT NOT NULL,
    team1_score INT DEFAULT 0,
    team2_score INT DEFAULT 0,
    status game_status DEFAULT 'pending',
    last_played TIMESTAMP WITH TIME ZONE,
    createdat TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_lobbyid FOREIGN KEY (lobbyid) REFERENCES lobbies(lobbyid) ON DELETE CASCADE,
    CONSTRAINT fk_team1_player1 FOREIGN KEY (team1_player1) REFERENCES lobbyplayers(playerid) ON DELETE RESTRICT,
    CONSTRAINT fk_team1_player2 FOREIGN KEY (team1_player2) REFERENCES lobbyplayers(playerid) ON DELETE RESTRICT,
    CONSTRAINT fk_team2_player1 FOREIGN KEY (team2_player1) REFERENCES lobbyplayers(playerid) ON DELETE RESTRICT,
    CONSTRAINT fk_team2_player2 FOREIGN KEY (team2_player2) REFERENCES lobbyplayers(playerid) ON DELETE RESTRICT
);

CREATE INDEX idx_games_lobbyid ON games(lobbyid);
CREATE INDEX idx_team1_player1 ON games(team1_player1);
CREATE INDEX idx_team1_player2 ON games(team1_player2);
CREATE INDEX idx_team2_player1 ON games(team2_player1);
CREATE INDEX idx_team2_player2 ON games(team2_player2);
