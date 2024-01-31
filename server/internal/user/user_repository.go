package user

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int

	query := "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) returning id"

	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&lastInsertId)

	if err != nil {
		return &User{}, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}
	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return &User{}, nil
	}

	return &u, nil
}

func (r *repository) FindUserByName(ctx context.Context, username string) (bool, error) {
	var userID int
	query := "SELECT id FROM users WHERE username = $1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (bool, error) {
	var userID string
	query := "SELECT id FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// func (r *repository) UpdateAllTokens(ctx context.Context, token string, refreshToken string, email string) (*User, error) {
// 	u, err := r.GetUserByEmail(ctx, email)
// 	if err != nil {
// 		return &User{}, err
// 	}

// 	query := "UPDATE users SET token = $1, refreshtoken = $2 WHERE email = $3"
// 	_, err = r.db.ExecContext(ctx, query, token, refreshToken, u.Email)
// 	if err != nil {
// 		return &User{}, err
// 	}

// 	u.Token = token
// 	u.RefreshToken = refreshToken

// 	return u, nil
// }
