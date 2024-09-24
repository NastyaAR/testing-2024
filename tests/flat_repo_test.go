package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/repo"
	"avito-test-task/pkg"
	mock_domain "avito-test-task/tests/mocks"
	"context"
	"errors"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"math/rand"
	"testing"
	"time"
)

type FlatRepoTest struct {
	suite.Suite
	mockLg *zap.Logger
}

func (f *FlatRepoTest) BeforeAll(t provider.T) {
	t.Log("Init log")
	f.mockLg = pkg.CreateMockLogger()
}

func (f *FlatRepoTest) TestNormalCreateFlat(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	txMock := mock_domain.NewMockTx(ctrl)
	rowMock := mock_domain.NewMockRow(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(poolMock, retryAdapter)

	userID, _ := uuid.Parse("019126ee-2b7d-758e-bb22-fe2e45b2db22")
	flat := domain.Flat{
		ID:          rand.Intn(100) + 20,
		HouseID:     3,
		UserID:      userID,
		Price:       10000000,
		Rooms:       2,
		Status:      domain.CreatedStatus,
		ModeratorID: 1,
	}

	query := `insert into flats(flat_id, house_id, user_id, price, rooms, status)
			values ($1, $2, $3, $4, $5, $6) 
			returning flat_id, house_id, user_id, price, rooms, status`
	execQuery := `update houses set update_flat_date=$1 where house_id=$2`

	poolMock.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(txMock, nil)
	txMock.EXPECT().QueryRow(context.Background(), query, flat.ID, flat.HouseID, flat.UserID,
		flat.Price, flat.Rooms, flat.Status).Return(rowMock)
	rowMock.EXPECT().Scan(gomock.Any()).Return(nil)
	txMock.EXPECT().Exec(context.Background(), execQuery, gomock.Any(), gomock.Any())
	txMock.EXPECT().Commit(context.Background())

	flat, err := flatRepo.Create(context.Background(), &flat, f.mockLg)

	t.Require().Nil(err)
}

func (f *FlatRepoTest) TestContextTimeoutCreateFlat(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	txMock := mock_domain.NewMockTx(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 1, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(poolMock, retryAdapter)

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
	poolMock.EXPECT().BeginTx(ctx, pgx.TxOptions{}).Return(txMock, errors.New("expired context"))

	flat, err := flatRepo.Create(ctx, &flat, f.mockLg)

	t.Require().Error(err)
}

func (f *FlatRepoTest) TestNormalDeleteByIdFlat(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)

	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(poolMock, retryAdapter)

	userID, _ := uuid.Parse("019126ee-2b7d-758e-bb22-fe2e45b2db22")
	flat := domain.Flat{
		ID:          190,
		HouseID:     3,
		UserID:      userID,
		Price:       10000000,
		Rooms:       2,
		Status:      domain.CreatedStatus,
		ModeratorID: 1,
	}
	query := `delete from flats where flat_id=$1 and house_id=$2`
	cTag := pgconn.CommandTag{}

	poolMock.EXPECT().Exec(context.Background(), query, flat.ID, flat.HouseID).Return(cTag, nil)

	err := flatRepo.DeleteByID(context.Background(), flat.ID, flat.HouseID, f.mockLg)

	t.Require().Nil(err)
}

func (f *FlatRepoTest) TestContextTimeoutDeleteByIdFlat(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)

	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(poolMock, retryAdapter)
	query := `delete from flats where flat_id=$1 and house_id=$2`

	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
	time.Sleep(time.Millisecond)
	poolMock.EXPECT().Exec(ctx, query, 1, 1).Return(pgconn.CommandTag{}, errors.New("expired context"))

	err := flatRepo.DeleteByID(ctx, 1, 1, f.mockLg)

	t.Require().Error(err)
}

func (f *FlatRepoTest) TestNormalGetByIdFlat(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowsMock := mock_domain.NewMockRows(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(poolMock, retryAdapter)

	query := `select flat_id, house_id, user_id, price, rooms, status
	from flats where flat_id=$1 and house_id=$2`

	poolMock.EXPECT().QueryRow(context.Background(), query, 1, 1).Return(rowsMock)
	rowsMock.EXPECT().Scan(gomock.Any()).Return(nil)

	_, err := flatRepo.GetByID(context.Background(), 1, 1, f.mockLg)

	t.Require().Nil(err)
}

func (f *FlatRepoTest) TestNormalGetAllFlat(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowsMock := mock_domain.NewMockRows(ctrl)

	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(poolMock, retryAdapter)
	query := `select flat_id, house_id, user_id, price, rooms, status from flats limit $1 offset $2`
	poolMock.EXPECT().Query(context.Background(), query, 10, 0).Return(rowsMock, nil)
	rowsMock.EXPECT().Next()
	rowsMock.EXPECT().Close()

	_, err := flatRepo.GetAll(context.Background(), 0, 10, f.mockLg)

	t.Require().Nil(err)
}

func (f *FlatRepoTest) TestContextTimeoutGetAllFlat(t provider.T) {
	ctrl := gomock.NewController(t)
	poolMock := mock_domain.NewMockIPool(ctrl)
	rowsMock := mock_domain.NewMockRows(ctrl)
	retryAdapter := repo.NewPostgresRetryAdapter(poolMock, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(poolMock, retryAdapter)
	ctx, _ := context.WithTimeout(context.Background(), 1)
	time.Sleep(1)
	query := `select flat_id, house_id, user_id, price, rooms, status from flats limit $1 offset $2`

	poolMock.EXPECT().Query(ctx, query, 1, 1).Return(rowsMock, errors.New("expired context"))
	rowsMock.EXPECT().Close().MinTimes(1)

	flats, err := flatRepo.GetAll(ctx, 1, 1, f.mockLg)

	t.Require().Error(err)
	t.Require().Equal(len(flats), 0)
}

func TestFlatSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(FlatRepoTest))
}
