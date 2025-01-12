package Repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	models "github.com/SubhamMurarka/chat_app/User/Models"
	util "github.com/SubhamMurarka/chat_app/User/Util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository interface {
	AddUserToRoom(ctx context.Context, user_id int64, channel_id int64) error
	GetAllMembers(ctx context.Context, channel_id int64) ([]models.Members, error)
	// IsRoom(ctx context.Context, room string) (int64, error)
	AddRoom(ctx context.Context, room string, typ string) (int64, error)
	GetAllUserRoom(ctx context.Context, user_id int64) ([]models.Room, error)
	IsUSerMember(ctx context.Context, user_id int64, channel_id int64) bool
	UpdateStatus(ctx context.Context, channelid int64, userid int64) error
}

type roomRepository struct {
	db *pgxpool.Pool
}

func NewRoomRepository(db *pgxpool.Pool) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) AddRoom(ctx context.Context, room string, typ string) (int64, error) {
	var lastInsertId int64

	query := `INSERT INTO channels(name, type) VALUES ($1, $2)
          ON CONFLICT(name) DO NOTHING RETURNING id`

	err := r.db.QueryRow(ctx, query, room, typ).Scan(&lastInsertId)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("room already exists : %v", err)
			return -1, util.ErrRoomExists
		}
		log.Printf("Error Fetching channel : %v", err)
		return -1, util.ErrInternal
	}

	return lastInsertId, nil
}

func (r *roomRepository) AddUserToRoom(ctx context.Context, user_id int64, channel_id int64) error {
	var channelExists, isAlreadyMember bool

	query := `WITH ins AS (
        INSERT INTO memberships (user_id, channel_id)
        SELECT $1, id FROM channels WHERE id = $2
        ON CONFLICT (user_id, channel_id) DO NOTHING
        RETURNING id
    ), chk AS (
        SELECT EXISTS(SELECT 1 FROM channels WHERE id = $2) as channel_exists,
               EXISTS(SELECT 1 FROM memberships WHERE user_id = $1 AND channel_id = $2) as already_member
    )
    SELECT chk.channel_exists, chk.already_member FROM ins
    RIGHT JOIN chk ON true;`

	err := r.db.QueryRow(ctx, query, user_id, channel_id).Scan(&channelExists, &isAlreadyMember)
	if err != nil {
		log.Printf("Error adding user to room: %v", err)
		return util.ErrInternal
	}

	if !channelExists {
		return util.ErrRoomExist
	}
	if isAlreadyMember {
		return util.ErrAlreadyJoined
	}
	return nil
}

func (r *roomRepository) GetAllMembers(ctx context.Context, channel_id int64) ([]models.Members, error) {
	var users []models.Members

	query := `
    SELECT 
        m.user_id, 
        m.channel_id, 
        u.username, 
        m.online_until
    FROM memberships m
    JOIN users u ON m.user_id = u.id
    WHERE m.channel_id = $1
    AND EXISTS (SELECT 1 FROM channels WHERE id = $1)
`

	rows, err := r.db.Query(ctx, query, channel_id)
	if err != nil {
		log.Printf("Error getting all members : %v", err)
		return nil, util.ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		var user models.Members

		if err := rows.Scan(&user.UserID, &user.ChannelID, &user.Username, &user.Online_At); err != nil {
			fmt.Printf("Error scanning rows : %v", err)
			return nil, util.ErrInternal
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error scanning rows : %v", err)
		return nil, util.ErrInternal
	}

	return users, nil
}

// func (r *roomRepository) IsRoom(ctx context.Context,  int64) (int64, error) {
// 	var insertedID int64

// 	query := "SELECT id FROM channels WHERE name = $1"

// 	err := r.db.QueryRow(ctx, query, room).Scan(&insertedID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Printf("No channel present: %v", err)
// 			return -1, util.ErrRoomExist
// 		}

// 		log.Printf("Error Fetching channel : %v", err)

// 		return -1, util.ErrInternal
// 	}

// 	return insertedID, nil
// }

func (r *roomRepository) GetAllUserRoom(ctx context.Context, user_id int64) ([]models.Room, error) {
	var rooms []models.Room

	query := `SELECT c.id AS room_id, c.name AS roomname, c.type AS type
			  FROM memberships m
			  JOIN channels c ON m.channel_id = c.id
	          WHERE m.user_id = $1`

	rows, err := r.db.Query(ctx, query, user_id)
	if err != nil {
		log.Printf("Error getting all rooms : %v", err)
		return nil, util.ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		var room models.Room

		if err := rows.Scan(&room.RoomID, &room.RoomName, &room.Typ); err != nil {
			fmt.Printf("Error scanning rows : %v", err)
			return nil, util.ErrInternal
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error scanning rows : %v", err)
		return nil, util.ErrInternal
	}

	return rooms, nil
}

func (r *roomRepository) IsUSerMember(ctx context.Context, user_id int64, channel_id int64) bool {
	var insertedID int64

	query := "SELECT id FROM memberships WHERE user_id = $1 and channel_id = $2"

	err := r.db.QueryRow(ctx, query, user_id, channel_id).Scan(&insertedID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No Entry found found for : ", user_id, channel_id)
			return false
		}
		log.Printf("Error getting data : %v", err)
		return false
	}

	return true
}

func (r *roomRepository) UpdateStatus(ctx context.Context, channelid int64, userid int64) error {

	query := `UPDATE memberships
			  SET online_at = current_timestamp
			  WHERE user_id = $1 AND channel_id = $2`

	cmdTag, err := r.db.Exec(ctx, query, userid, channelid)
	if err != nil {
		log.Printf("Error updating for user %d and channel %d Error : %w", userid, channelid, err)
		return util.ErrInternal
	}

	log.Println("row affected : ", cmdTag.RowsAffected())

	return nil
}
