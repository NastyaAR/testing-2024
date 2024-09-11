package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/repo"
	"avito-test-task/pkg"
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"log"
	"sync"
	"testing"
	"time"
)

type FlatRepoTest struct {
	suite.Suite
	pool     *pgxpool.Pool
	migrator *migrate.Migrate
	mtx      sync.Mutex
}

func (f *FlatRepoTest) BeforeAll(t provider.T) {
	t.Log("Init database connection and migrator")

	var err error
	connString := "postgres://test-user:test-password@localhost:5431/test-db?sslmode=disable"

	f.pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	f.migrator, err = migrate.New("file://../test_migrations", connString)
}

func (s *FlatRepoTest) BeforeEach(t provider.T) {
	t.Log("Up Migration")
	s.mtx.Lock()
	err := s.migrator.Up()
	s.mtx.Unlock()
	t.Log(err)
}

func (s *FlatRepoTest) AfterEach(t provider.T) {
	t.Log("Down Migration")
	s.mtx.Lock()
	err := s.migrator.Down()
	s.mtx.Unlock()
	t.Log(err)
}

func (s *FlatRepoTest) AfterAll(t provider.T) {
	t.Log("Close database connection")
	s.pool.Close()

}

func (f *FlatRepoTest) TestNormalCreateFlat(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(f.pool, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(f.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	userID, _ := uuid.Parse("019126ee-2b7d-758e-bb22-fe2e45b2db22")
	flat := domain.Flat{
		ID:          101,
		HouseID:     1,
		UserID:      userID,
		Price:       10000000,
		Rooms:       2,
		Status:      domain.CreatedStatus,
		ModeratorID: 1,
	}

	flat, err := flatRepo.Create(context.Background(), &flat, lg)

	actual, _ := flatRepo.GetByID(context.Background(), 101, 1, lg)
	t.Require().Nil(err)
	t.Require().Equal(flat, actual)
}

func (f *FlatRepoTest) TestContextTimeoutCreateFlat(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(f.pool, 1, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(f.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	userID, _ := uuid.Parse("019126ee-2b7d-758e-bb22-fe2e45b2db22")
	flat := domain.Flat{
		ID:          0,
		HouseID:     1,
		UserID:      userID,
		Price:       10000000,
		Rooms:       2,
		Status:      domain.CreatedStatus,
		ModeratorID: 1,
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
	time.Sleep(time.Millisecond)

	flat, err := flatRepo.Create(ctx, &flat, lg)

	t.Require().Error(err)
}

func (f *FlatRepoTest) TestNormalDeleteByIdFlat(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(f.pool, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(f.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	err := flatRepo.DeleteByID(context.Background(), 2, 1, lg)

	t.Require().Nil(err)
	flat, _ := flatRepo.GetByID(context.Background(), 2, 1, lg)

	t.Require().Equal(domain.Flat{}, flat)
}

func (f *FlatRepoTest) TestContextTimeoutDeleteByIdFlat(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(f.pool, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(f.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
	time.Sleep(time.Millisecond)

	err := flatRepo.DeleteByID(ctx, 1, 1, lg)

	t.Require().Error(err)
}

func (f *FlatRepoTest) TestNormalGetByIdFlat(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(f.pool, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(f.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	userID, _ := uuid.Parse("019126ee-2b7d-758e-bb22-fe2e45b2db22")
	flat := domain.Flat{
		ID:          1,
		HouseID:     1,
		UserID:      userID,
		Price:       100,
		Rooms:       2,
		Status:      domain.CreatedStatus,
		ModeratorID: 0,
	}

	flat, err := flatRepo.GetByID(context.Background(), 1, 1, lg)

	t.Require().Nil(err)
	t.Require().Equal(flat, flat)
}

func (f *FlatRepoTest) TestContextTimeoutGetByIdFlat(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(f.pool, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(f.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
	time.Sleep(time.Millisecond)

	flat, err := flatRepo.GetByID(ctx, 1, 1, lg)

	t.Require().Error(err)
	t.Require().Equal(domain.Flat{}, flat)
}

func (f *FlatRepoTest) TestNormalGetAllFlat(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(f.pool, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(f.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	flats, err := flatRepo.GetAll(context.Background(), 0, 10, lg)

	t.Require().Nil(err)
	t.Require().Equal(10, len(flats))
}

func (f *FlatRepoTest) TestNormalOutOfRangeGetAllFlat(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(f.pool, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(f.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	flats, err := flatRepo.GetAll(context.Background(), 20, 10, lg)

	t.Require().Nil(err)
	t.Require().Equal(len(flats), 0)
}

func (f *FlatRepoTest) TestContextTimeoutGetAllFlat(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(f.pool, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(f.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	flats, err := flatRepo.GetAll(ctx, 1, 1, lg)

	t.Require().Error(err)
	t.Require().Equal(len(flats), 0)
}

func TestFlatSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(FlatRepoTest))
}
