package stats

import (
	"database/sql"
	"errors"
	"plugin/internal/database/queries"
	"time"
)

type GamblingStats struct {
	TotalGambles int
	TotalWagered int
	TotalPaid    int
	LastUpdate   *time.Time
}

type GamblingStatsRepository interface {
	Init() error
	Get() (*GamblingStats, error)
	UpdateAfterGamble(wager int, paid int) error
	Reset() error
}

type gamblingRepository struct {
	db *sql.DB
}

func NewGamblingStats(db *sql.DB) GamblingStatsRepository {
	return &gamblingRepository{db: db}
}

func (r *gamblingRepository) Init() error {
	_, err := r.db.Exec(queries.InitGlobalStats)
	return err
}

func (r *gamblingRepository) Get() (*GamblingStats, error) {
	var stats GamblingStats

	err := r.db.QueryRow(queries.GetGlobalStats).Scan(
		&stats.TotalGambles,
		&stats.TotalWagered,
		&stats.TotalPaid,
		&stats.LastUpdate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &stats, nil
}

func (r *gamblingRepository) UpdateAfterGamble(wager int, paid int) error {
	_, err := r.db.Exec(
		queries.UpdateGlobalAfterGamble,
		wager,
		paid,
	)
	return err
}

func (r *gamblingRepository) Reset() error {
	_, err := r.db.Exec(queries.ResetGlobalStats)
	return err
}
