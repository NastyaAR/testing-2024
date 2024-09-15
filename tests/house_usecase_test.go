package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/ports"
	"avito-test-task/internal/usecase"
	"avito-test-task/pkg"
	mock_domain "avito-test-task/tests/mocks"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
	"time"
)

type HouseUsecaseTest struct {
	suite.Suite
	houseRepoMock  *mock_domain.MockHouseRepo
	notifyRepoMock *mock_domain.MockNotifyRepo
	lg             *zap.Logger
	flatMother     *FlatMother
}

func (h *HouseUsecaseTest) BeforeAll(t provider.T) {
	t.Log("Init mock")
	ctrl := gomock.NewController(t)
	h.houseRepoMock = mock_domain.NewMockHouseRepo(ctrl)
	h.notifyRepoMock = mock_domain.NewMockNotifyRepo(ctrl)
	h.lg, _ = pkg.CreateLogger("../log.log", "prod")
	h.flatMother = new(FlatMother)
}

func (h *HouseUsecaseTest) TestNormalCreateHouse(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	req := domain.CreateHouseRequest{
		HomeID:    2,
		Address:   "ул. Тестовая, д. 3",
		Year:      2021,
		Developer: "ООО ТестСтрой",
	}

	resp := domain.CreateHouseResponse{
		HomeID:    2,
		Address:   "ул. Тестовая, д. 3",
		Year:      2021,
		Developer: "ООО ТестСтрой",
		CreatedAt: "",
		UpdateAt:  "",
	}

	house := domain.House{
		HouseID:         2,
		Address:         "ул. Тестовая, д. 3",
		ConstructYear:   2021,
		Developer:       "ООО ТестСтрой",
		CreateHouseDate: time.Time{},
		UpdateFlatDate:  time.Time{},
	}

	h.houseRepoMock.EXPECT().Create(context.Background(), gomock.Any(), h.lg).Return(house, nil)

	created, err := houseUsecase.Create(context.Background(), &req, h.lg)

	t.Require().Nil(err)
	t.Require().Equal(resp.HomeID, created.HomeID)
	t.Require().Equal(resp.Address, created.Address)
	t.Require().Equal(resp.Year, created.Year)
	t.Require().Equal(resp.Developer, created.Developer)
}

func (h *HouseUsecaseTest) TestNilRequestCreateHouse(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	created, err := houseUsecase.Create(context.Background(), nil, h.lg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateHouseResponse{}, created)
}

func (h *HouseUsecaseTest) TestBadYearCreateHouse(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	req := domain.CreateHouseRequest{
		HomeID:    0,
		Address:   "ул. Тестовая, д. 3",
		Year:      -2021,
		Developer: "ООО ТестСтрой",
	}

	created, err := houseUsecase.Create(context.Background(), &req, h.lg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateHouseResponse{}, created)
}

func (h *HouseUsecaseTest) TestBadDeveloperCreateHouse(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	req := domain.CreateHouseRequest{
		HomeID:    0,
		Address:   "ул. Тестовая, д. 3",
		Year:      2021,
		Developer: "",
	}

	created, err := houseUsecase.Create(context.Background(), &req, h.lg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateHouseResponse{}, created)
}

func (h *HouseUsecaseTest) TestBadAddressCreateHouse(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	req := domain.CreateHouseRequest{
		HomeID:    0,
		Address:   "",
		Year:      2021,
		Developer: "OOO ТестСтрой",
	}

	created, err := houseUsecase.Create(context.Background(), &req, h.lg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateHouseResponse{}, created)
}

func (h *HouseUsecaseTest) TestBadRepoCallCreateHouse(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	req := domain.CreateHouseRequest{
		HomeID:    2,
		Address:   "ул. Тестовая, д. 3",
		Year:      2021,
		Developer: "ООО ТестСтрой",
	}

	h.houseRepoMock.EXPECT().Create(context.Background(), gomock.Any(), h.lg).Return(domain.House{}, errors.New("error"))

	created, err := houseUsecase.Create(context.Background(), &req, h.lg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateHouseResponse{}, created)
}

func (h *HouseUsecaseTest) TestNormalGetFlatsByHouseID(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	flats := []domain.Flat{}
	singleFlats := []domain.SingleFlatResponse{}
	for i := 1; i < 6; i++ {
		newFlat := h.flatMother.DefaultFlat(i, 10)
		flats = append(flats, newFlat)
		singleFlatResponse := h.flatMother.DefaultSingleFlatResponse(&newFlat)
		singleFlats = append(singleFlats, singleFlatResponse)
	}

	resp := domain.FlatsByHouseResponse{singleFlats}

	h.houseRepoMock.EXPECT().GetFlatsByHouseID(context.Background(), 10, domain.CreatedStatus, h.lg).Return(flats, nil)

	foundFlats, err := houseUsecase.GetFlatsByHouseID(context.Background(), 10, domain.CreatedStatus, h.lg)

	t.Require().Nil(err)
	t.Require().Equal(resp, foundFlats)
}

func (h *HouseUsecaseTest) TestBadRepoCallGetFlatsByHouseID(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	flats := []domain.Flat{}
	for i := 1; i < 6; i++ {
		flats = append(flats, h.flatMother.DefaultFlat(i, 10))
	}

	h.houseRepoMock.EXPECT().GetFlatsByHouseID(context.Background(), 10, domain.CreatedStatus, h.lg).Return([]domain.Flat{}, errors.New("error"))

	foundFlats, err := houseUsecase.GetFlatsByHouseID(context.Background(), 10, domain.CreatedStatus, h.lg)

	t.Require().Error(err)
	t.Require().Equal(domain.FlatsByHouseResponse{}, foundFlats)
}

func (h *HouseUsecaseTest) TestBadIDGetFlatsByHouseID(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	foundedFlats, err := houseUsecase.GetFlatsByHouseID(context.Background(), -1, domain.CreatedStatus, h.lg)

	t.Require().Error(err)
	t.Require().Equal(domain.FlatsByHouseResponse{}, foundedFlats)
}

func (h *HouseUsecaseTest) TestBadStatusGetFlatsByHouseID(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	foundedFlats, err := houseUsecase.GetFlatsByHouseID(context.Background(), 10, "test", h.lg)

	t.Require().Error(err)
	t.Require().Equal(domain.FlatsByHouseResponse{}, foundedFlats)
}

func (h *HouseUsecaseTest) TestNormalSubscribeByID(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	uid := uuid.New()
	h.houseRepoMock.EXPECT().SubscribeByID(context.Background(), 1, uid, h.lg).Return(nil)

	err := houseUsecase.SubscribeByID(context.Background(), 1, uid, h.lg)

	t.Require().Nil(err)
}

func (h *HouseUsecaseTest) TestBadIDSubscribeByID(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	uid := uuid.Nil
	err := houseUsecase.SubscribeByID(context.Background(), 1, uid, h.lg)

	t.Require().Error(err)
}

func (h *HouseUsecaseTest) TestBadRepoCallSubscribeByID(t provider.T) {
	notifySender := ports.NewSender()
	done := make(chan bool, 1)
	done <- true
	houseUsecase := usecase.NewHouseUsecase(h.houseRepoMock, notifySender, h.notifyRepoMock, done, time.Second, time.Second, h.lg)

	uid := uuid.New()
	h.houseRepoMock.EXPECT().SubscribeByID(context.Background(), 1, uid, h.lg).Return(errors.New("error"))

	err := houseUsecase.SubscribeByID(context.Background(), 1, uid, h.lg)

	t.Require().Error(err)
}

func TestHouseUsecaseSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(HouseUsecaseTest))
}
