package database

import "database/sql"

type AttendeeModel struct {
	db *sql.DB
}

type Attendee struct{
	Id int `json:"id"`
	UserId int `json:"userId"`
	EventId int `json:"eventId"`
}