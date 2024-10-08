package repo

import (
	"avito-test-task/internal/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type PostgresHouseRepo struct {
	db           IPool
	retryAdapter IPostgresRetryAdapter
}

func NewPostgresHouseRepo(db IPool, retryAdapter IPostgresRetryAdapter) *PostgresHouseRepo {
	return &PostgresHouseRepo{
		db:           db,
		retryAdapter: retryAdapter,
	}
}

func (p *PostgresHouseRepo) Create(ctx context.Context, house *domain.House, lg *zap.Logger) (domain.House, error) {
	lg.Info("create house", zap.Int("house_id", house.HouseID))

	var createdHouse domain.House
	query := `insert into houses(address, construct_year, developer, create_house_date, update_flat_date)
	values ($1, $2, $3, $4, $5) returning *`
	rows := p.db.QueryRow(ctx, query,
		house.Address, house.ConstructYear,
		house.Developer, house.CreateHouseDate,
		house.UpdateFlatDate)

	err := rows.Scan(&createdHouse.HouseID, &createdHouse.Address,
		&createdHouse.ConstructYear, &createdHouse.Developer,
		&createdHouse.CreateHouseDate, &createdHouse.UpdateFlatDate)
	if err != nil {
		lg.Warn("postgres house create error", zap.Error(err))
		return domain.House{}, err
	}

	return createdHouse, nil
}

func (p *PostgresHouseRepo) DeleteByID(ctx context.Context, id int, lg *zap.Logger) error {
	lg.Info("delete house", zap.Int("house_id", id))

	query := `delete from houses where house_id=$1`
	_, err := p.db.Exec(ctx, query, id)
	if err != nil {
		lg.Warn("postgres house delete error", zap.Error(err))
		return err
	}

	return nil
}

func (p *PostgresHouseRepo) Update(ctx context.Context, newHouseData *domain.House, lg *zap.Logger) error {
	lg.Info("update house", zap.Int("house_id", newHouseData.HouseID))

	query := `update houses set house_id=$1,
                  				address=$2,
                  				construct_year=$3,
                  				developer=$4,
                  				create_house_date=$5,
                  				update_flat_date=$6
                  				where house_id=$1`
	_, err := p.db.Exec(ctx, query, newHouseData.HouseID, newHouseData.Address,
		newHouseData.ConstructYear, newHouseData.Developer,
		newHouseData.CreateHouseDate, newHouseData.UpdateFlatDate)
	if err != nil {
		lg.Warn("postgres house update error", zap.Error(err))
		return err
	}

	return nil
}

func (p *PostgresHouseRepo) GetByID(ctx context.Context, id int, lg *zap.Logger) (domain.House, error) {
	lg.Info("get house by id", zap.Int("id", id))
	var house domain.House

	query := `select * from houses where house_id=$1`
	rows := p.db.QueryRow(ctx, query, id)

	err := rows.Scan(&house.HouseID, &house.Address,
		&house.ConstructYear, &house.Developer,
		&house.CreateHouseDate, &house.UpdateFlatDate)
	if err != nil {
		lg.Warn("postgres house get by id error", zap.Error(err))
		return domain.House{}, err
	}

	return house, nil
}

func (p *PostgresHouseRepo) GetAll(ctx context.Context, offset int, limit int, lg *zap.Logger) ([]domain.House, error) {
	lg.Info("get houses", zap.Int("offset", offset), zap.Int("limit", limit))

	query := `select * from houses order by house_id limit $1 offset $2`
	rows, err := p.db.Query(ctx, query, limit, offset)
	defer rows.Close()
	if err != nil {
		lg.Warn("postgres house get all error", zap.Error(err))
		return nil, err
	}

	var (
		houses []domain.House
		house  domain.House
	)
	for rows.Next() {
		err = rows.Scan(&house.HouseID, &house.Address, &house.ConstructYear,
			&house.Developer, &house.CreateHouseDate, &house.UpdateFlatDate)
		if err != nil {
			lg.Warn("postgres house get all error: scan house error")
			continue
		}
		houses = append(houses, house)
	}

	return houses, err
}

func (p *PostgresHouseRepo) GetFlatsByHouseID(ctx context.Context, id int, status string, lg *zap.Logger) ([]domain.Flat, error) {
	lg.Info("get flats by house id", zap.Int("house_id", id))

	query := ``
	var (
		rows pgx.Rows
		err  error
	)
	if status == domain.ModeratingStatus {
		query = `select flat_id, houses.house_id, user_id, price, rooms, status 
			from flats join houses
			on flats.house_id = houses.house_id
			where houses.house_id=$1 and flats.status=$2`
		rows, err = p.db.Query(ctx, query, id, status)
	} else {
		query = `select flat_id, houses.house_id, user_id, price, rooms, status 
			from flats join houses
			on flats.house_id = houses.house_id
			where houses.house_id=$1`
		rows, err = p.db.Query(ctx, query, id)
	}

	defer rows.Close()
	if err != nil {
		lg.Warn("postgres house repo: get flats by house id", zap.Error(err))
		return nil, fmt.Errorf("postgres house repo: get flats by house id: %v", err.Error())
	}

	var flats []domain.Flat
	for rows.Next() {
		flat := domain.Flat{}
		err = rows.Scan(&flat.ID, &flat.HouseID, &flat.UserID, &flat.Price, &flat.Rooms, &flat.Status)
		if err != nil {
			lg.Warn("postgres house repo: get all error: scan house error", zap.Error(err))
			continue
		}
		flats = append(flats, flat)
	}

	return flats, err
}

func (p *PostgresHouseRepo) SubscribeByID(ctx context.Context, houseID int, userID uuid.UUID, lg *zap.Logger) error {
	lg.Info("postgres house repo: subscribe by id")

	query := `insert into subscribers(user_id, house_id) values ($1, $2)`
	_, err := p.db.Exec(ctx, query, userID, houseID)
	if err != nil {
		lg.Warn("postgres house repo: subscribe by id error", zap.Error(err))
		return fmt.Errorf("postgres house repo: subscribe by id error: %v", err.Error())
	}

	return nil
}
