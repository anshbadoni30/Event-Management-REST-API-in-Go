package main

import (
	"net/http"
	"strconv"

	"github.com/anshbadoni30/event-management-app/internal/database"
	"github.com/gin-gonic/gin"
)

func (app *application) createEvent(c *gin.Context){
	var event database.Event
	if err:=c.ShouldBindJSON(&event);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	err:=app.models.Events.Insert(&event)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to create event"})
		return
	}
	c.JSON(http.StatusCreated,event)
}

func (app *application)getAllEvents(c *gin.Context){
	events,err:=app.models.Events.GetAll()
	
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"failed to retrive events"})
		return
	}

	c.JSON(http.StatusOK,events)
}

func (app *application) getEvent(c *gin.Context){
	id, err:= strconv.Atoi(c.Param("id"))

	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid event ID"})
		return
	}

	event,err:=app.models.Events.Get(id)

	if event==nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Event not found"})
		return
	}

	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to retrieve event"})
		return
	}
	c.JSON(http.StatusCreated,event)
}

func(app *application)updateEvent(c *gin.Context){
	id, err:= strconv.Atoi(c.Param("id"))

	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid event ID"})
		return
	}
	existingevent,err:=app.models.Events.Get(id)

	if existingevent==nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Event not found"})
		return
	}

	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to retrieve event"})
		return
	}

	updatedEvent:= &database.Event{}

	err=c.ShouldBindBodyWithJSON(updatedEvent)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	updatedEvent.Id=id
	errr:=app.models.Events.Update(updatedEvent)
	if errr!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to update event"})
		return
	}
	c.JSON(http.StatusCreated,updatedEvent)
}

func (app *application)deleteEvent(c *gin.Context){
	id, err:=strconv.Atoi(c.Param("id"))
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid Event ID"})
		return
	}

	err=app.models.Events.Delete(id)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to delete a event"})
		return
	}
	c.JSON(http.StatusNoContent,gin.H{"success":"OK"})
}
