package integration

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/repo"
	"avito-test-task/internal/usecase"
	"avito-test-task/pkg"
	"avito-test-task/tests"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/zap"
	"testing"
)

type FlatIntegrationTest struct {
	suite.Suite
	flatUsecase domain.FlatUsecase
	flatRepo    domain.FlatRepo
	db          repo.IPool
	mockLg      *zap.Logger
	flatMother  *tests.FlatMother
}

func (f *FlatIntegrationTest) BeforeAll(t provider.T) {
	connString := "postgres://test-user:test-password@127.0.0.1:5431/test-db?sslmode=disable"

	var err error
	f.db, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf("error while connecting to db: %v", err.Error())
	}
	f.flatRepo = repo.NewPostgresFlatRepo(f.db, nil)
	f.flatUsecase = usecase.NewFlatUsecase(f.flatRepo)
	f.mockLg = pkg.CreateMockLogger()
	f.flatMother = &tests.FlatMother{}
}

func (f *FlatIntegrationTest) AfterAll(t provider.T) {
	f.db.Close()
}

func (f *FlatIntegrationTest) TestNormalCreate(t provider.T) {
	newFlat := domain.CreateFlatRequest{
		FlatID:  11,
		HouseID: 1,
		Price:   1000000,
		Rooms:   2,
	}
	expected := f.flatMother.DefaultFlatResponseFromRequest(&newFlat)

	created, err := f.flatUsecase.Create(context.Background(),
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db22"), &newFlat, f.mockLg)

	t.Require().Nil(err)
	t.Require().Equal(expected, created)
}

func (f *FlatIntegrationTest) TestNoExistHouseFlatCreate(t provider.T) {
	newFlat := domain.CreateFlatRequest{
		FlatID:  11,
		HouseID: 100,
		Price:   1000000,
		Rooms:   2,
	}

	created, err := f.flatUsecase.Create(context.Background(),
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db22"), &newFlat, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, created)
}

func (f *FlatIntegrationTest) TestBadPriceCreate(t provider.T) {
	newFlat := domain.CreateFlatRequest{
		FlatID:  11,
		HouseID: 100,
		Price:   -1000000,
		Rooms:   2,
	}

	created, err := f.flatUsecase.Create(context.Background(),
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db22"), &newFlat, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, created)
}

func (f *FlatIntegrationTest) TestBadRoomsCreate(t provider.T) {
	newFlat := domain.CreateFlatRequest{
		FlatID:  11,
		HouseID: 100,
		Price:   1000000,
		Rooms:   -2,
	}

	created, err := f.flatUsecase.Create(context.Background(),
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db22"), &newFlat, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, created)
}

func (f *FlatIntegrationTest) TestNormalUpdate(t provider.T) {
	updFlat := domain.UpdateFlatRequest{
		ID:      1,
		HouseID: 1,
		Status:  domain.ApprovedStatus,
	}

	expected := domain.CreateFlatResponse{
		ID:      1,
		HouseID: 1,
		Price:   100,
		Rooms:   2,
		Status:  domain.ApprovedStatus,
	}

	updated, err := f.flatUsecase.Update(context.Background(),
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db23"),
		&updFlat, f.mockLg)

	t.Require().Nil(err)
	t.Require().Equal(expected, updated)
}

func (f *FlatIntegrationTest) TestNilRequestUpdate(t provider.T) {
	updated, err := f.flatUsecase.Update(context.Background(),
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db23"),
		nil, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, updated)
}

func (f *FlatIntegrationTest) TestBadIDUpdate(t provider.T) {
	updFlat := domain.UpdateFlatRequest{
		ID:      -1,
		HouseID: 1,
		Status:  domain.ApprovedStatus,
	}

	updated, err := f.flatUsecase.Update(context.Background(),
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db23"),
		&updFlat, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, updated)
}

func (f *FlatIntegrationTest) TestBadHouseIDUpdate(t provider.T) {
	updFlat := domain.UpdateFlatRequest{
		ID:      1,
		HouseID: -1,
		Status:  domain.ApprovedStatus,
	}

	updated, err := f.flatUsecase.Update(context.Background(),
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db23"),
		&updFlat, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, updated)
}

func (f *FlatIntegrationTest) TestBadStatusIDUpdate(t provider.T) {
	updFlat := domain.UpdateFlatRequest{
		ID:      1,
		HouseID: -1,
		Status:  "status",
	}

	updated, err := f.flatUsecase.Update(context.Background(),
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db23"),
		&updFlat, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, updated)
}

func (f *FlatIntegrationTest) TestNoExistFlatUpdate(t provider.T) {
	updFlat := domain.UpdateFlatRequest{
		ID:      1,
		HouseID: 1000,
		Status:  "status",
	}

	updated, err := f.flatUsecase.Update(context.Background(),
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db23"),
		&updFlat, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, updated)
}

func (f *FlatIntegrationTest) TestNoExistModeratorFlatUpdate(t provider.T) {
	updFlat := domain.UpdateFlatRequest{
		ID:      1,
		HouseID: 1000,
		Status:  "status",
	}

	updated, err := f.flatUsecase.Update(context.Background(),
		uuid.MustParse("019126ee-2b7d-758e-bb22-fe2e45b2db60"),
		&updFlat, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(domain.CreateFlatResponse{}, updated)
}

func TestFlatSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(FlatIntegrationTest))
}
