package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/repo"
	"avito-test-task/pkg"
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

type UserRepoTest struct {
	suite.Suite
	pool     *pgxpool.Pool
	migrator *migrate.Migrate
	mockLg   *zap.Logger
}

func (u *UserRepoTest) BeforeAll(t provider.T) {
	t.Log("Init database connection and migrator")

	var err error
	connString := "postgres://test-user:test-password@localhost:5431/test-db?sslmode=disable"

	u.pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	u.migrator, err = migrate.New("file://../test_migrations", connString)
	err = u.migrator.Up()
	if err != nil {
		log.Fatal(err)
	}
	u.mockLg = pkg.CreateMockLogger()
}

func (u *UserRepoTest) AfterAll(t provider.T) {
	t.Log("Close database connection")
	err := u.migrator.Down()
	if err != nil {
		log.Fatal(err)
	}
	u.pool.Close()
}

func (u *UserRepoTest) TestNormalCreateUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db40")
	user := domain.User{
		UserID:   userID,
		Mail:     "test@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}

	err := userRepo.Create(context.Background(), &user, u.mockLg)

	t.Require().Nil(err)
	usr, _ := userRepo.GetByID(context.Background(), userID, u.mockLg)
	t.Require().Equal(user, usr)
	userRepo.DeleteByID(context.Background(), userID, u.mockLg)
}

func (u *UserRepoTest) TestExistsCreateUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24")
	user := domain.User{
		UserID:   userID,
		Mail:     "test@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}

	err := userRepo.Create(context.Background(), &user, u.mockLg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestContextTimeoutCreateUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	err := userRepo.Create(ctx, &domain.User{}, u.mockLg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestNormalDeleteByID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db41")
	user := domain.User{
		UserID:   userID,
		Mail:     "test@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}
	_ = userRepo.Create(context.Background(), &user, u.mockLg)

	err := userRepo.DeleteByID(context.Background(), userID, u.mockLg)

	t.Require().Nil(err)
	usr, err := userRepo.GetByID(context.Background(), userID, u.mockLg)
	t.Require().Equal(domain.User{}, usr)
}

func (u *UserRepoTest) TestContextTimeoutDeleteByID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)
	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db41")

	err := userRepo.DeleteByID(ctx, userID, u.mockLg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestNormalUpdateUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db27")
	user := domain.User{
		UserID:   userID,
		Mail:     "newmail@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}

	err := userRepo.Update(context.Background(), &user, u.mockLg)

	t.Require().Nil(err)
	usr, _ := userRepo.GetByID(context.Background(), userID, u.mockLg)
	t.Require().Equal(user, usr)
}

func (u *UserRepoTest) TestContextTimeoutUpdateUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24")
	user := domain.User{
		UserID:   userID,
		Mail:     "newmail@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}

	err := userRepo.Update(ctx, &user, u.mockLg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestNormalGetByIDUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24")
	user := domain.User{
		UserID:   userID,
		Mail:     "user1@mail.ru",
		Password: "password1",
		Role:     domain.Client,
	}

	usr, err := userRepo.GetByID(context.Background(), userID, u.mockLg)

	t.Require().Nil(err)
	t.Require().Equal(user, usr)
}

func (u *UserRepoTest) TestNoExistsGetByIDUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db44")

	usr, err := userRepo.GetByID(context.Background(), userID, u.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.User{}, usr)
}

func (u *UserRepoTest) TestContextTimeoutGetByIDUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db44")

	usr, err := userRepo.GetByID(ctx, userID, u.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.User{}, usr)
}

func (u *UserRepoTest) TestNormalGetAllUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)

	expected := []domain.User{
		{
			UserID:   uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db22"),
			Mail:     "test@mail.ru",
			Password: "password",
			Role:     "client",
		},
		{
			UserID:   uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db23"),
			Mail:     "test@mail.ru",
			Password: "password",
			Role:     "moderator",
		},
		{
			UserID:   uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24"),
			Mail:     "user1@mail.ru",
			Password: "password1",
			Role:     "client",
		},
		{
			UserID:   uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db25"),
			Mail:     "user3@mail.ru",
			Password: "password3",
			Role:     "client",
		},
		{
			UserID:   uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db26"),
			Mail:     "user6@mail.ru",
			Password: "password6",
			Role:     "client",
		}}

	users, err := userRepo.GetAll(context.Background(), 0, 5, u.mockLg)

	t.Require().Nil(err)
	t.Require().Equal(expected, users)
}

func (u *UserRepoTest) TestContextTimeoutGetAllUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	users, err := userRepo.GetAll(ctx, 0, 5, u.mockLg)

	t.Require().Error(err)
	t.Require().Equal(0, len(users))
}

func TestUserSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(UserRepoTest))
}
