package services

import (
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/dto"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
)

func GetProductImages(productId string) ([]dto.ProductImageResponse, error) {
	images, err := repository.GetImagesByProductID(productId)

	responseImages := []dto.ProductImageResponse{}

	for _, img := range images {
		responseImages = append(responseImages, dto.ProductImageResponse{
			Id:        img.ID.String(),
			URL:       img.ImageUrl,
			IsPrimary: img.IsPrimary,
			PublicId:  img.PublicId,
		})
	}

	return responseImages, err
}
