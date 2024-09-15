package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/usecase"
	"avito-test-task/pkg"
	mock_domain "avito-test-task/tests/mocks"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
)

type UserUsecaseTest struct {
	suite.Suite
	userRepoMock *mock_domain.MockUserRepo
	lg           *zap.Logger
}

func (u *UserUsecaseTest) BeforeAll(t provider.T) {
	t.Log("Init mock")
	ctrl := gomock.NewController(t)
	u.userRepoMock = mock_domain.NewMockUserRepo(ctrl)
	u.lg, _ = pkg.CreateLogger("../log.log", "prod")
}

func (u *UserUsecaseTest) TestNormalRegister(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	req := domain.RegisterUserRequest{
		Email:    "test@mail.ru",
		Password: "password",
		UserType: domain.Client,
	}

	u.userRepoMock.EXPECT().Create(context.Background(), gomock.Any(), u.lg).Return(nil)

	_, err := userUsecase.Register(context.Background(), &req, u.lg)

	t.Require().Nil(err)
}

func (u *UserUsecaseTest) TestBadPasswordRegister(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	req := domain.RegisterUserRequest{
		Email:    "test@mail.ru",
		Password: "",
		UserType: domain.Client,
	}

	_, err := userUsecase.Register(context.Background(), &req, u.lg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestBadMailRegister(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	req := domain.RegisterUserRequest{
		Email:    "testmail.ru",
		Password: "password",
		UserType: domain.Client,
	}

	_, err := userUsecase.Register(context.Background(), &req, u.lg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestBadUserTypeRegister(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	req := domain.RegisterUserRequest{
		Email:    "test@mail.ru",
		Password: "password",
		UserType: "user",
	}

	_, err := userUsecase.Register(context.Background(), &req, u.lg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestBadRepoCallRegister(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	req := domain.RegisterUserRequest{
		Email:    "test@mail.ru",
		Password: "password",
		UserType: domain.Client,
	}

	u.userRepoMock.EXPECT().Create(context.Background(), gomock.Any(), u.lg).Return(errors.New("error"))

	_, err := userUsecase.Register(context.Background(), &req, u.lg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestNormalLogin(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	uid := uuid.New()
	usr := domain.User{
		UserID:   uid,
		Mail:     "test@mail.ru",
		Password: "$2a$10$uKeo8Mj.unKogJ6rV138heK1J/x./xqA97cVTepiq7evt9sER8EPG",
		Role:     domain.Client,
	}

	req := domain.LoginUserRequest{
		ID:       uid,
		Password: "mysecretpassword",
	}

	u.userRepoMock.EXPECT().GetByID(context.Background(), req.ID, u.lg).Return(usr, nil)

	_, err := userUsecase.Login(context.Background(), &req, u.lg)

	t.Require().Nil(err)
}

func (u *UserUsecaseTest) TestBadPasswordLogin(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	uid := uuid.New()
	usr := domain.User{
		UserID:   uid,
		Mail:     "test@mail.ru",
		Password: "$2a$10$uKeo8Mj.unKogJ6rV138heK1J/x./xqA97cVTepiq7evt9sER8EPG",
		Role:     domain.Client,
	}

	req := domain.LoginUserRequest{
		ID:       uid,
		Password: "password",
	}

	u.userRepoMock.EXPECT().GetByID(context.Background(), req.ID, u.lg).Return(usr, nil)

	_, err := userUsecase.Login(context.Background(), &req, u.lg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestBadRepoCallLogin(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	uid := uuid.New()
	req := domain.LoginUserRequest{
		ID:       uid,
		Password: "password",
	}

	u.userRepoMock.EXPECT().GetByID(context.Background(), req.ID, u.lg).Return(domain.User{}, errors.New("error"))

	_, err := userUsecase.Login(context.Background(), &req, u.lg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestNormalDummyLogin(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	u.userRepoMock.EXPECT().Create(context.Background(), gomock.Any(), u.lg)
	_, err := userUsecase.DummyLogin(context.Background(), domain.Moderator, u.lg)
	t.Require().Nil(err)
}

func (u *UserUsecaseTest) TestBadUserTypeDummyLogin(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	_, err := userUsecase.DummyLogin(context.Background(), "user", u.lg)
	t.Require().Error(err)
}

func TestUserUsecaseSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(UserUsecaseTest))
}
