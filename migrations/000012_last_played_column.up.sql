ALTER TABLE lobbies
ADD COLUMN last_played TIMESTAMP WITH TIME ZONE;

CREATE OR REPLACE FUNCTION update_lobby_last_played()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE lobbies
    SET last_played = NEW.createdat
    WHERE lobbyid = NEW.lobbyid;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_game_insert
AFTER INSERT ON games
FOR EACH ROW
EXECUTE FUNCTION update_lobby_last_played();
