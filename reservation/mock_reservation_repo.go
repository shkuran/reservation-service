package reservation

import (
	"errors"
	"time"

)

type MockReservationRepo struct {
	reservation []Reservation
}

func NewMockReservationRepo(res []Reservation) *MockReservationRepo {
	return &MockReservationRepo{reservation: res}
}

func (r *MockReservationRepo) GetAll() ([]Reservation, error) {
	if len(r.reservation) == 0 {
		return nil, errors.New("simulated error fetching reservations")
	}
	return r.reservation, nil
}

func (r *MockReservationRepo) GetById(id int64) (Reservation, error) {
	for _, res := range r.reservation {
		if res.ID == id {
			return res, nil
		}
	}
	return Reservation{}, errors.New("simulated error fetching reservation by id")
}

func (r *MockReservationRepo) Save(res Reservation) error {
	res.ID = int64(len(r.reservation)) + 1
	r.reservation = append(r.reservation, res)
	return nil
}

func (r *MockReservationRepo) UpdateReturnDate(id int64) error {
	returnDate := time.Now()
	for i := range r.reservation {
		if r.reservation[i].ID == id {
			r.reservation[i].ReturnDate = &returnDate
			return nil
		}
	}
	return errors.New("simulated error updating ReturnDate")
}
