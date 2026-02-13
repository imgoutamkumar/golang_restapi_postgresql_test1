package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
)

type UploadFileToCloudinaryResponse struct {
	ImageUrl  string
	Public_Id string
}

func UploadFileToCloudinary(fileHeader *multipart.FileHeader) (*UploadFileToCloudinaryResponse, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	uploadResult, err := config.CLD.Upload.Upload(
		context.Background(),
		file,
		uploader.UploadParams{})
	if err != nil {
		log.Printf("Failed to upload file to Cloudinary: %v", err)
		return nil, err
	}

	// return uploadResult.SecureURL, nil
	response := &UploadFileToCloudinaryResponse{
		ImageUrl:  uploadResult.SecureURL,
		Public_Id: uploadResult.PublicID,
	}
	return response, nil
}

func DeleteFileFromCloudinary(url string) error {
	publicID := "abc"
	_, err := config.CLD.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID: publicID,
		},
	)

	if err != nil {
		return fmt.Errorf("cloudinary delete failed: %w", err)
	}

	return nil
}

// optional
// uploader.UploadParams{
// 		Folder: "products",
//        PublicID: "product_" + uuid.New().String(),
// })
