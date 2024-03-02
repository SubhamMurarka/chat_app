package internal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SubhamMurarka/chat_app/db"
	"github.com/SubhamMurarka/chat_app/models"
	"github.com/SubhamMurarka/chat_app/reddis"
)

type repository struct {
	db db.DBTX
}

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByName(ctx context.Context, username string) (bool, error)
	FindUserByEmail(ctx context.Context, email string) (bool, error)
	AddUserToRoomRedis(ctx context.Context, room string, cl *models.Client) error
	RemoveUserFromRoomRedis(ctx context.Context, room string, cl *models.Client) error
	IsUserRoomRedis(ctx context.Context, room string, cl *models.Client) bool
	GetAllMembersRedis(ctx context.Context, room string) []string
	// CreateRooms(ctx context.Context, room *CreateRoomdb) (*CreateRoomRes, error)
	// FindRepobyName(ctx context.Context, roomname string) (bool, error)
	// AddUserToRoom(ctx context.Context, userID int64, roomID int64) error
	// RemoveUserFromRoom(ctx context.Context, userID int64, roomID int64) error
	// CheckClientInRoom(ctx context.Context, userID int64, roomID int64) (bool, error)
}

func NewRepository(db db.DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&lastInsertId)
	if err != nil {
		return &models.User{}, err
	}
	user.ID = int64(lastInsertId)
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
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

func (r *repository) AddUserToRoomRedis(ctx context.Context, room string, cl *models.Client) error {
	_, err := reddis.RedisClient.SAdd(ctx, room, cl.ID).Result()
	return err
}

func (r *repository) RemoveUserFromRoomRedis(ctx context.Context, room string, cl *models.Client) error {
	_, err := reddis.RedisClient.SRem(ctx, room, cl.ID).Result()
	return err
}

func (r *repository) IsUserRoomRedis(ctx context.Context, room string, cl *models.Client) bool {
	is, err := reddis.RedisClient.SIsMember(ctx, room, cl.ID).Result()
	if err != nil {
		return false
	}
	return is
}

func (r *repository) GetAllMembersRedis(ctx context.Context, room string) []string {
	members, err := reddis.RedisClient.SMembers(ctx, room).Result()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return members
}

// func (r *repository) CreateRooms(ctx context.Context, room *CreateRoomdb) (*CreateRoomRes, error) {
// 	var insertedID int
// 	query := "INSERT INTO rooms(room_name, created_by) VALUES ($1, $2) returning id"
// 	err := r.db.QueryRowContext(ctx, query, room.RoomName, room.CreatedBY).Scan(&insertedID)

// 	if err != nil {
// 		return &CreateRoomRes{}, err
// 	}
// 	res := &CreateRoomRes{
// 		Name:   room.RoomName,
// 		RoomID: strconv.Itoa(insertedID),
// 	}

// 	return res, nil
// }

// func (r *repository) FindRepobyName(ctx context.Context, roomname string) (bool, error) {
// 	var roomID string
// 	query := "SELECT room_id FROM rooms WHERE room_name = $1"
// 	err := r.db.QueryRowContext(ctx, query, roomname).Scan(&roomID)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return false, nil
// 		}
// 		return false, err
// 	}

// 	return true, nil
// }

// func (r *repository) AddUserToRoom(ctx context.Context, userID int64, roomID int64) error {
// 	query := "INSERT INTO room_members(user_id, room_id) VALUES ($1, $2)"

// 	_, err := r.db.ExecContext(ctx, query, userID, roomID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *repository) RemoveUserFromRoom(ctx context.Context, userID int64, roomID int64) error {
// 	query := "DELETE FROM room_members WHERE user_id = $1 AND room_id = $2"

// 	_, err := r.db.ExecContext(ctx, query, userID, roomID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *repository) CheckClientInRoom(ctx context.Context, userID int64, roomID int64) (bool, error) {
// 	var id int
// 	query := "SELECT user_id FROM room_members WHERE user_id = userID AND room_id = roomID"
// 	err := r.db.QueryRowContext(ctx, query, userID, roomID).Scan(&id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return false, nil
// 		}
// 		return false, err
// 	}

// 	return true, nil
// }
