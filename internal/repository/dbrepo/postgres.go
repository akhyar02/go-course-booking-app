package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/akhyar02/bookings/internal/models"
)

func (m *postgresDbRepo) AllUsers() ([]models.User, error) {
	query := `SELECT id, first_name, last_name, email, created_at, updated_at FROM users`
	users := []models.User{}
	rows, err := m.DB.Query(query)
	if err != nil {
		log.Println("Error getting all users:", err)
		return nil, err
	}

	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			err = rows.Close()
			if err != nil {
				log.Println("Error closing rows:", err)
				return users, err

			}
			return users, err
		}
		users = append(users, user)
	}
	err = rows.Close()
	if err != nil {
		log.Println("Error closing rows:", err)
		return users, err
	}
	return users, nil
}

func (m *postgresDbRepo) InsertReservation(res models.Reservation) (int, error) {
	var newId int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	row := m.DB.QueryRowContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomId, time.Now(), time.Now())
	err := row.Scan(&newId)
	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (m *postgresDbRepo) InsertRoomRestriction(rr models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := m.DB.ExecContext(ctx, stmt, rr.StartDate, rr.EndDate, rr.RoomId, rr.ReservationId, rr.RestrictionId, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDbRepo) SearchAvailibilityByDatesAndRoomId(start time.Time, end time.Time, roomId int) (bool, error) {
	var numRows int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT COUNT(id) FROM room_restrictions
	WHERE room_id = $1
	AND ((start_date >= $2 AND start_date <= $3)
	OR (end_date >= $2 AND end_date <= $3)
	OR (start_date < $2 AND end_date > $3))`

	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows != 0 {
		return false, nil
	}
	return true, nil
}

func (m *postgresDbRepo) SearchAvailibilityByDates(start time.Time, end time.Time) ([]models.Room, error) {
	var availableRooms []models.Room
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT r.id, r.room_name FROM rooms WHERE r.id NOT IN (SELECT rr.room_id FROM room_restrictions rr WHERE 
		((rr.start_date >= $2 AND rr.start_date <= $3)
			OR (rr.end_date >= $2 AND rr.end_date <= $3)
			OR (rr.start_date < $2 AND rr.end_date > $3))
		)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		room := models.Room{}
		err := rows.Scan(&room.Id, &room.RoomName)
		if err != nil {
			return availableRooms, err
		}
		availableRooms = append(availableRooms, room)
	}

	return availableRooms, nil
}
