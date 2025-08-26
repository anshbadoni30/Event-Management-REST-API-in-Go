package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()
	v1 := g.Group("/api/v1")
	{
		//events
		v1.GET("/events", app.getAllEvents) //Print all events
		v1.GET("/events/:id", app.getEvent) //Print Sepcific Event
		//user
		v1.POST("/auth/register", app.registerUser) // Register a User and put in user table
		v1.POST("/auth/login", app.login)           // Login user
		v1.POST("/user/:id", app.getUser)           //Print info of User
		//attendees
		v1.GET("/attendees/:id/events", app.getEventsByAttendee)  //Print all events associated with an attendee (taking user id)
		v1.GET("/events/:id/attendees", app.getAttendeesForEvent) //Print all attendees associated with an event
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.POST("/events", app.createEvent)                                    //Creating a event (require a user)
		authGroup.PUT("/events/:id", app.updateEvent)                                 //Update an event by passing full updated event info
		authGroup.DELETE("/events/:id", app.deleteEvent)                              //delete an event
		authGroup.POST("/events/:id/attendees/:userid", app.addAttendeeToEvent)       //Add attendee in attendees table
		authGroup.DELETE("/events/:id/attendees/:userid", app.deleteAtendeeFromEvent) // Delete an attendee
	}

	g.GET("/swagger/*any",func(c *gin.Context){
		if c.Request.RequestURI=="/swagger/"{
			c.Redirect(302,"/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler,ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})
	return g
}
