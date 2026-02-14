package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/helper"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type RegisterRequest struct {
	Fullname string `json:"fullname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Gender   string `json:"gender" binding:"required,oneof=male female other"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	// Avatar   *multipart.FileHeader `form:"avatar"`
}

type LoginRequest struct {
	// Username string `json:"username" binding:"required"`
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
	defaultRoleId := config.DefaultUserRoleID
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

	// if req.Avatar != nil {
	// 	err := c.SaveUploadedFile(req.Avatar, "./uploads/"+req.Avatar.Filename)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"error": "failed to save uploaded file",
	// 		})
	// 		return
	// 	}
	// }

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.User{
		Fullname: req.Fullname,
		Username: req.Username,
		Email:    req.Email,
		Gender:   req.Gender,
		Password: hashedPassword, // hash password
		RoleId:   defaultRoleId,
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
	var req LoginRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "error in request body", nil)
		return
	}

	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "User not found", nil)
		return
	}

	isPasswordValid, err := utils.CheckPassword(user.Password, req.Password)
	if err != nil {
		fmt.Println("Error:", err)
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if !isPasswordValid {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid Credential", nil)
		return
	}

	claims := utils.JWTClaims{
		UserID: user.ID.String(),
		Roles:  []string{"user", "admin"}, // dynamic from DB
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "e-commerce-app",
		},
	}
	token, err := utils.CreateToken(claims)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "could not create token", nil)
		return
	}

	userResponse := utils.ToUserResponse(user)
	// utils.ResponseSuccess(c, http.StatusOK, "loggedin successfully", userResponse)

	// utils.ResponseSuccess(c, http.StatusOK, "loggedin successfully", gin.H{
	// 	"data":  userResponse,
	// 	"token": token,
	// })
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "logged in successfully",
		"data":    userResponse,
		"token":   token,
	})
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

func GetFilterAndSearchUsers(c *gin.Context) {
	// 1. Parse and set defaults for pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	// 2. Create the params object
	params := helper.UserFilterParams{
		ProductName: c.Query("productName"),
		FullName:    c.Query("fullname"),
		Page:        page,
		Limit:       limit,
	}
	users, total, err := repository.FilterAndSearchUsers(params)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}
	// 3. Return response with metadata
	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", gin.H{
		"users": users,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func SendOtpRequest(c *gin.Context) {
	var req helper.PasswordResetRequest
	var reset_password models.PasswordReset
	randomOtp := utils.GenerateOtp()
	otpHash, err := utils.HashOtp(randomOtp)

	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "User not found", nil)
		return
	}

	reset_password.UserID = user.ID.String()
	reset_password.OTPHash = otpHash
	reset_password.ExpiresAt = time.Now().Add(15 * time.Minute) // OTP valid for 15 minutes

	if err := repository.CreatePasswordReset(&reset_password); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create password reset request", nil)
		return
	}

	// Here you would send the OTP to the user's email using an email service.
	// For this example, we will just return the OTP in the response (DO NOT do this in production).
	utils.ResponseSuccess(c, http.StatusOK, "OTP sent successfully", gin.H{
		"otp": randomOtp, // In production, do not return OTP in response
	})
}

func VerifyOtpRequest(c *gin.Context) {
	// This handler would verify the OTP and allow the user to reset their password.
	// Implementation would involve checking the OTP against the hashed value in the database,
	// ensuring it hasn't expired, and then allowing the user to set a new password.
	// genereate reset token and return to user if OTP is valid

	var req helper.VerifyOtpRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "User not found", nil)
		return
	}

	resetRecord, err := repository.GetPasswordResetByUserID(user.ID.String())
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "No OTP request found for this user", nil)
		return
	}

	if resetRecord.AttemptCount >= 5 {
		utils.ResponseError(c, http.StatusTooManyRequests, "Too many wrong attempts. Try later.", nil)
		return
	}

	if time.Now().After(resetRecord.CreatedAt.Add(15 * time.Minute)) {
		utils.ResponseError(c, http.StatusBadRequest, "OTP expired", nil)
		return
	}

	isValidOtp, err := utils.CheckOtpHash(req.Otp, resetRecord.OTPHash)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Failed to verify OTP", nil)
		return
	}
	if !isValidOtp {
		resetRecord.AttemptCount += 1
		if err := repository.UpdatePasswordReset(resetRecord); err != nil {
			utils.ResponseError(c, http.StatusInternalServerError, "Failed to update OTP attempt count", nil)
		}
		utils.ResponseError(c, http.StatusBadRequest, "Invalid OTP", nil)
		return
	}

	// create reset token (JWT or UUID) and return to user
	resetTokenClaim := utils.JWTClaims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "e-commerce-app",
		},
	}
	resetToken, err := utils.PasswordResetToken(resetTokenClaim)

	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create reset token", nil)
		return
	}

	// invalidate OTP after successful verification
	resetRecord.OTPHash = ""
	resetRecord.AttemptCount = 0 // reset attempt count on success

	if err := repository.UpdatePasswordReset(resetRecord); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "OTP verified successfully", gin.H{
		"reset_token": resetToken,
	})
}

func ResendOtpRequest(c *gin.Context) {
	// reset attempt to 0 and generate new OTP and send to user email
}

func PasswordReset(c *gin.Context) {
	// This handler would allow the user to reset their password after verifying the OTP.
	// It would take the reset token generated in the VerifyOtpRequest handler, verify it,
	// and then allow the user to set a new password. After resetting, it should invalidate the reset token.
	// Implementation would involve updating the user's password in the database and ensuring that the reset token cannot be reused.
	// validate reset token, hash new password and update user record, invalidate reset token
	// delete or invalidate the password reset record after successful password reset
}
