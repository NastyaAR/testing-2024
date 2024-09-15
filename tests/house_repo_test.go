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

type HouseRepoTest struct {
	suite.Suite
	pool     *pgxpool.Pool
	migrator *migrate.Migrate
}

func IsEqualHouses(expected *domain.House, actual *domain.House) bool {
	return (expected.Developer == actual.Developer &&
		expected.Address == actual.Address && expected.ConstructYear == actual.ConstructYear)
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
	f.migrator.Up()
}

func (s *HouseRepoTest) AfterAll(t provider.T) {
	t.Log("Close database connection")
	s.migrator.Down()
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
	t.Require().Equal(true, IsEqualHouses(&house, &created))
	houseRepo.DeleteByID(context.Background(), created.HouseID, lg)
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

	err = houseRepo.DeleteByID(context.Background(), created.HouseID, lg)

	t.Require().Nil(err)
	house, err = houseRepo.GetByID(context.Background(), created.HouseID, lg)
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
		t.Require().Equal(true, IsEqualHouses(&houses[i], &actualHouses[i]))
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

func (h *HouseRepoTest) TestNormalOffsetGetAll(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	now := time.Now().UTC().Truncate(time.Second)
	houses := []domain.House{
		{HouseID: 2, Address: "ул. Спортивная, д. 2", ConstructYear: 2020, Developer: "ЗАО Строительство", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 3, Address: "ул. Спортивная, д. 3", ConstructYear: 2019, Developer: "ИП Строитель", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 4, Address: "ул. Спортивная, д. 4", ConstructYear: 2022, Developer: "ЗАО Новострой", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 5, Address: "ул. Спортивная, д. 5", ConstructYear: 2021, Developer: "OOO Строй", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 6, Address: "ул. Спортивная, д. 6", ConstructYear: 2023, Developer: "Компания Реал", CreateHouseDate: now, UpdateFlatDate: now},
	}

	actualHouses, err := houseRepo.GetAll(context.Background(), 1, 5, lg)

	t.Require().Nil(err)
	for i := 0; i < len(houses); i++ {
		t.Require().Equal(true, IsEqualHouses(&houses[i], &actualHouses[i]))
	}
}

func (h *HouseRepoTest) TestNormalGetFlatsByHouseID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	flats := []domain.Flat{
		{ID: 10, HouseID: 1, UserID: uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db22"), Price: 100, Rooms: 2, Status: "created", ModeratorID: 0},
		{ID: 1, HouseID: 1, UserID: uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db22"), Price: 100, Rooms: 2, Status: "created", ModeratorID: 0},
		{ID: 2, HouseID: 1, UserID: uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24"), Price: 150, Rooms: 3, Status: "approved", ModeratorID: 0},
		{ID: 3, HouseID: 1, UserID: uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24"), Price: 200, Rooms: 2, Status: "declined", ModeratorID: 0},
	}

	actualFlats, err := houseRepo.GetFlatsByHouseID(context.Background(), 1, domain.AnyStatus, lg)

	t.Require().Nil(err)
	t.Require().Equal(flats, actualFlats)
}

func (h *HouseRepoTest) TestNormalOnModerationGetFlatsByHouseID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	flats := []domain.Flat{
		{ID: 4, HouseID: 2, UserID: uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db25"), Price: 250, Rooms: 4, Status: "on moderation", ModeratorID: 0},
	}

	actualFlats, err := houseRepo.GetFlatsByHouseID(context.Background(), 2, domain.ModeratingStatus, lg)

	t.Require().Nil(err)
	t.Require().Equal(flats, actualFlats)
}

func (h *HouseRepoTest) TestContextTimeoutGetFlatsByHouseID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	_, err := houseRepo.GetFlatsByHouseID(ctx, 2, domain.ModeratingStatus, lg)

	t.Require().Error(err)
}

func (h *HouseRepoTest) TestNormalNoExistsGetFlatsByHouseID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	flats, err := houseRepo.GetFlatsByHouseID(context.Background(), 11, domain.AnyStatus, lg)

	t.Require().Nil(err)
	t.Require().Equal(0, len(flats))
}

func (h *HouseRepoTest) TestNormalSubscribeByID(t provider.T) {
	retryAdapter := repo.NewPostgresRetryAdapter(h.pool, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(h.pool, retryAdapter)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	err := houseRepo.SubscribeByID(context.Background(), 1, uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24"), lg)

	t.Require().Nil(err)
}

func TestHouseSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(HouseRepoTest))
}
