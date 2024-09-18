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
	mockLg       *zap.Logger
}

func (u *UserUsecaseTest) BeforeAll(t provider.T) {
	t.Log("Init mock")
	ctrl := gomock.NewController(t)
	u.userRepoMock = mock_domain.NewMockUserRepo(ctrl)
	u.mockLg = pkg.CreateMockLogger()
}

func (u *UserUsecaseTest) TestNormalRegister(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	req := domain.RegisterUserRequest{
		Email:    "test@mail.ru",
		Password: "password",
		UserType: domain.Client,
	}

	u.userRepoMock.EXPECT().Create(context.Background(), gomock.Any(), u.mockLg).Return(nil)

	_, err := userUsecase.Register(context.Background(), &req, u.mockLg)

	t.Require().Nil(err)
}

func (u *UserUsecaseTest) TestBadPasswordRegister(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	req := domain.RegisterUserRequest{
		Email:    "test@mail.ru",
		Password: "",
		UserType: domain.Client,
	}

	_, err := userUsecase.Register(context.Background(), &req, u.mockLg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestBadMailRegister(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	req := domain.RegisterUserRequest{
		Email:    "testmail.ru",
		Password: "password",
		UserType: domain.Client,
	}

	_, err := userUsecase.Register(context.Background(), &req, u.mockLg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestBadUserTypeRegister(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	req := domain.RegisterUserRequest{
		Email:    "test@mail.ru",
		Password: "password",
		UserType: "user",
	}

	_, err := userUsecase.Register(context.Background(), &req, u.mockLg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestBadRepoCallRegister(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	req := domain.RegisterUserRequest{
		Email:    "test@mail.ru",
		Password: "password",
		UserType: domain.Client,
	}

	u.userRepoMock.EXPECT().Create(context.Background(), gomock.Any(), u.mockLg).Return(errors.New("error"))

	_, err := userUsecase.Register(context.Background(), &req, u.mockLg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestNormalLogin(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)
	clientBuilder := NormalClientUserBuilder{}

	clientBuilder.SetRole()
	clientBuilder.SetMail()
	clientBuilder.SetPassword()
	clientBuilder.SetUid("019126ee-2b7d-758e-bb22-fe2e45b2db40")
	usr := clientBuilder.GetUser()

	req := domain.LoginUserRequest{
		ID:       usr.UserID,
		Password: "password",
	}

	u.userRepoMock.EXPECT().GetByID(context.Background(), req.ID, u.mockLg).Return(usr, nil)

	_, err := userUsecase.Login(context.Background(), &req, u.mockLg)

	t.Require().Nil(err)
}

func (u *UserUsecaseTest) TestBadPasswordLogin(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)
	clientBuilder := NormalClientUserBuilder{}

	clientBuilder.SetRole()
	clientBuilder.SetMail()
	clientBuilder.SetPassword()
	clientBuilder.SetUid("019126ee-2b7d-758e-bb22-fe2e45b2db40")
	usr := clientBuilder.GetUser()

	req := domain.LoginUserRequest{
		ID:       usr.UserID,
		Password: "badpassword",
	}

	u.userRepoMock.EXPECT().GetByID(context.Background(), req.ID, u.mockLg).Return(usr, nil)

	_, err := userUsecase.Login(context.Background(), &req, u.mockLg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestBadRepoCallLogin(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	uid := uuid.New()
	req := domain.LoginUserRequest{
		ID:       uid,
		Password: "password",
	}

	u.userRepoMock.EXPECT().GetByID(context.Background(), req.ID, u.mockLg).Return(domain.User{}, errors.New("error"))

	_, err := userUsecase.Login(context.Background(), &req, u.mockLg)

	t.Require().Error(err)
}

func (u *UserUsecaseTest) TestNormalDummyLogin(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	u.userRepoMock.EXPECT().Create(context.Background(), gomock.Any(), u.mockLg)
	_, err := userUsecase.DummyLogin(context.Background(), domain.Moderator, u.mockLg)
	t.Require().Nil(err)
}

func (u *UserUsecaseTest) TestBadUserTypeDummyLogin(t provider.T) {
	userUsecase := usecase.NewUserUsecase(u.userRepoMock)

	_, err := userUsecase.DummyLogin(context.Background(), "user", u.mockLg)
	t.Require().Error(err)
}

func TestUserUsecaseSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(UserUsecaseTest))
}
