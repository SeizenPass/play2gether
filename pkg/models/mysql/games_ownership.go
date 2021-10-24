package mysql

import (
	"database/sql"
	"errors"
	"github.com/SeizenPass/play2gether/pkg/models"
)

type GameOwnershipModel struct {
	DB *sql.DB
}

// Insert ownerships into DB
func (m *GameOwnershipModel) Insert(gameID, userID int) (int, error) {
	stmt := `INSERT INTO games_ownership (game_id, user_id)
	VALUES(?, ?)`

	transaction, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := transaction.Exec(stmt, gameID, userID)
	if err != nil {
		transaction.Rollback()
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	err = transaction.Commit()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

// Get func
func (m *GameOwnershipModel) Get(id int) (*models.GameOwnership, error) {
	stmt := `SELECT id, game_id, user_id FROM games_ownership
	WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	// init a pointer to new Snippet struct
	s := &models.GameOwnership{}

	err := row.Scan(&s.ID, &s.GameID, &s.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// Get func
func (m *GameOwnershipModel) GetByGameID(id int) ([]*models.GameOwnership, error) {
	stmt := `SELECT id, game_id, user_id FROM games_ownership
	WHERE game_id = ?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ows := []*models.GameOwnership{}

	for rows.Next() {
		s := &models.GameOwnership{}
		err = rows.Scan(&s.ID, &s.GameID, &s.UserID)
		if err != nil {
			return nil, err
		}
		ows = append(ows, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ows, nil
}

func (m *GameOwnershipModel) GetByUserID(id int) ([]*models.GameOwnership, error) {
	stmt := `SELECT id, game_id, user_id FROM games_ownership
	WHERE user_id = ?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ows := []*models.GameOwnership{}

	for rows.Next() {
		s := &models.GameOwnership{}
		err = rows.Scan(&s.ID, &s.GameID, &s.UserID)
		if err != nil {
			return nil, err
		}
		ows = append(ows, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ows, nil
}

func (m *GameOwnershipModel) GetByUserIDAndGameID(userID, gameID int) (*models.GameOwnership, error) {
	stmt := `SELECT id, game_id, user_id FROM games_ownership
	WHERE user_id = ? AND game_id = ?`

	row := m.DB.QueryRow(stmt, userID, gameID)

	// init a pointer to new Snippet struct
	s := &models.GameOwnership{}

	err := row.Scan(&s.ID, &s.GameID, &s.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *GameOwnershipModel) Remove(id int) (error) {
	stmt := `DELETE FROM games_ownership
	WHERE id = ?`
	_, err := m.DB.Exec(stmt, id)
	return err
}

