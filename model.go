package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mujz/ipp/authentication"
)

type User struct {
	ID       int     `jsonapi:"primary,User"`
	Email    string  `jsonapi:"attr,email"`
	Password string  `jsonapi:"attr,password"`
	Number   *Number `jsonapi:"relation,number"`
}

type Number struct {
	ID    int `jsonapi:"primary,number"`
	Value int `jsonapi:"attr,value"`
}

type Model struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
	*sql.DB
}

func (m *Model) Open() error {
	url := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		m.User, m.Password, m.Host, m.Port, m.DBName, m.SSLMode,
	)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	m.DB = db
	return m.DB.Ping()
}

func (m *Model) GetNumber(id int) (*Number, error) {
	number := &Number{ID: id}
	err := m.QueryRow("SELECT num FROM users WHERE id = $1", id).Scan(&number.Value)
	return number, err
}

func (m *Model) IncrementNumber(id int) (*Number, error) {
	number := &Number{ID: id}
	err := m.QueryRow("UPDATE users SET num = num + 1 WHERE id = $1 RETURNING num", id).Scan(&number.Value)
	return number, err
}

func (m *Model) UpdateNumber(userID, newValue int) (*Number, error) {
	number := &Number{ID: userID}
	err := m.QueryRow(
		"UPDATE users SET num = $1 WHERE id = $2 RETURNING num",
		newValue, userID,
	).Scan(&number.Value)
	return number, err
}

func (m Model) Get(email string) (authentication.User, error) {
	user := authentication.User{}
	err := m.QueryRow(
		"SELECT id, email, password FROM users WHERE email = $1", email,
	).Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}

func (m Model) Create(email, password string) (authentication.User, error) {
	user := authentication.User{}

	err := m.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email",
		email, password,
	).Scan(&user.ID, &user.Username)

	return user, err
}

func (m Model) CreateFacebookUser(fbID string) (int, error) {
	var id int
	// Insert user if it doesn't exist
	q := fmt.Sprintf(
		`INSERT INTO users (facebook_id)
		SELECT '%s' WHERE NOT EXISTS (
			SELECT id FROM users WHERE facebook_id = '%s'
		)`, fbID, fbID)

	if _, err := m.Exec(q); err != nil {
		return 0, err
	}

	err := m.QueryRow("SELECT id FROM users WHERE facebook_id = $1", fbID).Scan(&id)
	return id, err
}
