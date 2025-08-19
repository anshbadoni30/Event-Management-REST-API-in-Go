package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeeModel struct {
	db *sql.DB
}

type Attendee struct {
	Id      int `json:"id"`
	UserId  int `json:"userId"`
	EventId int `json:"eventId"`
}

func (m *AttendeeModel) Insert(attendee *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "insert into attendees (event_id,user_id) values (?,?)"
	result, err := m.db.ExecContext(ctx, query, attendee.EventId, attendee.UserId)

	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	attendee.Id = int(id)
	return attendee, nil
}

func (m *AttendeeModel) GetByEventAndAttendee(eventid, userid int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select * from attendees where event_id= ?  and user_id = ?"

	var attendee Attendee
	err := m.db.QueryRowContext(ctx, query, eventid, userid).Scan(&attendee.Id, &attendee.UserId, &attendee.EventId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &attendee, nil
}

func (m *AttendeeModel) GetAttendeesByEvent(eventid int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "select u.id,u.name,u.email from users u JOIN attendees a ON u.id=a.user_id where a.event_id=? "

	rows, err := m.db.QueryContext(ctx, query, eventid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (m *AttendeeModel) Delete(userId, eventId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "delete from attendees where user_id = ? and event_id = ?"
	_, err := m.db.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return err
	}
	return nil
}

func (m *AttendeeModel) GetByAttendee(attendeeid int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "Select e.id,e.owner_id,e.name,e.description,e.date,e.location from events e JOIN attendees a on e.id=a.event_id where a.user_id=?"

	rows, err := m.db.QueryContext(ctx, query, attendeeid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}
