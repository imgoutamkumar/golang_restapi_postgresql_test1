package config

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

// CLOUDINARY_URL=cloudinary://<your_api_key>:<your_api_secret>@dfd9vbo0o
var CLD *cloudinary.Cloudinary

func InitializeCloudinary() {
	env := LoadEnv()

	url := "cloudinary://" + env.CLOUDINARTY_API_KEY + ":" +
		env.CLOUDINARTY_API_SECRET + "@" +
		env.CLOUDINARTY_CLOUD_NAME

	var err error
	CLD, err = cloudinary.NewFromURL(url)
	if err != nil {
		log.Fatal("Cloudinary init failed:", err)
	}
}
