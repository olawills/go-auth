package controllers

import (
	database "auth_api_with_Go/db"
	"auth_api_with_Go/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterUser(context *gin.Context) {
	var user model.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{"id": uuid.New(), "first_name": user.FirstName, "last_name": user.LastName, "email": user.Email, "username": user.Username})
}

func LoginUser(login *model.Login, context *gin.Context) (*model.User, error) {
	user := &model.User{}

	if err := database.Instance.Where("email = ?", login.Email).First(&user).Error; err != nil {
		return nil, err
	}
	if err := user.CheckPassword(login.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Incorrect Password"})
		context.Abort()
		return nil, err
	}
	// context.JSON(http.StatusCreated, gin.H{"id": uuid.New(), "first_name": user.FirstName, "last_name": user.LastName, "email": user.Email, "username": user.Username})
	return user, nil
}
