package wallet

import (
	"fmt"
	"plugin/internal/repository/wallet"
)

type Service struct {
	repo wallet.WalletRepository
}

func New(repo wallet.WalletRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateWallet(playerID int, initialBalance int) error {
	exists, err := s.repo.Exists(playerID)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("wallet already exists for player %d", playerID)
	}

	return s.repo.Create(playerID, initialBalance)
}

func (s *Service) GetBalance(playerID int) (int, error) {
	return s.repo.GetBalance(playerID)
}

func (s *Service) SetBalance(playerID int, amount int) error {
	return s.repo.SetBalance(playerID, amount)
}

func (s *Service) Deposit(playerID int, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}
	return s.repo.Deposit(playerID, amount)
}

func (s *Service) Withdraw(playerID int, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("withdraw amount must be positive")
	}

	return s.repo.Withdraw(playerID, amount)
}

func (s *Service) DeleteWallet(playerID int) error {
	return s.repo.Delete(playerID)
}

func (s *Service) Exists(playerID int) (bool, error) {
	return s.repo.Exists(playerID)
}

func (s *Service) GetTop5RichestWallets() ([]wallet.PlayerWallet, error) {
	return s.repo.GetTopWallets(5)
}

func (s *Service) GetTop5PoorestWallets() ([]wallet.PlayerWallet, error) {
	return s.repo.GetBottomWallets(5)
}

func (s *Service) GetTop10RichestWallets() ([]wallet.PlayerWallet, error) {
	return s.repo.GetTopWallets(10)
}

func (s *Service) GetTop10PoorestWallets() ([]wallet.PlayerWallet, error) {
	return s.repo.GetBottomWallets(10)
}
