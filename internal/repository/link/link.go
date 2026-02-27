package link

import (
	"database/sql"
	"errors"
	"plugin/internal/database/queries"
	"time"
)

type Link struct {
	PlayerID int
	Code     string
}

type LinkRepository interface {
	Create(playerID int, code string, expiresAt time.Time) error
	GetPlayerIDByCode(code string) (int, error)
	GetCodeByPlayerID(playerID int) (string, error)
	DeleteByPlayerID(playerID int) error
	DeleteByCode(code string) error
}

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) LinkRepository {
	return &repository{db: db}
}

func (r *repository) Create(playerID int, code string, expiresAt time.Time) error {
	if err := r.DeleteByPlayerID(playerID); err != nil {
		return err
	}

	_, err := r.db.Exec(queries.CreateLink, playerID, code, expiresAt)
	return err
}

func (r *repository) GetPlayerIDByCode(code string) (int, error) {
	var playerID int
	err := r.db.QueryRow(
		queries.GetIDByCode,
		code,
	).Scan(&playerID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return playerID, nil
}

func (r *repository) GetCodeByPlayerID(playerID int) (string, error) {
	var code string
	err := r.db.QueryRow(
		queries.GetCodeByID,
		playerID,
	).Scan(&code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return code, nil
}

func (r *repository) DeleteByPlayerID(playerID int) error {
	_, err := r.db.Exec(
		queries.DeleteLinkByPlayer,
		playerID,
	)
	return err
}

func (r *repository) DeleteByCode(code string) error {
	_, err := r.db.Exec(
		queries.DeleteLinkByCode,
		code,
	)
	return err
}
