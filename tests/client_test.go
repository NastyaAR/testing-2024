//go:build e2e
// +build e2e

package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/ports"
	"avito-test-task/internal/repo"
	"avito-test-task/internal/usecase"
	"avito-test-task/pkg"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

type ClientTest struct {
	suite.Suite
	userUsecase  domain.UserUsecase
	flatUsecase  domain.FlatUsecase
	houseUsecase domain.HouseUsecase
	notifySender domain.NotifySender
	userRepo     domain.UserRepo
	flatRepo     domain.FlatRepo
	houseRepo    domain.HouseRepo
	notifyRepo   domain.NotifyRepo
	skipped      bool

	db   repo.IPool
	lg   *zap.Logger
	done chan bool
}

func (c *ClientTest) BeforeAll(t provider.T) {
	host := os.Getenv("POSTGRES_TEST_HOST")
	port := os.Getenv("POSTGRES_TEST_PORT")
	connString := "postgres://test-user:test-password@" + host + ":" + port + "/test-db?sslmode=disable"

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c.db, err = pgxpool.New(ctx, connString)
	if err != nil {
		t.Fatalf("error while connecting to db: %v", err.Error())
	}

	c.lg, _ = pkg.CreateLogger("./log.log", "debug")

	c.userRepo = repo.NewPostrgesUserRepo(c.db, nil)
	c.flatRepo = repo.NewPostgresFlatRepo(c.db, nil)
	c.houseRepo = repo.NewPostgresHouseRepo(c.db, nil)
	c.notifyRepo = repo.NewPostgresNotifyRepo(c.db, nil)

	c.userUsecase = usecase.NewUserUsecase(c.userRepo)
	c.flatUsecase = usecase.NewFlatUsecase(c.flatRepo)
	c.notifySender = ports.NewSender()

	c.done = make(chan bool, 1)
	c.done <- true
	c.houseUsecase = usecase.NewHouseUsecase(c.houseRepo, c.notifySender,
		c.notifyRepo, c.done, time.Second, time.Second, c.lg)

	args := os.Args
	for _, arg := range args {
		if arg == "skipped" {
			c.skipped = true
		}
	}
}

func (c *ClientTest) AfterAll(t provider.T) {
	c.db.Close()
}

func (c *ClientTest) TestCreateNewFlat(t provider.T) {
	if c.skipped {
		t.Skip()
	}

	registerRequest := domain.RegisterUserRequest{
		Email:    "masha@mail.ru",
		Password: "12345",
		UserType: domain.Client,
	}

	ctx := context.Background()

	registerResponse, registerErr := c.userUsecase.Register(ctx, &registerRequest, c.lg)
	loginRequest := domain.LoginUserRequest{
		ID:       registerResponse.UserID,
		Password: registerRequest.Password,
	}

	_, loginErr := c.userUsecase.Login(ctx, &loginRequest, c.lg)

	createRequest := domain.CreateFlatRequest{
		FlatID:  100,
		HouseID: 6,
		Price:   100000,
		Rooms:   3,
	}

	expected := domain.Flat{
		ID:          100,
		HouseID:     6,
		UserID:      registerResponse.UserID,
		Price:       100000,
		Rooms:       3,
		Status:      domain.CreatedStatus,
		ModeratorID: 0,
	}

	_, createdErr := c.flatUsecase.Create(ctx, registerResponse.UserID,
		&createRequest, c.lg)

	flat, _ := c.flatRepo.GetByID(ctx, createRequest.FlatID, createRequest.HouseID, c.lg)
	t.Require().Nil(registerErr)
	t.Require().Nil(loginErr)
	t.Require().Nil(createdErr)
	t.Require().Equal(expected, flat)
}

func TestClientSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ClientTest))
}
