package stats

import (
	"database/sql"
	"errors"
	"plugin/internal/database/queries"
	"time"
)

type PlayerStats struct {
	PlayerID     int
	TotalGambles int
	TotalWagered int
	TotalWon     int
	TotalLost    int
	Wins         int
	Losses       int
	LastGamble   *time.Time
}

type PlayerStatsRepository interface {
	Create(playerID int) error
	Get(playerID int) (*PlayerStats, error)
	RecordWin(playerID int, wager int, won int) error
	RecordLoss(playerID int, wager int, lost int) error
	Reset(playerID int) error
}

type playerRepository struct {
	db *sql.DB
}

func NewPlayerStats(db *sql.DB) PlayerStatsRepository {
	return &playerRepository{db: db}
}

func (r *playerRepository) Create(playerID int) error {
	_, err := r.db.Exec(queries.CreatePlayerStats, playerID)
	return err
}

func (r *playerRepository) Get(playerID int) (*PlayerStats, error) {
	var stats PlayerStats
	var lastGamble sql.NullTime

	err := r.db.QueryRow(queries.GetPlayerStats, playerID).Scan(
		&stats.PlayerID,
		&stats.TotalGambles,
		&stats.TotalWagered,
		&stats.TotalWon,
		&stats.TotalLost,
		&stats.Wins,
		&stats.Losses,
		&lastGamble,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if lastGamble.Valid {
		stats.LastGamble = &lastGamble.Time
	}

	return &stats, nil
}

func (r *playerRepository) RecordWin(playerID int, wager int, won int) error {
	_, err := r.db.Exec(
		queries.UpdateStatsWin,
		wager,
		won,
		playerID,
	)
	return err
}

func (r *playerRepository) RecordLoss(playerID int, wager int, lost int) error {
	_, err := r.db.Exec(
		queries.UpdateStatsLoss,
		wager,
		lost,
		playerID,
	)
	return err
}

func (r *playerRepository) Reset(playerID int) error {
	_, err := r.db.Exec(queries.ResetPlayerStats, playerID)
	return err
}
