package mysql

import (
	"database/sql"
	"errors"
	"github.com/SeizenPass/play2gether/pkg/models"
)

// ChatMessageModel struct
type ChatMessageModel struct {
	DB *sql.DB
}

func (m *ChatMessageModel) Insert(content string, senderID, receiverID int) (int, error) {
	stmt := `INSERT INTO chat_messages (sender_id, receiver_id, content)
	VALUES(?, ?, ?)`

	transaction, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := transaction.Exec(stmt, senderID, receiverID, content)
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
func (m *ChatMessageModel) Get(id int) (*models.ChatMessage, error) {
	stmt := `SELECT id, sender_id, receiver_id, content, is_read, created_at FROM chat_messages
	WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	// init a pointer to new Snippet struct
	s := &models.ChatMessage{}

	err := row.Scan(&s.ID, &s.SenderID, &s.ReceiverID, &s.Content, &s.IsRead, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil

}

func (m *ChatMessageModel) GetAllByUserID(id int) ([]*models.ChatMessage, error) {
	stmt := `SELECT id, sender_id, receiver_id, content, is_read, created_at 
			FROM chat_messages
			WHERE sender_id = ? OR receiver_id = ?
			ORDER BY created_at DESC, id DESC`
	//TODO: to not use order by id

	rows, err := m.DB.Query(stmt, id, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cs := []*models.ChatMessage{}

	for rows.Next() {
		s := &models.ChatMessage{}
		err = rows.Scan(&s.ID, &s.SenderID, &s.ReceiverID, &s.Content, &s.IsRead, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		cs = append(cs, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cs, nil
}

func (m *ChatMessageModel) GetDialogue(firstID int, secondID int) ([]*models.ChatMessage, error) {
	stmt := `SELECT id, sender_id, receiver_id, content, is_read, created_at 
			FROM chat_messages
			WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
			ORDER BY created_at`

	rows, err := m.DB.Query(stmt, firstID, secondID, secondID, firstID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cs := []*models.ChatMessage{}

	for rows.Next() {
		s := &models.ChatMessage{}
		err = rows.Scan(&s.ID, &s.SenderID, &s.ReceiverID, &s.CreatedAt, &s.IsRead, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		cs = append(cs, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cs, nil
}

