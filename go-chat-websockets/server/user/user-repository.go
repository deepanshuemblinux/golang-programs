package user

import (
	"context"
	"database/sql"
	"fmt"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	query := "SELECT id FROM users where username=$1"
	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Username).Scan(&id)
	if err == nil {
		return nil, fmt.Errorf("username %s not available", user.Username)
	}
	query = "INSERT INTO users (username, password, email) VALUES ($1, $2, $3) returning ID"

	err = r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&id)
	if err != nil {

		return nil, err
	}
	user.ID = id
	fmt.Println(user)
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := User{}
	query := "Select id, email, username, password from users where email=$1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email id %s does not exist", email)
		}
	}
	return &user, nil
}
