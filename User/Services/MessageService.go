package Services

import (
	"context"

	models "github.com/SubhamMurarka/chat_app/User/Models"
	"github.com/SubhamMurarka/chat_app/User/Repository"
)

type MsgService interface {
	GetNewMessages(ctx context.Context, channelID int64) ([]models.MessageOP, error)
	PaginationMessages(ctx context.Context, channelID int64, lastid uint64) ([]models.MessageOP, error)
}

type msgService struct {
	msgrepo Repository.RepoInterface
}

func NewMsgService(repo Repository.RepoInterface) MsgService {
	return &msgService{
		msgrepo: repo,
	}
}

func (m *msgService) GetNewMessages(ctx context.Context, channelID int64) ([]models.MessageOP, error) {
	msg, err := m.msgrepo.GetNewMessages(ctx, channelID)
	return msg, err
}

func (m *msgService) PaginationMessages(ctx context.Context, channelID int64, lastid uint64) ([]models.MessageOP, error) {
	msg, err := m.msgrepo.PaginationMessages(ctx, channelID, lastid)
	return msg, err
}
