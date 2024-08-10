package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deepanshuemblinux/go-chat-websockets/util"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(r Repository) Service {
	return &service{
		r,
		time.Second * 2,
	}
}

func (s *service) CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUserResp, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	hashed_password, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed_password,
	}
	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return &CreateUserResp{
		ID:       r.ID,
		Username: r.Username,
		Email:    r.Email,
	}, nil
}

func (s *service) Login(ctx context.Context, req *LoginUserReq) (*LoginUserResp, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	user, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid user name or password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		ID:       strconv.Itoa(int(user.ID)),
		UserName: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "deepanshu",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	})
	ss, err := token.SignedString([]byte("deepanshu"))
	if err != nil {
		return nil, err
	}
	resp := &LoginUserResp{
		accessToken: ss,
		ID:          user.ID,
		Username:    user.Username,
	}
	return resp, nil
}
