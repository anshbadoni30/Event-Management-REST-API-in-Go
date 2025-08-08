package database

import "database/sql"

type Models struct {
	Users     UserModel
	Events    EventModel
	Attendees AttendeeModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:     UserModel{db: db},
		Events:    EventModel{db: db},
		Attendees: AttendeeModel{db: db},
	}
}
