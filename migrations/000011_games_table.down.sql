DROP INDEX IF EXISTS idx_team1_player1;
DROP INDEX IF EXISTS idx_team1_player2;
DROP INDEX IF EXISTS idx_team2_player1;
DROP INDEX IF EXISTS idx_team2_player2;
DROP INDEX IF EXISTS idx_games_lobbyid;

DROP TABLE IF EXISTS games;

DROP TYPE IF EXISTS game_status;

ALTER TABLE lobbyplayers DROP CONSTRAINT IF EXISTS unique_lobby_user;

ALTER TABLE lobbyplayers DROP COLUMN IF EXISTS playerid;
