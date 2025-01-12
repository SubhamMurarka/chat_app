package Services

import (
	"context"

	models "github.com/SubhamMurarka/chat_app/User/Models"
	"github.com/SubhamMurarka/chat_app/User/Repository"
)

type RoomService interface {
	AddUserToRoom(ctx context.Context, userid int64, channelid int64) error
	GetAllMembers(ctx context.Context, channel_id int64) ([]models.Members, error)
	// IsRoom(ctx context.Context, room string) (int64, error)
	AddRoom(ctx context.Context, room string, typ string) (int64, error)
	IsUSerMember(ctx context.Context, userid int64, channelid int64) bool
	GetAllUserRoom(ctx context.Context, user_id int64) ([]models.Room, error)
}

type roomService struct {
	roomrepo Repository.RoomRepository
}

func NewRoomService(roomRepository Repository.RoomRepository) RoomService {
	return &roomService{
		roomrepo: roomRepository,
	}
}

func (r *roomService) AddRoom(ctx context.Context, room string, typ string) (int64, error) {
	id, err := r.roomrepo.AddRoom(ctx, room, typ)
	return id, err
}

func (r *roomService) AddUserToRoom(ctx context.Context, userid int64, channelid int64) error {
	err := r.roomrepo.AddUserToRoom(ctx, userid, channelid)
	return err
}

// func (r *roomService) IsRoom(ctx context.Context, room string) (int64, error) {
// 	channel_id, err := r.roomrepo.IsRoom(ctx, room)
// 	return channel_id, err
// }

func (r *roomService) GetAllMembers(ctx context.Context, channel_id int64) ([]models.Members, error) {
	users, err := r.roomrepo.GetAllMembers(ctx, channel_id)
	return users, err
}

func (r *roomService) GetAllUserRoom(ctx context.Context, user_id int64) ([]models.Room, error) {
	room, err := r.roomrepo.GetAllUserRoom(ctx, user_id)
	return room, err
}

func (r *roomService) IsUSerMember(ctx context.Context, userid int64, channelid int64) bool {
	is := r.roomrepo.IsUSerMember(ctx, userid, channelid)
	return is
}
