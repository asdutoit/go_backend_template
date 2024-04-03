package models

import (
	"errors"

	"github.com/asdutoit/go_backend_template/db"
	"github.com/asdutoit/go_backend_template/utils"
)

type User struct {
	ID       int64
	Username string
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := `
	INSERT INTO users(username, email, password) 
	VALUES ($1, $2, $3) RETURNING id`

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	err = db.DB.QueryRow(query, u.Username, u.Email, hashedPassword).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (*User, error) {
	query := `
	SELECT id, username, email, password
	FROM users
	WHERE email = $1`

	row := db.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, err
}

func GetAllUsers() ([]User, error) {
	query := `
	SELECT id, username, email
	FROM users`

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *User) DeleteUser() error {
	query := `
	DELETE FROM users
	WHERE email = $1`

	_, err := db.DB.Exec(query, u.Email)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) ValidateCredentials() (int64, error) {
	user, err := GetUserByEmail(u.Email)

	if err != nil {
		return 0, err
	}

	valid := utils.CheckPasswordHash(u.Password, user.Password)

	if !valid {
		return 0, errors.New("invalid credentials")
	}

	return user.ID, nil
}
