package handlers

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

type RegisterRequest struct {
	Fullname string                `json:"fullname" binding:"required"`
	Username string                `json:"username" binding:"required"`
	Email    string                `json:"email" binding:"required,email"`
	Password string                `json:"password" binding:"required,min=6"`
	Avatar   *multipart.FileHeader `form:"avatar"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterHeaders struct {
	ClientVersion string `header:"X-App-Version"`
	RequestID     string `header:"X-Request-Id"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	var headers RegisterHeaders

	if err := c.ShouldBindHeader(&headers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// directly access from gin:
	// file, err := c.FormFile("avatar")
	// if err == nil {
	// 	c.SaveUploadedFile(file, "./uploads/"+file.Filename)
	// }

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Avatar != nil {
		err := c.SaveUploadedFile(req.Avatar, "./uploads/"+req.Avatar.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to save uploaded file",
			})
			return
		}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.User{
		Fullname: req.Fullname,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword, // hash password
	}
	userData, err := repository.Register(&user)

	if err != nil {
		if msg, ok := utils.ParsePostgresError(err); ok {
			utils.ResponseError(c, http.StatusConflict, msg, nil)
			return
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		return
	}
	userResponse := utils.ToUserResponse(userData)
	utils.ResponseSuccess(c, http.StatusOK, "user registered successfully", userResponse)

}

func Login(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "error in request bpdy", nil)
	}

	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid Credential", nil)
		return
	}

	if isPasswordValid := utils.CheckPassword(req.Password, user.Password); isPasswordValid == false {
		utils.ResponseError(c, http.StatusInternalServerError, "Invalid Credential", nil)
		return
	}
	userResponse := utils.ToUserResponse(user)
	utils.ResponseSuccess(c, http.StatusOK, "loggedin successfully", userResponse)
}

func GetAllUsers(c *gin.Context) {
	users, err := repository.GetAllUsers()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id") // returns string
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := repository.GetUserByUUID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err, "status": "Failure", "message": "Something went wrong"})
		return
	}
	userResponse := utils.ToUserResponse(user)
	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", userResponse)
}

func GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err, "status": "Failure", "message": "Something went wrong"})
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", user)
}
