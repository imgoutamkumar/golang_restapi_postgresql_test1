package dto

import "mime/multipart"

type CreateBannerRequest struct {
	Name     string `form:"name" binding:"required"`
	Type     string `form:"type" binding:"required"`
	IsActive bool   `form:"is_active"`

	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`

	// multiple images
	Images []*multipart.FileHeader `form:"images"`
}