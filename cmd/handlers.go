package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Chained/auth-service/cmd/utils"
	pb "github.com/Chained/auth-service/github.com/Chained/auth-service"
	"github.com/Chained/auth-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type userInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser handler parses data into user variable
// and passes it to the gRPC server
func (app *Application) CreateUser(c *gin.Context) {
	var u *models.User
	var err error

	if err = c.ShouldBindJSON(&u); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := &pb.User{
		FirstName: u.Firstname,
		LastName:  u.Lastname,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := app.Client.
		Authorize(
			ctx,
			user,
		)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"user": res,
	})
}

// GetUser returns user
func (app *Application) GetUser(c *gin.Context) {
	var getUserRequest *pb.GetUserRequest
	var err error
	id := c.Param("id")

	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	getUserRequest.Id = int64(intId)

	res, err := app.Client.GetUser(context.Background(), getUserRequest)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"res": res.User,
	})
}

// EditUser updates information about user
func (app *Application) EditUser(c *gin.Context) {
	var updatedUser *pb.UpdateUserRequest
	var err error

	strId := c.Param("id")

	id, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	updatedUser.Id = int64(id)
	if err = c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = app.Client.UpdateUser(context.Background(), updatedUser)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Your info was successfully changed.",
	})
}

// DeleteUser deletes user from database
func (app *Application) DeleteUser(c *gin.Context) {
	var pdu *pb.DeleteUserRequest
	var err error

	err = jsonpb.Unmarshal(c.Request.Body, pdu)
	if err != nil {
		c.JSON(http.StatusTeapot, gin.H{
			"error":   true,
			"message": fmt.Sprintf("Error unmarshalling json to protobuf message"),
		})
	}

	err = app.DB.DeleteUser(int(pdu.Id))
	if err != nil {
		c.JSON(http.StatusTeapot, gin.H{
			"error":   true,
			"message": fmt.Sprintf("Error deleting user from database. Please, try again."),
		})
		return
	}

}

// CreateAuthenticationToken creates access token for user
func (app *Application) CreateAuthenticationToken(c *gin.Context) {
	var userInp userInput

	if err := c.ShouldBindJSON(&userInp); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := app.DB.GetUserByEmail(userInp.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	ok, err := utils.PasswordMatches(user.Password, userInp.Password)
	if err != nil {
		utils.InvalidCredentials(c)
		return
	}

	if !ok {
		utils.InvalidCredentials(c)
		return
	}

	token, err := models.GenerateToken(user.Id, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"error":                false,
		"message":              fmt.Sprintf("Token for %s successfully created", user.Email),
		"authentication_token": token,
	})
}

// CheckAuthentication verifying authentication status
func (app *Application) CheckAuthentication(c *gin.Context) {
	// validate the token
	user, err := app.authenticateToken(c.Request)
	if err != nil {
		utils.InvalidCredentials(c)
		return
	}

	c.JSON(500, gin.H{
		"error":   false,
		"message": fmt.Sprintf("Authenticated user: %v\n", user.Email),
	})
}

// authenticateToken verifies if access token from header is valid
func (app *Application) authenticateToken(r *http.Request) (*pb.User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("received wrong authorization header")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("authorization wrong size")
	}

	user, err := app.DB.GetUserByAccessToken(token)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}
