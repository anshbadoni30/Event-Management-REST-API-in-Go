package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (app *application) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(app.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("what")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userId := claims["userId"].(float64)
		user, err := app.models.Users.Get(int(userId))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

/*
Explaination of above code:
1 This code is used to protect specific routes which is used by clients in order to protect the session handling from attackers
2 It uses JWT token authentication
3 The process starts when clients login, then in auth.go file in login fucntion, jwt token is created for a specific user.
4 After this, every time when client wants to access a specific route (made for user access), the session is started for client by using their unique jwt token which is created during login
5 Through this file, when we want to access a route as client we passed the token in authentication field
6 Firstly, the authentication will be checked if it has the prefix "Bearer" or not.
7 If it has that prefix, it will remove it from and store the jwt token in token string
8 Now here comes process of jwt.parse function:
	i) A JWT has 3 parts, separated by dots:
	Header → JSON describing algorithm, e.g. {"alg":"HS256","typ":"JWT"}
	Payload → JSON with claims (userId, exp, etc.)
	Signature → cryptographic proof that header+payload weren’t tampered with

	ii)Parse the header + payload
	jwt.Parse reads the first two parts (header, payload).
	From the header it sees "alg":"HS256", so it sets token.Method = SigningMethodHS256.
	Call your inner Func
	You check the algorithm really is HMAC.
	Then you return your server’s secret key (app.jwtSecret).

	iii)Verify the signature
	The library takes:
		The secret key you just returned.
		The decoded header+payload.
		The algorithm object (HS256).
		It re-computes:	signature = HMAC-SHA256(header.payload, secret)
	Then compares this computed signature against the 3rd field from the JWT string.

	iv)Mark token as valid/invalid
	If the signature matches → token.Valid = true.
	If it doesn’t match (wrong secret, tampered header/payload, or wrong alg) → invalid.

9 After this, claims variable retrive the json format paylod data through jwt.mapclaims
10 Then we will retrieve userid from json data and verifies it by using our database
*/
