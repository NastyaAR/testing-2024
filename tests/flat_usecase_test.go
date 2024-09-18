package tests

import (
	"avito-test-task/internal/domain"
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
)

type FlatUsecaseTest struct {
	suite.Suite
	flatRepoMock *mock_domain.MockFlatRepo
	mockLg       *zap.Logger
}

func (f *FlatUsecaseTest) BeforeAll(t provider.T) {
	t.Log("Init mock")
	ctrl := gomock.NewController(t)
	f.flatRepoMock = mock_domain.NewMockFlatRepo(ctrl)
	f.mockLg = pkg.CreateMockLogger()
}

func (f *FlatUsecaseTest) TestNormalCreateFlat(t provider.T) {
	userUsecase := usecase.NewFlatUsecase(f.flatRepoMock)
	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db28")
	flat := domain.Flat{
		ID:      1000,
		UserID:  userID,
		HouseID: 4,
		Price:   1000,
		Rooms:   2,
		Status:  domain.CreatedStatus,
	}

	req := domain.CreateFlatRequest{
		FlatID:  1000,
		HouseID: 4,
		Price:   1000,
		Rooms:   2,
	}
	resp := domain.CreateFlatResponse{
		ID:      1000,
		HouseID: 4,
		Price:   1000,
		Rooms:   2,
		Status:  domain.CreatedStatus,
	}

	f.flatRepoMock.EXPECT().Create(gomock.Any(), &flat, f.mockLg).Return(flat, nil)

	created, err := userUsecase.Create(context.Background(), userID, &req, f.mockLg)

	t.Require().Nil(err)
	t.Require().Equal(resp, created)
}

func (f *FlatUsecaseTest) TestEmptyRequestCreateFlat(t provider.T) {
	userUsecase := usecase.NewFlatUsecase(f.flatRepoMock)
	created, err := userUsecase.Create(context.Background(), uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db28"), nil, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, created)
}

func (f *FlatUsecaseTest) TestBadFlatIDCreateFlat(t provider.T) {
	userUsecase := usecase.NewFlatUsecase(f.flatRepoMock)
	req := domain.CreateFlatRequest{
		FlatID:  0,
		HouseID: 4,
		Price:   1000,
		Rooms:   2,
	}

	created, err := userUsecase.Create(context.Background(), uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db28"), &req, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, created)
}

func (f *FlatUsecaseTest) TestBadHouseIDCreateFlat(t provider.T) {
	userUsecase := usecase.NewFlatUsecase(f.flatRepoMock)
	req := domain.CreateFlatRequest{
		FlatID:  10,
		HouseID: 0,
		Price:   1000,
		Rooms:   2,
	}

	created, err := userUsecase.Create(context.Background(), uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db28"), &req, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, created)
}

func (f *FlatUsecaseTest) TestBadRoomsCreateFlat(t provider.T) {
	userUsecase := usecase.NewFlatUsecase(f.flatRepoMock)
	req := domain.CreateFlatRequest{
		FlatID:  10,
		HouseID: 10,
		Price:   1000,
		Rooms:   0,
	}

	created, err := userUsecase.Create(context.Background(), uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db28"), &req, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, created)
}

func (f *FlatUsecaseTest) TestBadRepoCallCreateFlat(t provider.T) {
	userUsecase := usecase.NewFlatUsecase(f.flatRepoMock)
	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db28")
	flat := domain.Flat{
		ID:      1000,
		UserID:  userID,
		HouseID: 4,
		Price:   1000,
		Rooms:   2,
		Status:  domain.CreatedStatus,
	}

	req := domain.CreateFlatRequest{
		FlatID:  1000,
		HouseID: 4,
		Price:   1000,
		Rooms:   2,
	}

	f.flatRepoMock.EXPECT().Create(gomock.Any(), &flat, f.mockLg).Return(domain.Flat{}, errors.New("flat repo error"))

	created, err := userUsecase.Create(context.Background(), userID, &req, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, created)
}

func (f *FlatUsecaseTest) TestNormalUpdateFlat(t provider.T) {
	userUsecase := usecase.NewFlatUsecase(f.flatRepoMock)
	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db28")
	flat := domain.Flat{
		ID:      1000,
		HouseID: 4,
		Status:  domain.ModeratingStatus,
	}

	req := domain.UpdateFlatRequest{
		ID:      1000,
		HouseID: 4,
		Status:  domain.ModeratingStatus,
	}
	resp := domain.CreateFlatResponse{
		ID:      1000,
		HouseID: 4,
		Status:  domain.ModeratingStatus,
	}

	f.flatRepoMock.EXPECT().Update(gomock.Any(), userID, &flat, f.mockLg).Return(flat, nil)

	upd, err := userUsecase.Update(context.Background(), userID, &req, f.mockLg)

	t.Require().Nil(err)
	t.Require().Equal(resp, upd)
}

func (f *FlatUsecaseTest) TestBadStatusUpdateFlat(t provider.T) {
	userUsecase := usecase.NewFlatUsecase(f.flatRepoMock)
	userID := uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db28")

	req := domain.UpdateFlatRequest{
		ID:      1000,
		HouseID: 4,
		Status:  "bla bla",
	}

	upd, err := userUsecase.Update(context.Background(), userID, &req, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, upd)
}

func (f *FlatUsecaseTest) TestBadRepoCallUpdateFlat(t provider.T) {
	userUsecase := usecase.NewFlatUsecase(f.flatRepoMock)
	userID := uuid.Nil
	flat := domain.Flat{
		ID:          1000,
		UserID:      userID,
		HouseID:     4,
		Status:      domain.ApprovedStatus,
		ModeratorID: 0,
	}

	req := domain.UpdateFlatRequest{
		ID:      1000,
		HouseID: 4,
		Status:  domain.ApprovedStatus,
	}

	f.flatRepoMock.EXPECT().Update(gomock.Any(), userID, &flat, f.mockLg).Return(domain.Flat{}, errors.New("flat repo error"))

	created, err := userUsecase.Update(context.Background(), userID, &req, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, created)
}

func TestFlatUsecaseSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(FlatUsecaseTest))
}
