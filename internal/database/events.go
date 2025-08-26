package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type EventModel struct {
	db *sql.DB
}

type Event struct {
	Id          int    `json:"id"`
	OwnerId     int    `json:"owner_id"`
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Date        string `json:"date" binding:"required"`
	Location    string `json:"location" binding:"required,min=3"`
}

func (e *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO events (owner_id, name, description, date, location) VALUES (?,?,?,?,?)"

	result, err := e.db.ExecContext(ctx, query, event.OwnerId, event.Name, event.Description, event.Date, event.Location)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	event.Id = int(id)
	return nil
}

func (m *EventModel) GetAll() ([]Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "SELECT * from events"

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (m *EventModel) Get(id int) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "Select * from events where id=?"
	var event Event

	err := m.db.QueryRowContext(ctx, query, id).Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no event found with id %s", fmt.Sprint(id))
		}
		return nil, err
	}
	return &event, nil

}

func (e *EventModel) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "update events set name=$1,description=$2,date=$3, location=$4 where id=$5"
	_, err := e.db.ExecContext(ctx, query, event.Name, event.Description, event.Date, event.Location, event.Id)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "Delete from events where id=$1"
	_, err := e.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
