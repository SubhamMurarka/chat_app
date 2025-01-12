package Services

import (
	"context"
	"log"
	"strconv"

	models "github.com/SubhamMurarka/chat_app/User/Models"
	Repositories "github.com/SubhamMurarka/chat_app/User/Repository"
	util "github.com/SubhamMurarka/chat_app/User/Util"
)

type UserService interface {
	CreateUser(c context.Context, req *models.CreateUserReq) (*models.CreateUserRes, error)
	Login(c context.Context, req *models.LoginUserReq) (*models.LoginUserRes, error)
}

type userService struct {
	repo Repositories.UserRepository
}

func NewUserService(userRepository Repositories.UserRepository) UserService {
	return &userService{
		repo: userRepository,
	}
}

func (s *userService) CreateUser(c context.Context, req *models.CreateUserReq) (*models.CreateUserRes, error) {

	if req.Email == "" || req.Username == "" {
		log.Println("No username and password")
		return nil, util.ErrUserEmailData
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Println("Error hashing password")
		return nil, util.ErrInternal
	}

	u := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.repo.CreateUser(c, u)
	if err != nil {
		return nil, err
	}

	res := &models.CreateUserRes{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *userService) Login(c context.Context, req *models.LoginUserReq) (*models.LoginUserRes, error) {
	u, err := s.repo.GetUserByEmail(c, req.Email)
	if err != nil {
		return &models.LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &models.LoginUserRes{}, err
	}

	token, err := util.GenerateAllTokens(strconv.FormatInt(u.ID, 10), u.Username, u.Email)
	if err != nil {
		return &models.LoginUserRes{}, err
	}

	logRes := &models.LoginUserRes{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		Token:    token,
	}

	return logRes, nil
}
