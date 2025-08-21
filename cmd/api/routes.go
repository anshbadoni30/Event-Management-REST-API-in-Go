package main

import ( 
	"net/http"
	"github.com/gin-gonic/gin"
)
func (app *application) routes() http.Handler{
	g:= gin.Default()
	v1:=g.Group("/api/v1")
	{
		//events
		v1.POST("/events",app.createEvent)  //Creating a event (require a user)
		v1.GET("/events",app.getAllEvents) //Print all events
		v1.GET("/events/:id",app.getEvent) //Print Sepcific Event
		v1.PUT("/events/:id",app.updateEvent) //Update an event by passing full updated event info
		v1.DELETE("/events/:id",app.deleteEvent) //delete an event
		//user
		v1.POST("/auth/register",app.registerUser) // Register a User and put in user table
		v1.POST("/user/:id",app.getUser) //Print info of User
		//attendees
		v1.POST("/events/:id/attendees/:userid",app.addAttendeeToEvent) //Add attendee in attendees table 
		v1.DELETE("/events/:id/attendees/:userid",app.deleteAtendeeFromEvent) // Delete an attendee
		v1.GET("/attendees/:id/events",app.getEventsByAttendee) //Print all events associated with an attendee (taking user id)
		v1.GET("/events/:id/attendees",app.getAttendeesForEvent) //Print all attendees associated with an event
		v1.POST("/auth/login",app.login)
	}

	authGroup:=v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.POST("/events",app.createEvent)
		authGroup.PUT("/events/:id",app.updateEvent) 
		authGroup.DELETE("/events/:id",app.deleteEvent)
		authGroup.POST("/events/:id/attendees/:userid",app.addAttendeeToEvent)
		authGroup.DELETE("/events/:id/attendees/:userid",app.deleteAtendeeFromEvent)
	}
	return g
}