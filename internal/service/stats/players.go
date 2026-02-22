package stats

import (
	"fmt"
	"plugin/internal/repository/stats"
)

type PlayeStatsService struct {
	repo stats.PlayerStatsRepository
}

func NewPlayerStats(repo stats.PlayerStatsRepository) *PlayeStatsService {
	return &PlayeStatsService{repo: repo}
}

func (s *PlayeStatsService) Init(playerID int) error {
	stats, err := s.repo.Get(playerID)
	if err != nil {
		return err
	}
	if stats != nil {
		return nil
	}
	return s.repo.Create(playerID)
}

func (s *PlayeStatsService) GetStats(playerID int) (*stats.PlayerStats, error) {
	return s.repo.Get(playerID)
}

func (s *PlayeStatsService) Win(playerID int, wager int, payout int) error {
	if wager <= 0 || payout <= 0 {
		return fmt.Errorf("invalid wager or payout")
	}
	return s.repo.RecordWin(playerID, wager, payout)
}

func (s *PlayeStatsService) Loss(playerID int, wager int) error {
	if wager <= 0 {
		return fmt.Errorf("invalid wager")
	}
	return s.repo.RecordLoss(playerID, wager, wager)
}

func (s *PlayeStatsService) Reset(playerID int) error {
	return s.repo.Reset(playerID)
}
