package mysql

import (
	"database/sql"
	"github.com/SeizenPass/play2gether/pkg/models"
)

type ReviewModel struct {
	DB *sql.DB
}

func (m *ReviewModel) Insert(text string, reviewerID, reviewedID int) (int, error) {
	stmt := `INSERT INTO reviews (review_text, reviewer_id, reviewed_id)
	VALUES(?, ?, ?)`

	transaction, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := transaction.Exec(stmt, text, reviewerID, reviewedID)
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

func (m *ReviewModel) Get(id int) (*models.Review, error) {
	s := &models.Review{}

	stmt := `SELECT id, review_text, reviewer_id, reviewed_id FROM reviews WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.ReviewText, &s.ReviewerID, &s.ReviewedID)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *ReviewModel) GetByReviewedID(id int) ([]*models.Review, error) {
	stmt := `SELECT id, review_text, reviewer_id, reviewed_id FROM reviews
	WHERE reviewed_id = ?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	revs := []*models.Review{}

	for rows.Next() {
		s := &models.Review{}
		err = rows.Scan(&s.ID, &s.ReviewText, &s.ReviewerID, &s.ReviewedID)
		if err != nil {
			return nil, err
		}
		revs = append(revs, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return revs, nil
}