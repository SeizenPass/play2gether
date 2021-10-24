package mysql

import (
	"database/sql"
	"errors"
	"github.com/SeizenPass/play2gether/pkg/models"
)

// GameModel Game struct
type GameModel struct {
	DB *sql.DB
}

func (m *GameModel) Insert(title, imageLink, description string) (int, error) {
	stmt := `INSERT INTO games (title, image_link, description)
	VALUES(?, ?, ?)`

	transaction, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := transaction.Exec(stmt, title, imageLink, description)
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
func (m *GameModel) Get(id int) (*models.Game, error) {
	stmt := `SELECT id, title, image_link, description FROM games
	WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	// init a pointer to new Snippet struct
	s := &models.Game{}

	err := row.Scan(&s.ID, &s.Title, &s.ImageLink, &s.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil

}

func (m *GameModel) GetAll() ([]*models.Game, error) {
	stmt := `SELECT g.id, g.title, g.image_link, g.description, COUNT(go.user_id)
			FROM games g
			LEFT JOIN games_ownership go on g.id = go.game_id
			GROUP BY g.id
			ORDER BY 5 DESC, 2`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	games := []*models.Game{}

	for rows.Next() {
		s := &models.Game{}
		err = rows.Scan(&s.ID, &s.Title, &s.ImageLink, &s.Description, &s.Players)
		if err != nil {
			return nil, err
		}
		games = append(games, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

