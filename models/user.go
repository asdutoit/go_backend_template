package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/asdutoit/go_backend_template/db"
	"github.com/asdutoit/go_backend_template/utils"
)

type User struct {
	ID         int64
	Username   string
	Email      string `binding:"required"`
	Password   string
	First_name string
	Last_name  string
	Picture    string
}

func (u *User) Save() error {
	query := `
	INSERT INTO users(username, email, password, first_name, last_name, picture) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	err = db.DB.QueryRow(query, u.Username, u.Email, hashedPassword, u.First_name, u.Last_name, u.Picture).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (*User, error) {
	query := `
	SELECT id, username, email, first_name, last_name, picture, password
	FROM users
	WHERE email = $1`

	row := db.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.First_name, &user.Last_name, &user.Picture, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			// No user was found, return nil for the user
			return nil, nil
		} else {
			// A real error occurred
			return nil, err
		}
	}

	return &user, err
}

func GetAllUsers() ([]User, error) {
	query := `
	SELECT id, username, email, first_name, last_name, picture
	FROM users`

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.First_name, &user.Last_name, &user.Picture)
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

func (u User) String() string {
	return fmt.Sprintf("User{ID: %d, Email: %s, FirstName: %s, LastName: %s, Picture: %s}", u.ID, u.Email, u.First_name, u.Last_name, u.Picture)
}

func (u *User) ValidateCredentials() (int64, error) {
	user, err := GetUserByEmail(u.Email)

	// Print the user to console
	fmt.Println(user.String())

	if err != nil {
		return 0, err
	}

	valid := utils.CheckPasswordHash(u.Password, user.Password)

	if !valid {
		return 0, errors.New("invalid credentials")
	}

	return user.ID, nil
}

func (u *User) Update() error {
	query := `
	UPDATE users
	SET first_name = $1, last_name = $2, picture = $3
	WHERE email = $4`

	_, err := db.DB.Exec(query, u.First_name, u.Last_name, u.Picture, u.Email)

	if err != nil {
		return err
	}

	return nil
}
