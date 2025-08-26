package main

import (
	"net/http"
	"strconv"

	"github.com/anshbadoni30/event-management-app/internal/database"
	"github.com/gin-gonic/gin"
)

// CreateEvent creates a new event
//
//	@Summary		Creates a new event
//	@Description	Creates a new event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		database.Event	true	"Event"
//	@Success		201		{object}	database.Event
//	@Router			/api/v1/events [post]
//	@Security		BearerAuth
func (app *application) createEvent(c *gin.Context) {
	var event database.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user:=app.GetUserFromContext(c)
	event.OwnerId=user.Id
	err := app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}
	c.JSON(http.StatusCreated, event)
}

// getEvents return all events
//
// @Summary Returns all events
// @Description Resturns all events
// @Tags Events
// @Accept json
// @Produce json
// @Success 200 {object} []database.Event
// @Router /api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrive events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEvent returns a single event
//
//	@Summary		Returns a single event
//	@Description	Returns a single event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	database.Event
//	@Router			/api/v1/events/{id} [get]
func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.Get(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	c.JSON(http.StatusCreated, event)
}

// UpdateEvent updates an existing event
//
//	@Summary		Updates an existing event
//	@Description	Updates an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			event	body		database.Event	true	"Event"
//	@Success		200	{object}	database.Event
//	@Router			/api/v1/events/{id} [put]
//	@Security		BearerAuth
func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user:=app.GetUserFromContext(c)
	existingevent, err := app.models.Events.Get(id)

	if existingevent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if user.Id!=existingevent.Id{
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	updatedEvent := &database.Event{}

	err = c.ShouldBindBodyWithJSON(updatedEvent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedEvent.Id = id
	errr := app.models.Events.Update(updatedEvent)
	if errr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}
	c.JSON(http.StatusCreated, updatedEvent)
}

// DeleteEvent deletes an existing event
//
//	@Summary		Deletes an existing event
//	@Description	Deletes an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		204
//	@Router			/api/v1/events/{id} [delete]
//	@Security		BearerAuth
func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Event ID"})
		return
	}

	user:=app.GetUserFromContext(c)
	existingEvent,err:=app.models.Events.Get(id)
	if existingEvent==nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Event"})
		return
	}

	if user.Id!=existingEvent.OwnerId{
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this event"})
		return
	}
	err = app.models.Events.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete a event"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"success": "OK"})
}


// AddAttendeeToEvent adds an attendee to an event
// @Summary		Adds an attendee to an event
// @Description	Adds an attendee to an event
// @Tags			attendees
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Event ID"
// @Param			userId	path		int	true	"User ID"
// @Success		201		{object}	database.Attendee
// @Router			/api/v1/events/{id}/attendees/{userId} [post]
// @Security		BearerAuth
func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Event ID"})
		return
	}

	userid, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	//Checking of getting event details from event id
	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	//Checking of getting user details from userid
	userToAdd, err := app.models.Users.Get(userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user details"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user:=app.GetUserFromContext(c)
	if user.Id!=event.OwnerId{
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add an attendee"})
		return
	}

	//checking of user in attendee table
	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Attendee already exist"})
		return
	}

	//Insertion of user in attendees table
	attendee := database.Attendee{
		UserId:  userToAdd.Id,
		EventId: event.Id,
	}
	attendeeResult, err := app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add Attendee"})
		return
	}
	c.JSON(http.StatusCreated, attendeeResult)
}

// GetAttendeesForEvent returns all attendees for a given event
//
//	@Summary		Returns all attendees for a given event
//	@Description	Returns all attendees for a given event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	[]database.User
//	@Router			/api/v1/events/{id}/attendees [get]
func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) //Taking eventid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ivalid Event Id"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees for event"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// DeleteAttendeeFromEvent deletes an attendee from an event
// @Summary		Deletes an attendee from an event
// @Description	Deletes an attendee from an event
// @Tags			attendees
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Event ID"
// @Param			userId	path		int	true	"User ID"
// @Success		204
// @Router			/api/v1/events/{id}/attendees/{userId} [delete]
// @Security		BearerAuth
func (app *application) deleteAtendeeFromEvent(c *gin.Context) {
	eventid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}
	userid, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	event,err:=app.models.Events.Get(eventid)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Something went wrong"})
		return
	}
	if event==nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Event not found"})
		return
	}
	user:=app.GetUserFromContext(c)
	if user.Id!=event.OwnerId{
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete an attendee"})
		return
	}
	err = app.models.Attendees.Delete(userid, eventid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Attendee"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"success": "OK"})
}

// GetEventsByAttendee returns all events for a given attendee
//
//	@Summary		Returns all events for a given attendee
//	@Description	Returns all events for a given attendee
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Attendee ID"
//	@Success		200	{object}	[]database.Event
//	@Router			/api/v1/attendees/{id}/events [get]
func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Attendee Id"})
		return
	}

	events, err := app.models.Attendees.GetByAttendee(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees"})
		return
	}

	c.JSON(http.StatusOK, events)
}
