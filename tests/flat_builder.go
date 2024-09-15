package tests

import (
	"avito-test-task/internal/domain"
	"github.com/google/uuid"
)

type FlatMother struct{}

func (f *FlatMother) DefaultFlat(id int, houseID int) domain.Flat {
	return domain.Flat{
		ID:          id,
		HouseID:     houseID,
		UserID:      uuid.New(),
		Price:       10000,
		Rooms:       2,
		Status:      domain.CreatedStatus,
		ModeratorID: 1,
	}
}

func (f *FlatMother) DefaultFlatResponse(flat *domain.Flat) domain.CreateFlatResponse {
	return domain.CreateFlatResponse{
		ID:      flat.ID,
		HouseID: flat.HouseID,
		Price:   flat.Price,
		Rooms:   flat.Rooms,
		Status:  flat.Status,
	}
}

func (f *FlatMother) DefaultSingleFlatResponse(flat *domain.Flat) domain.SingleFlatResponse {
	return domain.SingleFlatResponse{
		ID:      flat.ID,
		HouseID: flat.HouseID,
		Price:   flat.Price,
		Rooms:   flat.Rooms,
		Status:  flat.Status,
	}
}
