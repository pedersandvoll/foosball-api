DROP TRIGGER IF EXISTS after_game_insert ON games;

DROP FUNCTION IF EXISTS update_lobby_last_played();

ALTER TABLE lobbies
DROP COLUMN IF EXISTS last_played;
