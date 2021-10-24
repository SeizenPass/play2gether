package mysql

import (
	"database/sql"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/SeizenPass/play2gether/pkg/models"
	"github.com/go-sql-driver/mysql"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password, imageLink string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created, image_link)
	VALUES(?, ?, ?, UTC_TIMESTAMP(), ?)`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword), imageLink)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

func (m *UserModel) Update(id int, bio string) error {

	stmt := `UPDATE users SET bio = ? WHERE id = ?`

	_, err := m.DB.Exec(stmt, bio, id)
	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}

	stmt := `SELECT id, name, email, created, image_link, bio FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Email, &s.Created, &s.ImageLink, &s.Bio)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *UserModel) GetAll() ([]*models.User, error) {
	stmt := `SELECT id, name, email, image_link, bio
			FROM users`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		s := &models.User{}
		err = rows.Scan(&s.ID, &s.Name, &s.Email, &s.ImageLink, &s.Bio)
		if err != nil {
			return nil, err
		}
		users = append(users, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}