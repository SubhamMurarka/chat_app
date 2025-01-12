package Repository

import (
	"context"
	"log"

	models "github.com/SubhamMurarka/chat_app/User/Models"
	util "github.com/SubhamMurarka/chat_app/User/Util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	var lastInsertId int

	query := `INSERT INTO users(username, email, password) VALUES ($1, $2, $3)
	          ON CONFLICT(email) DO NOTHING RETURNING id`

	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password).Scan(&lastInsertId)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Error : User already exists, %v", err)
			return &models.User{}, util.ErrEmailExists
		}
		log.Printf("Error : Internal error, %v", err)
		return &models.User{}, util.ErrInternal
	}

	user.ID = int64(lastInsertId)

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}

	query := "SELECT id, email, username, password FROM users WHERE email = $1"

	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Invalid Credentials : %v", err)
			return &models.User{}, util.ErrCredentials
		}
		log.Printf("Error : Internal error, %v", err)
		return &models.User{}, util.ErrInternal
	}

	return u, nil
}
