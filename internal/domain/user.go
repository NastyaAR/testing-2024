package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	Moderator = "moderator"
	Client    = "client"
)

var (
	DummyMail     = "dummy@mail.ru"
	DummyPassword = "dummy_password"
)

var (
	ErrUser_BadType     = errors.New("bad user type")
	ErrUser_BadRequest  = errors.New("bad nil request")
	ErrUser_BadMail     = errors.New("bad mail")
	ErrUser_BadPassword = errors.New("bad password")
	ErrUser_BadId       = errors.New("bad user id ")
)

type User struct {
	UserID   uuid.UUID
	Mail     string
	Password string
	Role     string
}

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

type RegisterUserResponse struct {
	UserID uuid.UUID `json:"user_id"`
}

type LoginUserRequest struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

type FinalLoginUserRequest struct {
	ID   uuid.UUID `json:"id"`
	Code int       `json:"code"`
}

type FinalLoginUserResponse struct {
	Token string `json:"token"`
}

type LoginUserResponse struct {
	Message string `json:"message"`
}

type DummyLoginRequest struct {
	UserType string `json:"user_type"`
}

type UserUsecase interface {
	Register(ctx context.Context, userReq *RegisterUserRequest, lg *zap.Logger) (RegisterUserResponse, error)
	Login(ctx context.Context, userReq *LoginUserRequest, lg *zap.Logger) (LoginUserResponse, error)
	DummyLogin(ctx context.Context, userType string, lg *zap.Logger) (LoginUserResponse, error)
	FinalLogin(ctx context.Context, userReq *FinalLoginUserRequest, lg *zap.Logger) (FinalLoginUserResponse, error)
}

type UserRepo interface {
	Create(ctx context.Context, user *User, lg *zap.Logger) error
	DeleteByID(ctx context.Context, id uuid.UUID, lg *zap.Logger) error
	Update(ctx context.Context, newUserData *User, lg *zap.Logger) error
	GetByID(ctx context.Context, id uuid.UUID, lg *zap.Logger) (User, error)
	GetAll(ctx context.Context, offset int, limit int, lg *zap.Logger) ([]User, error)
	CreateCode(ctx context.Context, user *User, codeHash string, lg *zap.Logger) error
	GetHashCode(ctx context.Context, user *User, lg *zap.Logger) (string, error)
}

type CodeSender interface {
	SendCode(ctx context.Context, user *User, code int, lg *zap.Logger) error
}
