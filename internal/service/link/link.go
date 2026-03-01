package link

import (
	"plugin/internal/repository/link"
	"time"
)

type Service struct {
	repo link.LinkRepository
}

func New(repo link.LinkRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateLink(playerID int, code string) error {
	expires := time.Now().Add(2 * time.Minute)
	return s.repo.Create(playerID, code, expires)
}

func (s *Service) GetPlayerIDByCode(code string) (int, error) {
	return s.repo.GetPlayerIDByCode(code)
}

func (s *Service) GetCodeByPlayerID(playerID int) (string, error) {
	return s.repo.GetCodeByPlayerID(playerID)
}

func (s *Service) DeleteByPlayerID(playerID int) error {
	return s.repo.DeleteByPlayerID(playerID)
}

func (s *Service) DeleteByCode(code string) error {
	return s.repo.DeleteByCode(code)
}
