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
	OwnerId     int    `json:"owner_id" binding:"required"`
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02"`
	Location     string `json:"location" binding:"required,min=3"`
}

func (e *EventModel) Insert(event *Event) error {
	ctx, cancel:= context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	query:= "INSERT INTO events (owner_id, name, description, date, location) VALUES ($1,$2,$3,$4,$5)"

	return e.db.QueryRowContext(ctx,query,event.OwnerId,event.Name,event.Description,event.Date,event.Location).Scan(&event.Id)
}


func (m *EventModel) GetAll() ([]Event,error){
	ctx,cancel:= context.WithTimeout(context.Background(),3*time.Second)

	defer cancel()

	query:="SELECT * from events"

	rows,err:= m.db.QueryContext(ctx,query)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	var events []Event
	
	for rows.Next(){
		var event Event
		err:=rows.Scan(&event.Id,&event.OwnerId,&event.Name,&event.Description,&event.Date,&event.Location)
		if err!=nil{
			return nil,err
		}
		events=append(events, event)
	}
	if err=rows.Err();err!=nil{
		return nil,err
	}
	return events,nil
}

func (m* EventModel) Get(id int) (*Event,error){
	ctx,cancel:= context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	query:="Select * from events where id=$1"
	var event Event

	err:=m.db.QueryRowContext(ctx,query,id).Scan(&event.Id,&event.OwnerId,&event.Name,&event.Description,&event.Date,&event.Location)
	if err!=nil{
		if err==sql.ErrNoRows{
			return nil,fmt.Errorf("no student found with id %s",fmt.Sprint(id))
		}
		return nil,err
	}
	return nil,nil

}

func (e *EventModel) Update(event *Event) error{
	ctx,cancel:= context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	query:= "update events set name=$1,description=$2,date=$3, location=$4 where id=$5"
	_,err:=e.db.ExecContext(ctx,query,event.Name,event.Description,event.Date,event.Location)
	if err!=nil{
		return err
	}
	return nil
}

func (e *EventModel) Delete(id int) error{
	ctx,cancel:= context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()
	query:= "Delete from events where id=$1"
	_,err:=e.db.ExecContext(ctx,query,id)
	if err!=nil{
		return err
	}
	return nil
}