package repositories

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

type User struct {
	ID       int
	Username string
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (repository UserRepository) GetUserByUsername(username string) (User, error) {
	sqlQuery := `SELECT * FROM users WHERE username = $1`

	var user User
	row := repository.db.QueryRow(sqlQuery, username)
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (repository UserRepository) IsUsernameTaken(username string) (bool, error) {
	sqlQuery := `SELECT 1 FROM users WHERE username = $1`
	res, err := repository.db.Exec(sqlQuery, username)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repository UserRepository) AddUser(username string) (int, error) {
	sqlQuery := `INSERT INTO users (username) VALUES ($1) RETURNING id`
	id := -1
	err := repository.db.QueryRow(sqlQuery, username).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (repository UserRepository) RemoveUser(username string) error {
	sqlQuery := `DELETE FROM users WHERE id = $1`
	_, err := repository.db.Exec(sqlQuery, username)
	if err != nil {
		return err
	}

	return nil
}
