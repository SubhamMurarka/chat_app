package repositories

import (
	"context"
	"fmt"

	"github.com/SubhamMurarka/chat_app/models"
	"github.com/redis/go-redis/v9"
)

type RoomRepository interface {
	AddUserToRoomRedis(ctx context.Context, room string, cl *models.Client) error
	RemoveUserFromRoomRedis(ctx context.Context, room string, cl *models.Client) error
	IsUserRoomRedis(ctx context.Context, room string, cl *models.Client) bool
	GetAllMembersRedis(ctx context.Context, room string) []string
}

type roomRepository struct {
	redisClient *redis.Client
}

func NewRoomRepository(redisClient *redis.Client) RoomRepository {
	return &roomRepository{redisClient: redisClient}
}

func (r *roomRepository) AddUserToRoomRedis(ctx context.Context, room string, cl *models.Client) error {
	_, err := r.redisClient.SAdd(ctx, room, cl.ID).Result()
	return err
}

func (r *roomRepository) RemoveUserFromRoomRedis(ctx context.Context, room string, cl *models.Client) error {
	_, err := r.redisClient.SRem(ctx, room, cl.ID).Result()
	return err
}

func (r *roomRepository) IsUserRoomRedis(ctx context.Context, room string, cl *models.Client) bool {
	is, err := r.redisClient.SIsMember(ctx, room, cl.ID).Result()
	if err != nil {
		return false
	}
	return is
}

func (r *roomRepository) GetAllMembersRedis(ctx context.Context, room string) []string {
	members, err := r.redisClient.SMembers(ctx, room).Result()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return members
}
