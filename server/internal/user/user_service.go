package user

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/SubhamMurarka/chat_app/util"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
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
		return nil, errors.New("Email already exists")
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

	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)

	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}
	return res, nil
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &LoginUserRes{}, errors.New("invalid credentials")
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &LoginUserRes{}, errors.New("invalid credentials")
	}

	token, err := util.GenerateAllTokens(u.Username, u.Email)

	if err != nil {
		return &LoginUserRes{}, errors.New("try to login again")
	}

	logRes := &LoginUserRes{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		Token:    token,
	}

	return logRes, nil
}
