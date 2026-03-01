package dto

type UserFilterParams struct {
	ProductName string
	FullName    string
	Page        int
	Limit       int
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email validate:"required,email"`
}

type VerifyOtpRequestBody struct {
	Email string `json:"email" binding:"required,email"`
	Otp   string `json:"otp" binding:"required,len=6"`
}