package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // hash later
	}

	if err := repository.Register(&user); err != nil {
		fmt.Println("registration error occured")
		if msg, ok := utils.ParsePostgresError(err); ok {
			utils.Error(c, http.StatusConflict, msg, nil)
			return
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "user registered successfully",
		"data":    user,
	})
}

func Login(c *gin.Context) {

}

func GetAllUsers(c *gin.Context) {

}
