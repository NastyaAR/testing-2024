package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/repo"
	"avito-test-task/pkg"
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"log"
	"sync"
	"testing"
	"time"
)

type HouseRepoTest struct {
	suite.Suite
	pool     *pgxpool.Pool
	migrator *migrate.Migrate
	mtx      sync.Mutex
}

func (f *HouseRepoTest) BeforeAll(t provider.T) {
	t.Log("Init database connection and migrator")

	var err error
	connString := "postgres://test-user:test-password@localhost:5431/test-db?sslmode=disable"

	f.pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	f.migrator, err = migrate.New("file://../test_migrations", connString)
}

func (s *HouseRepoTest) BeforeEach(t provider.T) {
	t.Log("Up Migration")
	s.mtx.Lock()
	err := s.migrator.Up()
	s.mtx.Unlock()
	t.Log(err)
}

func (s *HouseRepoTest) AfterEach(t provider.T) {
	t.Log("Down Migration")
	s.mtx.Lock()
	err := s.migrator.Down()
	s.mtx.Unlock()
	t.Log(err)
}

func (s *HouseRepoTest) AfterAll(t provider.T) {
	t.Log("Close database connection")
	s.pool.Close()

}

func (h *HouseRepoTest) TestNormalCreateHouse(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	now := time.Now().UTC().Truncate(time.Second)
	house := domain.House{
		HouseID:         11,
		Address:         "ул Спортивная, д 5",
		ConstructYear:   2012,
		Developer:       "ООО НСК",
		CreateHouseDate: now,
		UpdateFlatDate:  now,
	}

	created, err := houseRepo.Create(context.Background(), &house, lg)

	t.Require().Nil(err)
	t.Require().Equal(house, created)
}

func (h *HouseRepoTest) TestContextTimeoutCreateHouse(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	now := time.Now().UTC().Truncate(time.Second)
	house := domain.House{
		HouseID:         11,
		Address:         "ул Спортивная, д 5",
		ConstructYear:   2012,
		Developer:       "ООО НСК",
		CreateHouseDate: now,
		UpdateFlatDate:  now,
	}

	created, err := houseRepo.Create(ctx, &house, lg)

	t.Require().Error(err)
	t.Require().Equal(domain.House{}, created)
}

func (h *HouseRepoTest) TestNormalDeleteByIdHouse(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	err := houseRepo.DeleteByID(context.Background(), 10, lg)

	t.Require().Nil(err)
	house, err := houseRepo.GetByID(context.Background(), 10, lg)
	t.Require().Error(err)
	t.Require().Equal(domain.House{}, house)
}

func (h *HouseRepoTest) TestContextTimeoutDeleteByIdHouse(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	err := houseRepo.DeleteByID(ctx, 10, lg)

	t.Require().Error(err)
}

func (h *HouseRepoTest) TestNoExistDeleteByIdHouse(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	err := houseRepo.DeleteByID(context.Background(), 11, lg)

	t.Require().Nil(err)
}

func (h *HouseRepoTest) TestNormalUpdateHouse(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	now := time.Now().UTC().Truncate(time.Second)
	house := domain.House{
		HouseID:         10,
		Address:         "ул Спортивная, д 5",
		ConstructYear:   2012,
		Developer:       "ООО НСК",
		CreateHouseDate: now,
		UpdateFlatDate:  now,
	}

	err := houseRepo.Update(context.Background(), &house, lg)

	t.Require().Nil(err)
	updatedHouse, _ := houseRepo.GetByID(context.Background(), 10, lg)
	t.Require().Equal(house, updatedHouse)
}

func (h *HouseRepoTest) TestContextTimeoutUpdateHouse(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	now := time.Now().UTC().Truncate(time.Second)
	house := domain.House{
		HouseID:         10,
		Address:         "ул Спортивная, д 5",
		ConstructYear:   2012,
		Developer:       "ООО НСК",
		CreateHouseDate: now,
		UpdateFlatDate:  now,
	}

	err := houseRepo.Update(ctx, &house, lg)

	t.Require().Error(err)
}

func (h *HouseRepoTest) TestNoExistUpdateHouse(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	now := time.Now().UTC().Truncate(time.Second)
	house := domain.House{
		HouseID:         11,
		Address:         "ул Спортивная, д 5",
		ConstructYear:   2012,
		Developer:       "ООО НСК",
		CreateHouseDate: now,
		UpdateFlatDate:  now,
	}

	err := houseRepo.Update(context.Background(), &house, lg)

	t.Require().Nil(err)
}

func (h *HouseRepoTest) TestNormalGetByID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	_, err := houseRepo.GetByID(context.Background(), 10, lg)

	t.Require().Nil(err)
}

func (h *HouseRepoTest) TestNoExistGetByID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	_, err := houseRepo.GetByID(context.Background(), 11, lg)

	t.Require().Error(err)
}

func (h *HouseRepoTest) TestContextTimeoutGetByID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	_, err := houseRepo.GetByID(ctx, 10, lg)

	t.Require().Error(err)
}

func (h *HouseRepoTest) TestNormalGetAll(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	now := time.Now().UTC().Truncate(time.Second)
	houses := []domain.House{
		{HouseID: 1, Address: "ул. Спортивная, д. 1", ConstructYear: 2021, Developer: "OOO Строй", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 2, Address: "ул. Спортивная, д. 2", ConstructYear: 2020, Developer: "ЗАО Строительство", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 3, Address: "ул. Спортивная, д. 3", ConstructYear: 2019, Developer: "ИП Строитель", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 4, Address: "ул. Спортивная, д. 4", ConstructYear: 2022, Developer: "ЗАО Новострой", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 5, Address: "ул. Спортивная, д. 5", ConstructYear: 2021, Developer: "OOO Строй", CreateHouseDate: now, UpdateFlatDate: now},
	}

	actualHouses, err := houseRepo.GetAll(context.Background(), 0, 5, lg)

	t.Require().Nil(err)
	for i := 0; i < len(houses); i++ {
		t.Require().Equal(houses[i].HouseID, actualHouses[i].HouseID)
		t.Require().Equal(houses[i].Developer, actualHouses[i].Developer)
		t.Require().Equal(houses[i].ConstructYear, actualHouses[i].ConstructYear)
		t.Require().Equal(houses[i].Address, actualHouses[i].Address)
	}
}

func (h *HouseRepoTest) TestContextTimeoutGetAll(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	_, err := houseRepo.GetAll(ctx, 0, 10, lg)

	t.Require().Error(err)
}

func TestHouseSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(HouseRepoTest))
}
