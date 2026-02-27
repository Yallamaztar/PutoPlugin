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

func (s *Service) CreateLink(id int, code string) error {
	expires := time.Now().Add(2 * time.Minute)
	return s.repo.Create(id, code, expires)
}

func (s *Service) GetPlayerIDByCode(code string) (int, error) {
	return s.repo.GetPlayerIDByCode(code)
}

func (s *Service) GetCodeByPlayerID(id int) (string, error) {
	return s.repo.GetCodeByPlayerID(id)
}

func (s *Service) DeleteByPlayerID(id int) error {
	return s.repo.DeleteByPlayerID(id)
}

func (s *Service) DeleteByCode(code string) error {
	return s.repo.DeleteByCode(code)
}
