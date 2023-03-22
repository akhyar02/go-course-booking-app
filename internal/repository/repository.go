package repository

import (
	"time"

	"github.com/akhyar02/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() ([]models.User, error)
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(rr models.RoomRestriction) error
	SearchAvailibilityByDates(start time.Time, end time.Time) ([]models.Room, error)
	SearchAvailibilityByDatesAndRoomId(start time.Time, end time.Time, roomId int) (bool, error)
}
