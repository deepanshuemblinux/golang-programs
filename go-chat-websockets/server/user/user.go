package user

import "context"

type User struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Repository interface {
	CreateUser(context.Context, *User) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
}

type Service interface {
	CreateUser(context.Context, *CreateUserReq) (*CreateUserResp, error)
	Login(context.Context, *LoginUserReq) (*LoginUserResp, error)
}

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResp struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResp struct {
	accessToken string
	ID          int64  `json:"id"`
	Username    string `json:"username"`
}
