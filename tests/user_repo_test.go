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
	"log"
	"testing"
	"time"
)

type UserRepoTest struct {
	suite.Suite
	pool     *pgxpool.Pool
	migrator *migrate.Migrate
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
	u.migrator.Up()
}

func (u *UserRepoTest) AfterAll(t provider.T) {
	t.Log("Close database connection")
	u.migrator.Down()
	u.pool.Close()
}

func (u *UserRepoTest) TestNormalCreateUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db40")
	user := domain.User{
		UserID:   userID,
		Mail:     "test@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}

	err := userRepo.Create(context.Background(), &user, lg)

	t.Require().Nil(err)
	usr, _ := userRepo.GetByID(context.Background(), userID, lg)
	t.Require().Equal(user, usr)
	userRepo.DeleteByID(context.Background(), userID, lg)
}

func (u *UserRepoTest) TestExistsCreateUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24")
	user := domain.User{
		UserID:   userID,
		Mail:     "test@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}

	err := userRepo.Create(context.Background(), &user, lg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestContextTimeoutCreateUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	err := userRepo.Create(ctx, &domain.User{}, lg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestNormalDeleteByID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db41")
	user := domain.User{
		UserID:   userID,
		Mail:     "test@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}
	_ = userRepo.Create(context.Background(), &user, lg)

	err := userRepo.DeleteByID(context.Background(), userID, lg)

	t.Require().Nil(err)
	usr, err := userRepo.GetByID(context.Background(), userID, lg)
	t.Require().Equal(domain.User{}, usr)
}

func (u *UserRepoTest) TestContextTimeoutDeleteByID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)
	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db41")

	err := userRepo.DeleteByID(ctx, userID, lg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestNormalUpdateUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24")
	user := domain.User{
		UserID:   userID,
		Mail:     "newmail@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}

	err := userRepo.Update(context.Background(), &user, lg)

	t.Require().Nil(err)
	usr, _ := userRepo.GetByID(context.Background(), userID, lg)
	t.Require().Equal(user, usr)
}

func (u *UserRepoTest) TestContextTimeoutUpdateUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24")
	user := domain.User{
		UserID:   userID,
		Mail:     "newmail@mail.ru",
		Password: "password",
		Role:     domain.Client,
	}

	err := userRepo.Update(ctx, &user, lg)

	t.Require().Error(err)
}

func (u *UserRepoTest) TestNormalGetByIDUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24")
	user := domain.User{
		UserID:   userID,
		Mail:     "user1@mail.ru",
		Password: "password1",
		Role:     domain.Client,
	}

	usr, err := userRepo.GetByID(context.Background(), userID, lg)

	t.Require().Nil(err)
	t.Require().Equal(user, usr)
}

func (u *UserRepoTest) TestNoExistsGetByIDUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db44")

	usr, err := userRepo.GetByID(context.Background(), userID, lg)

	t.Require().Error(err)
	t.Require().Equal(domain.User{}, usr)
}

func (u *UserRepoTest) TestContextTimeoutGetByIDUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db44")

	usr, err := userRepo.GetByID(ctx, userID, lg)

	t.Require().Error(err)
	t.Require().Equal(domain.User{}, usr)
}

func (u *UserRepoTest) TestNormalGetAllUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

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
			UserID:   uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db30"),
			Mail:     "user2@mail.ru",
			Password: "password2",
			Role:     "moderator",
		},
		{
			UserID:   uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db25"),
			Mail:     "user3@mail.ru",
			Password: "password3",
			Role:     "client",
		}}

	users, err := userRepo.GetAll(context.Background(), 0, 5, lg)

	t.Require().Nil(err)
	t.Require().Equal(expected, users)
}

func (u *UserRepoTest) TestContextTimeoutGetAllUser(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(u.pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(u.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	users, err := userRepo.GetAll(ctx, 0, 5, lg)

	t.Require().Error(err)
	t.Require().Equal(0, len(users))
}

func TestUserSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(UserRepoTest))
}
