package balance

import (
	"errors"
	"fmt"
	"plugin/internal/config"
	"plugin/internal/service/player"
	"plugin/internal/service/wallet"
)

type Result struct {
	PlayerID int
	Name     string
	Balance  int
	Message  string
}

func Balance(
	playerID int,
	cfg config.Config,
	player *player.Service,
	wallet *wallet.Service,
) (*Result, error) {
	p, err := player.GetPlayerByID(playerID)
	if err != nil {
		return nil, errors.New("error occurred, please try again later")
	}

	balance, err := wallet.GetBalance(playerID)
	if err != nil {
		return nil, errors.New("error occurred, please try again later")
	}

	return &Result{
		PlayerID: playerID,
		Name:     p.Name,
		Balance:  balance,
		Message:  fmt.Sprintf("Your balance is %s%d", cfg.Gambling.Currency, balance),
	}, nil
}
