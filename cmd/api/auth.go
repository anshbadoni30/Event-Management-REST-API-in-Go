package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anshbadoni30/event-management-app/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

type loginRequest struct{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginResponse struct{
	Token string `json:"token"`
}

// Login logs in a user
//
//	@Summary		Logs in a user
//	@Description	Logs in a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body	loginRequest	true	"User"
//	@Success		200	{object}	loginResponse
//	@Router			/api/v1/auth/login [post]
func (app *application)login(c *gin.Context){
	var auth loginRequest
	if err:=c.ShouldBindJSON(&auth); err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
		return
	}
	//Email checking
	existingUser,err:=app.models.Users.GetByEmail(auth.Email)
	if existingUser==nil{
		c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid email or password"})
		return
	}
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Something went wrong"})
		return
	}

	//Password checking
	err=bcrypt.CompareHashAndPassword([]byte(existingUser.Password),[]byte(auth.Password))
	if err==nil{
		c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid email or password"})
		return
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"userId":existingUser.Id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString,err:=token.SignedString([]byte(app.jwtSecret))
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"error generating token"})
		return
	}
	c.JSON(http.StatusOK,loginResponse{Token: tokenString})
}

// RegisterUser registers a new user
// @Summary		Registers a new user
// @Description	Registers a new user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user	body		registerRequest	true	"User"
// @Success		201	{object}	database.User
// @Router			/api/v1/auth/register [post]
func (app *application) registerUser(c *gin.Context){
	var register registerRequest
	if err:=c.ShouldBindJSON(&register); err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
		return
	}

	hashedPassword,err:=bcrypt.GenerateFromPassword([]byte(register.Password),bcrypt.DefaultCost)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Something went wrong"})
		return
	}
	register.Password=string(hashedPassword)
	user:=database.User{
		Email: register.Email,
		Password: register.Password,
		Name: register.Name,
	}
	
	err=app.models.Users.Insert(&user)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not registered successfully"})
		return
	}
	c.JSON(http.StatusCreated,user)

}

func (app * application)getUser(c *gin.Context){
	id,err:=strconv.Atoi(c.Param("id"))
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid User ID"})
		return
	}

	user,err:=app.models.Users.Get(id)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to retrieve user details"})
		return
	}
	c.JSON(http.StatusOK,user)
}