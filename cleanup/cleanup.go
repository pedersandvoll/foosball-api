package cleanup

import (
	"log"
	"pedersandvoll/foosballapi/config"
	"time"
)

type LobbyCleanupService struct {
	db             *config.Database
	checkInterval  time.Duration // How often to check for inactive lobbies
	inactiveWindow time.Duration // How long a lobby can be inactive before deletion
	stop           chan struct{} // Channel to signal shutdown
}

func NewLobbyCleanupService(db *config.Database, checkInterval, inactiveWindow time.Duration) *LobbyCleanupService {
	return &LobbyCleanupService{
		db:             db,
		checkInterval:  checkInterval,
		inactiveWindow: inactiveWindow,
		stop:           make(chan struct{}),
	}
}

func (s *LobbyCleanupService) Start() {
	go s.cleanupLoop()
}

func (s *LobbyCleanupService) Stop() {
	close(s.stop)
}

func (s *LobbyCleanupService) cleanupLoop() {
	ticker := time.NewTicker(s.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupInactiveLobbies()
		case <-s.stop:
			log.Println("Lobby cleanup service stopping")
			return
		}
	}
}

func (s *LobbyCleanupService) cleanupInactiveLobbies() {
	cutoffTime := time.Now().Add(-s.inactiveWindow)

	query := `DELETE FROM lobbies 
              WHERE (last_played IS NOT NULL AND last_played < $1)
              OR (last_played IS NULL AND created_at < $1)`

	result, err := s.db.Exec(query, cutoffTime)
	if err != nil {
		log.Printf("Error cleaning up inactive lobbies: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("Cleaned up %d inactive lobbies", rowsAffected)
	}
}
