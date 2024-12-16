//go:build integration
// +build integration

package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/repo"
	"avito-test-task/internal/usecase"
	"avito-test-task/pkg"
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/zap"
)

type UserIntegrationTest struct {
	suite.Suite
	userUsecase domain.UserUsecase
	userRepo    domain.UserRepo
	db          repo.IPool
	mockLg      *zap.Logger
	skipped     bool
}

func (u *UserIntegrationTest) BeforeAll(t provider.T) {
	host := os.Getenv("POSTGRES_TEST_HOST")
	port := os.Getenv("POSTGRES_TEST_PORT")
	connString := "postgres://test-user:test-password@" + host + ":" + port + "/test-db?sslmode=disable"

	var err error
	u.db, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf("error while connecting to db: %v", err.Error())
	}
	u.userRepo = repo.NewPostrgesUserRepo(u.db, nil)
	u.userUsecase = usecase.NewUserUsecase(u.userRepo)
	u.mockLg = pkg.CreateMockLogger()

	args := os.Args
	for _, arg := range args {
		if arg == "skipped" {
			u.skipped = true
		}
	}
}

func (u *UserIntegrationTest) AfterAll(t provider.T) {
	u.db.Close()
}

func (u *UserIntegrationTest) TestNormalRegister(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	userReq := domain.RegisterUserRequest{
		Email:    "user@mail.ru",
		Password: "password",
		UserType: domain.Client,
	}

	created, err := u.userUsecase.Register(context.Background(), &userReq, u.mockLg)

	t.Require().Nil(err)
	t.Require().NotEmpty(created)
}

func (u *UserIntegrationTest) TestBadNilRequestRegister(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	created, err := u.userUsecase.Register(context.Background(), nil, u.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.RegisterUserResponse{}, created)
}

func (u *UserIntegrationTest) TestBadUserTypeRegister(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	userReq := domain.RegisterUserRequest{
		Email:    "user@mail.ru",
		Password: "password",
		UserType: "type",
	}

	created, err := u.userUsecase.Register(context.Background(), &userReq, u.mockLg)

	t.Require().ErrorIs(err, domain.ErrUser_BadType)
	t.Require().Equal(domain.RegisterUserResponse{}, created)
}

func (u *UserIntegrationTest) TestBadMailRegister(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	userReq := domain.RegisterUserRequest{
		Email:    "kjfjkshfh",
		Password: "password",
		UserType: domain.Client,
	}

	created, err := u.userUsecase.Register(context.Background(), &userReq, u.mockLg)

	t.Require().ErrorIs(err, domain.ErrUser_BadMail)
	t.Require().Equal(domain.RegisterUserResponse{}, created)
}

func (u *UserIntegrationTest) TestBadPasswordRegister(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	userReq := domain.RegisterUserRequest{
		Email:    "user@mail.ru",
		Password: "",
		UserType: domain.Client,
	}

	created, err := u.userUsecase.Register(context.Background(), &userReq, u.mockLg)

	t.Require().ErrorIs(err, domain.ErrUser_BadPassword)
	t.Require().Equal(domain.RegisterUserResponse{}, created)
}

func (u *UserIntegrationTest) TestNormalClientLogin(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	userReq := domain.RegisterUserRequest{
		Email:    "userlogin@mail.ru",
		Password: "password",
		UserType: domain.Client,
	}

	created, _ := u.userUsecase.Register(context.Background(), &userReq, u.mockLg)

	loginReq := domain.LoginUserRequest{
		ID:       created.UserID,
		Password: userReq.Password,
	}

	loginResp, err := u.userUsecase.Login(context.Background(), &loginReq, u.mockLg)

	userID, _ := pkg.ExtractPayloadFromToken(loginResp.Token, "userID")
	userRole, _ := pkg.ExtractPayloadFromToken(loginResp.Token, "role")
	t.Require().Nil(err)
	t.Require().Equal(created.UserID, uuid.MustParse(userID))
	t.Require().Equal(userReq.UserType, userRole)
}

func (u *UserIntegrationTest) TestNormalModeratorLogin(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	userReq := domain.RegisterUserRequest{
		Email:    "modlogin@mail.ru",
		Password: "password",
		UserType: domain.Moderator,
	}

	created, _ := u.userUsecase.Register(context.Background(), &userReq, u.mockLg)

	loginReq := domain.LoginUserRequest{
		ID:       created.UserID,
		Password: userReq.Password,
	}

	loginResp, err := u.userUsecase.Login(context.Background(), &loginReq, u.mockLg)

	userID, _ := pkg.ExtractPayloadFromToken(loginResp.Token, "userID")
	userRole, _ := pkg.ExtractPayloadFromToken(loginResp.Token, "role")
	t.Require().Nil(err)
	t.Require().Equal(created.UserID, uuid.MustParse(userID))
	t.Require().Equal(userReq.UserType, userRole)
}

func (u *UserIntegrationTest) TestBadNilRequestLogin(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	loginResp, err := u.userUsecase.Login(context.Background(), nil, u.mockLg)

	t.Require().ErrorIs(err, domain.ErrUser_BadRequest)
	t.Require().Equal(domain.LoginUserResponse{}, loginResp)
}

func (u *UserIntegrationTest) TestBadAuthLogin(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	loginReq := domain.LoginUserRequest{
		ID:       uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db28"),
		Password: "jdfhjdsg",
	}

	loginResp, err := u.userUsecase.Login(context.Background(), &loginReq, u.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.LoginUserResponse{}, loginResp)
}

func (u *UserIntegrationTest) TestBadIDLogin(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	loginReq := domain.LoginUserRequest{
		ID:       uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db90"),
		Password: "jdfhjdsg",
	}

	loginResp, err := u.userUsecase.Login(context.Background(), &loginReq, u.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.LoginUserResponse{}, loginResp)
}

func (u *UserIntegrationTest) TestNormalDummyLogin(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	loginResp, err := u.userUsecase.DummyLogin(context.Background(), domain.Client, u.mockLg)

	t.Require().Nil(err)
	t.Require().NotEmpty(loginResp)
}

func (u *UserIntegrationTest) TestBadUserTypeDummyLogin(t provider.T) {
	if u.skipped {
		t.Skip()
	}

	loginResp, err := u.userUsecase.DummyLogin(context.Background(),
		"type", u.mockLg)

	t.Require().Error(err)
	t.Require().Empty(loginResp)
}

func TestUserIntegrationSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(UserIntegrationTest))
}
