package internal

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/SubhamMurarka/chat_app/models"
	"github.com/SubhamMurarka/chat_app/util"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(8) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *models.CreateUserReq) (*models.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	if req.Email == "" || req.Username == "" {
		return nil, errors.New("email and username cannot be empty")
	}

	exists, err := s.Repository.FindUserByEmail(ctx, req.Email)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("email already exists")
	}

	exists, err = s.Repository.FindUserByName(ctx, req.Username)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)

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

func (s *service) Login(c context.Context, req *models.LoginUserReq) (*models.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &models.LoginUserRes{}, errors.New("invalid credentials")
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &models.LoginUserRes{}, errors.New("invalid credentials")
	}

	token, err := util.GenerateAllTokens(strconv.FormatInt(u.ID, 10), u.Username, u.Email)

	if err != nil {
		return &models.LoginUserRes{}, errors.New("try to login again")
	}

	logRes := &models.LoginUserRes{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		Token:    token,
	}

	return logRes, nil
}
