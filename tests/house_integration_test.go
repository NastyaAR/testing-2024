//go:build integration
// +build integration

package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/ports"
	"avito-test-task/internal/repo"
	"avito-test-task/internal/usecase"
	"avito-test-task/pkg"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

type HouseIntegrationTest struct {
	suite.Suite
	houseUsecase domain.HouseUsecase
	houseRepo    domain.HouseRepo
	notifyRepo   domain.NotifyRepo
	notifySender domain.NotifySender
	db           repo.IPool
	mockLg       *zap.Logger
	skipped      bool
}

func (h *HouseIntegrationTest) BeforeAll(t provider.T) {
	connString := "postgres://test-user:test-password@127.0.0.1:5431/test-db?sslmode=disable"

	var err error
	h.db, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf("error while connecting to db: %v", err.Error())
	}
	h.notifyRepo = repo.NewPostgresNotifyRepo(h.db, nil)
	h.notifySender = ports.NewSender()
	h.houseRepo = repo.NewPostgresHouseRepo(h.db, nil)

	done := make(chan bool, 1)
	h.mockLg = pkg.CreateMockLogger()
	h.houseUsecase = usecase.NewHouseUsecase(h.houseRepo, h.notifySender, h.notifyRepo, done, time.Minute, time.Minute, h.mockLg)

	args := os.Args
	for _, arg := range args {
		if arg == "skipped" {
			h.skipped = true
		}
	}
}

func (h *HouseIntegrationTest) AfterAll(t provider.T) {
	h.db.Close()
}

func (h *HouseIntegrationTest) TestNormalCreate(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	houseReq := domain.CreateHouseRequest{
		HomeID:    0,
		Address:   "ул. Спортивная, д. 11",
		Year:      2020,
		Developer: "ООО Комфорт-плюс",
	}
	expected := domain.CreateHouseResponse{
		Address:   "ул. Спортивная, д. 11",
		Year:      2020,
		Developer: "ООО Комфорт-плюс",
	}

	created, err := h.houseUsecase.Create(context.Background(), &houseReq, h.mockLg)

	t.Require().Nil(err)
	t.Require().Equal(expected.Address, created.Address)
	t.Require().Equal(expected.Year, created.Year)
	t.Require().Equal(expected.Developer, created.Developer)
}

func (h *HouseIntegrationTest) TestBadNilRequestCreate(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	created, err := h.houseUsecase.Create(context.Background(), nil, h.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateHouseResponse{}, created)
}

func (h *HouseIntegrationTest) TestBadYearCreate(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	houseReq := domain.CreateHouseRequest{
		HomeID:    0,
		Address:   "ул. Спортивная, д. 11",
		Year:      -2020,
		Developer: "ООО Комфорт-плюс",
	}

	created, err := h.houseUsecase.Create(context.Background(), &houseReq, h.mockLg)

	t.Require().ErrorIs(err, domain.ErrHouse_BadYear)
	t.Require().Equal(domain.CreateHouseResponse{}, created)
}

func (h *HouseIntegrationTest) TestBadDeveloperCreate(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	houseReq := domain.CreateHouseRequest{
		HomeID:    0,
		Address:   "ул. Спортивная, д. 11",
		Year:      2020,
		Developer: "",
	}

	created, err := h.houseUsecase.Create(context.Background(), &houseReq, h.mockLg)

	t.Require().ErrorIs(err, domain.ErrHouse_BadDeveloper)
	t.Require().Equal(domain.CreateHouseResponse{}, created)
}

func (h *HouseIntegrationTest) TestBadAddressCreate(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	houseReq := domain.CreateHouseRequest{
		HomeID:    0,
		Address:   "",
		Year:      2020,
		Developer: "ООО Комфорт-плюс",
	}

	created, err := h.houseUsecase.Create(context.Background(), &houseReq, h.mockLg)

	t.Require().ErrorIs(err, domain.ErrHouse_BadAddress)
	t.Require().Equal(domain.CreateHouseResponse{}, created)
}

func (h *HouseIntegrationTest) TestNormalGetFlatsByHouseID(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	expectedFlats := []domain.SingleFlatResponse{
		domain.SingleFlatResponse{
			ID:      8,
			HouseID: 4,
			Price:   450,
			Rooms:   4,
			Status:  domain.ModeratingStatus,
		},
	}

	expected := domain.FlatsByHouseResponse{Flats: expectedFlats}

	flats, err := h.houseUsecase.GetFlatsByHouseID(context.Background(), 4, domain.ModeratingStatus, h.mockLg)

	t.Require().Nil(err)
	t.Require().Equal(expected, flats)
}

func (h *HouseIntegrationTest) TestBadIDGetFlatsByHouseID(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	flats, err := h.houseUsecase.GetFlatsByHouseID(context.Background(), -1, domain.ModeratingStatus, h.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.FlatsByHouseResponse{}, flats)
}

func (h *HouseIntegrationTest) TestBadStatusGetFlatsByHouseID(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	flats, err := h.houseUsecase.GetFlatsByHouseID(context.Background(), 1, "bad status", h.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.FlatsByHouseResponse{}, flats)
}

func (h *HouseIntegrationTest) TestNormalSubscribeByID(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	err := h.houseUsecase.SubscribeByID(context.Background(), 1, uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24"), h.mockLg)

	t.Require().Nil(err)
}

func (h *HouseIntegrationTest) TestBadIDSubscribeByID(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	err := h.houseUsecase.SubscribeByID(context.Background(), -1, uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db24"), h.mockLg)

	t.Require().Error(err)
}

func (h *HouseIntegrationTest) TestBadUserIDSubscribeByID(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	err := h.houseUsecase.SubscribeByID(context.Background(), -1,
		uuid.Nil, h.mockLg)

	t.Require().Error(err)
}

func (h *HouseIntegrationTest) TestNotExistUserIDSubscribeByID(t provider.T) {
	if h.skipped {
		t.Skip()
	}

	err := h.houseUsecase.SubscribeByID(context.Background(), -1,
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db70"), h.mockLg)

	t.Require().Error(err)
}

func TestHouseIntegrationSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(HouseIntegrationTest))
}
