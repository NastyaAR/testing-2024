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

type HouseRepoTest struct {
	suite.Suite
	mockLg *zap.Logger
}

func (f *HouseRepoTest) BeforeAll(t provider.T) {
	t.Log("Init log")
	f.mockLg = pkg.CreateMockLogger()
}

func (h *HouseRepoTest) TestNormalCreateHouse(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowMock := mock_domain.NewMockRow(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)

	now := time.Now().UTC().Truncate(time.Second)
	house := domain.House{
		HouseID:         11,
		Address:         "ул Спортивная, д 5",
		ConstructYear:   2012,
		Developer:       "ООО НСК",
		CreateHouseDate: now,
		UpdateFlatDate:  now,
	}
	poolMock.EXPECT().QueryRow(context.Background(), gomock.Any(),
		house.Address, house.ConstructYear, house.Developer, house.CreateHouseDate, house.UpdateFlatDate).Return(rowMock)
	rowMock.EXPECT().Scan(gomock.Any()).Return(nil)

	_, err := houseRepo.Create(context.Background(), &house, h.mockLg)

	t.Require().Nil(err)
}

func (h *HouseRepoTest) TestContextTimeoutCreateHouse(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowMock := mock_domain.NewMockRow(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)
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
	poolMock.EXPECT().QueryRow(ctx, gomock.Any(),
		house.Address, house.ConstructYear, house.Developer, house.CreateHouseDate, house.UpdateFlatDate).Return(rowMock)
	rowMock.EXPECT().Scan(gomock.Any()).Return(errors.New("expired context"))

	_, err := houseRepo.Create(ctx, &house, h.mockLg)

	t.Require().Error(err)

}

func (h *HouseRepoTest) TestNormalDeleteByIdHouse(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)

	now := time.Now().UTC().Truncate(time.Second)
	house := domain.House{
		HouseID:         11,
		Address:         "ул Спортивная, д 5",
		ConstructYear:   2012,
		Developer:       "ООО НСК",
		CreateHouseDate: now,
		UpdateFlatDate:  now,
	}
	poolMock.EXPECT().Exec(context.Background(), gomock.Any(), house.HouseID).Return(pgconn.CommandTag{}, nil)

	err := houseRepo.DeleteByID(context.Background(), house.HouseID, h.mockLg)

	t.Require().Nil(err)
}

func (h *HouseRepoTest) TestContextTimeoutDeleteByIdHouse(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)
	poolMock.EXPECT().Exec(ctx, gomock.Any(), 10).Return(pgconn.CommandTag{}, errors.New("expired context"))

	err := houseRepo.DeleteByID(ctx, 10, h.mockLg)

	t.Require().Error(err)
}

func (h *HouseRepoTest) TestNormalUpdateHouse(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)

	now := time.Now().UTC().Truncate(time.Second)
	house := domain.House{
		HouseID:         10,
		Address:         "ул Спортивная, д 5",
		ConstructYear:   2012,
		Developer:       "ООО НСК",
		CreateHouseDate: now,
		UpdateFlatDate:  now,
	}
	poolMock.EXPECT().Exec(context.Background(), gomock.Any(),
		house.HouseID, house.Address, house.ConstructYear, house.Developer, house.CreateHouseDate,
		house.UpdateFlatDate).Return(pgconn.CommandTag{}, nil)

	err := houseRepo.Update(context.Background(), &house, h.mockLg)

	t.Require().Nil(err)
}

func (h *HouseRepoTest) TestContextTimeoutUpdateHouse(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)
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
	poolMock.EXPECT().Exec(ctx, gomock.Any(), house.HouseID, house.Address, house.ConstructYear,
		house.Developer, house.CreateHouseDate, house.UpdateFlatDate).Return(pgconn.CommandTag{}, errors.New("expired context"))

	err := houseRepo.Update(ctx, &house, h.mockLg)

	t.Require().Error(err)
}

func (h *HouseRepoTest) TestNormalGetByID(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowMock := mock_domain.NewMockRow(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)
	poolMock.EXPECT().QueryRow(context.Background(), gomock.Any(), 10).Return(rowMock)
	rowMock.EXPECT().Scan(gomock.Any()).Return(nil)

	_, err := houseRepo.GetByID(context.Background(), 10, h.mockLg)

	t.Require().Nil(err)
}

func (h *HouseRepoTest) TestContextTimeoutGetByID(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowMock := mock_domain.NewMockRow(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)
	poolMock.EXPECT().QueryRow(ctx, gomock.Any(), 10).Return(rowMock)
	rowMock.EXPECT().Scan(gomock.Any()).Return(errors.New("expired context"))

	_, err := houseRepo.GetByID(ctx, 10, h.mockLg)

	t.Require().Error(err)
}

func (h *HouseRepoTest) TestNormalGetAll(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowsMock := mock_domain.NewMockRows(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)
	now := time.Now().UTC().Truncate(time.Second)
	_ = []domain.House{
		{HouseID: 1, Address: "ул. Спортивная, д. 1", ConstructYear: 2021, Developer: "OOO Строй", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 2, Address: "ул. Спортивная, д. 2", ConstructYear: 2020, Developer: "ЗАО Строительство", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 3, Address: "ул. Спортивная, д. 3", ConstructYear: 2019, Developer: "ИП Строитель", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 4, Address: "ул. Спортивная, д. 4", ConstructYear: 2022, Developer: "ЗАО Новострой", CreateHouseDate: now, UpdateFlatDate: now},
		{HouseID: 5, Address: "ул. Спортивная, д. 5", ConstructYear: 2021, Developer: "OOO Строй", CreateHouseDate: now, UpdateFlatDate: now},
	}
	poolMock.EXPECT().Query(context.Background(), gomock.Any(), 5, 0).Return(rowsMock, nil)
	rowsMock.EXPECT().Next()
	rowsMock.EXPECT().Close().AnyTimes()

	_, err := houseRepo.GetAll(context.Background(), 0, 5, h.mockLg)

	t.Require().Nil(err)
}

func (h *HouseRepoTest) TestContextTimeoutGetAll(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowsMock := mock_domain.NewMockRows(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)
	poolMock.EXPECT().Query(ctx, gomock.Any(), 10, 0).Return(rowsMock, errors.New("expired context"))
	rowsMock.EXPECT().Close()

	_, err := houseRepo.GetAll(ctx, 0, 10, h.mockLg)

	t.Require().Error(err)
}

func (h *HouseRepoTest) TestNormalNonModeratingGetFlatsByHouseID(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowsMock := mock_domain.NewMockRows(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)
	_ = []domain.Flat{
		{ID: 10, HouseID: 1, UserID: uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db22"), Price: 100, Rooms: 2, Status: "created", ModeratorID: 0},
		{ID: 1, HouseID: 1, UserID: uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db22"), Price: 100, Rooms: 2, Status: "created", ModeratorID: 0},
		{ID: 2, HouseID: 1, UserID: uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24"), Price: 150, Rooms: 3, Status: "approved", ModeratorID: 0},
		{ID: 3, HouseID: 1, UserID: uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24"), Price: 200, Rooms: 2, Status: "declined", ModeratorID: 0},
	}

	poolMock.EXPECT().Query(context.Background(), gomock.Any(), 1).Return(rowsMock, nil)
	rowsMock.EXPECT().Next()
	rowsMock.EXPECT().Close().MinTimes(1)

	_, err := houseRepo.GetFlatsByHouseID(context.Background(), 1, domain.AnyStatus, h.mockLg)

	t.Require().Nil(err)
}

func (h *HouseRepoTest) TestNormalModeratingGetFlatsByHouseID(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowsMock := mock_domain.NewMockRows(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)

	poolMock.EXPECT().Query(context.Background(), gomock.Any(), 2, domain.ModeratingStatus).Return(rowsMock, nil)
	rowsMock.EXPECT().Next()
	rowsMock.EXPECT().Close()

	_, err := houseRepo.GetFlatsByHouseID(context.Background(), 2, domain.ModeratingStatus, h.mockLg)

	t.Require().Nil(err)
}

func (h *HouseRepoTest) TestContextTimeoutGetFlatsByHouseID(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowsMock := mock_domain.NewMockRows(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)

	poolMock.EXPECT().Query(ctx, gomock.Any(), 2, domain.ModeratingStatus).Return(rowsMock, errors.New("expired context"))
	rowsMock.EXPECT().Close()

	_, err := houseRepo.GetFlatsByHouseID(ctx, 2, domain.ModeratingStatus, h.mockLg)

	t.Require().Error(err)
}

func (h *HouseRepoTest) TestNormalSubscribeByID(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	houseRepo := repo.NewPostgresHouseRepo(poolMock, retryAdapter)

	poolMock.EXPECT().Exec(context.Background(), gomock.Any(), uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24"), 1).Return(pgconn.CommandTag{}, nil)

	err := houseRepo.SubscribeByID(context.Background(), 1, uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24"), h.mockLg)

	t.Require().Nil(err)
}

func TestHouseSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(HouseRepoTest))
}
