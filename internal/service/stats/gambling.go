package stats

import (
	"fmt"
	"plugin/internal/repository/stats"
)

type GamblingStatsService struct {
	repo stats.GamblingStatsRepository
}

func NewGamblingStats(repo stats.GamblingStatsRepository) *GamblingStatsService {
	return &GamblingStatsService{repo: repo}
}

func (s *GamblingStatsService) Init() error {
	stats, err := s.repo.Get()
	if err != nil {
		return err
	}

	if stats != nil {
		return nil
	}

	return s.repo.Init()
}

func (s *GamblingStatsService) GetStats() (*stats.GamblingStats, error) {
	return s.repo.Get()
}

func (s *GamblingStatsService) RecordGamble(wager int, paid int) error {
	if wager <= 0 {
		return fmt.Errorf("wager must be positive")
	}
	if paid < 0 {
		return fmt.Errorf("paid cannot be negative")
	}
	return s.repo.UpdateAfterGamble(wager, paid)
}

func (s *GamblingStatsService) Reset() error {
	return s.repo.Reset()
}
