//go:build unit
// +build unit

package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/repo"
	"avito-test-task/pkg"
	mock_domain "avito-test-task/tests/mocks"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
	"time"
)

type UserRepoTest struct {
	suite.Suite
	mockLg *zap.Logger
}

func (u *UserRepoTest) BeforeAll(t provider.T) {
	t.Log("Init mocks")
	u.mockLg = pkg.CreateMockLogger()
}

func (u *UserRepoTest) TestNormalCreateUser(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db40")
	user := domain.User{
		UserID:   userID,
		Mail:     "test@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}
	poolMock.EXPECT().Exec(context.Background(), gomock.Any(), user.UserID,
		user.Mail, user.Password, user.Role).Return(pgconn.CommandTag{}, nil)

	err := userRepo.Create(context.Background(), &user, u.mockLg)

	t.Require().Nil(err)
}

func (u *UserRepoTest) TestExistsCreateUser(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24")
	user := domain.User{
		UserID:   userID,
		Mail:     "test@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}
	poolMock.EXPECT().Exec(context.Background(), gomock.Any(), user.UserID,
		user.Mail, user.Password, user.Role).Return(pgconn.CommandTag{}, errors.New("such user exists"))

	err := userRepo.Create(context.Background(), &user, u.mockLg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestContextTimeoutCreateUser(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)
	poolMock.EXPECT().Exec(ctx, gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{},
		errors.New("expired context"))

	err := userRepo.Create(ctx, &domain.User{}, u.mockLg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestNormalDeleteByID(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db41")
	user := domain.User{
		UserID:   userID,
		Mail:     "test@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}
	poolMock.EXPECT().Exec(context.Background(), gomock.Any(),
		userID).Return(pgconn.CommandTag{}, nil)

	err := userRepo.DeleteByID(context.Background(), user.UserID, u.mockLg)

	t.Require().Nil(err)
}

func (u *UserRepoTest) TestContextTimeoutDeleteByID(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)
	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db41")
	poolMock.EXPECT().Exec(ctx, gomock.Any(), userID).Return(pgconn.CommandTag{},
		errors.New("expired context"))

	err := userRepo.DeleteByID(ctx, userID, u.mockLg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestNormalUpdateUser(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db27")
	user := domain.User{
		UserID:   userID,
		Mail:     "newmail@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}
	poolMock.EXPECT().Exec(context.Background(), gomock.Any(), user.UserID,
		user.Mail, user.Password, user.Role).Return(pgconn.CommandTag{}, nil)

	err := userRepo.Update(context.Background(), &user, u.mockLg)

	t.Require().Nil(err)
}

func (u *UserRepoTest) TestContextTimeoutUpdateUser(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24")
	user := domain.User{
		UserID:   userID,
		Mail:     "newmail@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}
	poolMock.EXPECT().Exec(ctx, gomock.Any(), user.UserID, user.Mail,
		user.Password, user.Role).Return(pgconn.CommandTag{}, errors.New("expired context"))

	err := userRepo.Update(ctx, &user, u.mockLg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestNormalGetByIDUser(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowMock := mock_domain.NewMockRow(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24")
	user := domain.User{
		UserID:   userID,
		Mail:     "user1@mail.ru",
		Password: "password1",
		Role:     domain.Client,
	}
	poolMock.EXPECT().QueryRow(context.Background(), gomock.Any(), user.UserID).Return(rowMock)
	rowMock.EXPECT().Scan(gomock.Any()).Return(nil)

	_, err := userRepo.GetByID(context.Background(), userID, u.mockLg)

	t.Require().Nil(err)
}

func (u *UserRepoTest) TestNoExistsGetByIDUser(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowMock := mock_domain.NewMockRow(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db44")
	poolMock.EXPECT().QueryRow(context.Background(), gomock.Any(), userID).Return(rowMock)
	rowMock.EXPECT().Scan(gomock.Any()).Return(errors.New("empty scan error"))

	usr, err := userRepo.GetByID(context.Background(), userID, u.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.User{}, usr)
}

func (u *UserRepoTest) TestContextTimeoutGetByIDUser(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowMock := mock_domain.NewMockRow(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db44")
	poolMock.EXPECT().QueryRow(ctx, gomock.Any(), userID).Return(rowMock)
	rowMock.EXPECT().Scan(gomock.Any()).Return(errors.New("expired context"))

	usr, err := userRepo.GetByID(ctx, userID, u.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.User{}, usr)
}

func (u *UserRepoTest) TestNormalGetAllUser(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowsMock := mock_domain.NewMockRows(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)

	poolMock.EXPECT().Query(context.Background(), gomock.Any(), 5, 0).Return(rowsMock, nil)
	rowsMock.EXPECT().Next()
	rowsMock.EXPECT().Close()

	_, err := userRepo.GetAll(context.Background(), 0, 5, u.mockLg)

	t.Require().Nil(err)
}

func (u *UserRepoTest) TestContextTimeoutGetAllUser(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowsMock := mock_domain.NewMockRows(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(poolMock, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)
	poolMock.EXPECT().Query(ctx, gomock.Any(), 5, 0).Return(rowsMock, errors.New("expired context"))
	rowsMock.EXPECT().Close().MinTimes(1)

	users, err := userRepo.GetAll(ctx, 0, 5, u.mockLg)

	t.Require().Error(err)
	t.Require().Equal(0, len(users))
}

func TestUserSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(UserRepoTest))
}
