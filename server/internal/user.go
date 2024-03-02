package internal

import (
	"context"

	"github.com/SubhamMurarka/chat_app/models"
)

// type authentication struct {
// 	Token string `json:"token"`
// }

type Service interface {
	CreateUser(c context.Context, req *models.CreateUserReq) (*models.CreateUserRes, error)
	Login(c context.Context, req *models.LoginUserReq) (*models.LoginUserRes, error)
}
