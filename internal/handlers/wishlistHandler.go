package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/dto"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/services"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

func AddItemToWishlist(c *gin.Context) {
	var req dto.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	val, ok := c.Get("userId")
	if !ok {
		utils.ResponseError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userId, err := uuid.Parse(val.(string))
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid User ID", err.Error())
		return
	}
	wishlist, err := repository.GetDefaultWishlistByUserId(userId)
	if err != nil {
		wishlist, err = repository.CreateWishlist(&models.Wishlist{
			UserID: userId.String(),
			Name:   "My Wishlist",
		})
		if err != nil {
			utils.ResponseError(c, http.StatusInternalServerError, "Failed to create wishlist", err.Error())
			return
		}
	}

	isExist, err := services.IsItemExists(wishlist.ID, req.VariantID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to check wishlist", err.Error())
		return
	}
	if isExist {
		utils.ResponseError(c, http.StatusConflict, "Item already exists in wishlist", nil)
		return
	}
	var wishlistItem models.WishlistItem
	wishlistItem.WishlistID = wishlist.ID
	wishlistItem.VariantID = req.VariantID

	if err := repository.AddNewWishlistItem(&wishlistItem); err != nil {
		// Handle race condition (duplicate insert)
		if strings.Contains(err.Error(), "duplicate") {
			utils.ResponseError(c, http.StatusConflict, "Item already exists in wishlist", nil)
			return
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Wishlist Item Failed", err.Error())
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "Wishlist Item Added", nil)
}

func RemoveItemFromWishlist(c *gin.Context) {
	var req dto.RemoveFromWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	val, ok := c.Get("userId")
	if !ok {
		utils.ResponseError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	userId, err := uuid.Parse(val.(string))
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid User ID", err.Error())
		return
	}
	wishlist, err := repository.GetWishlistByUserId(userId)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to get wishlist", err.Error())
		return
	}
	if err := repository.RemoveWishlistItem(wishlist.ID, req.VariantID); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to remove item from wishlist", err.Error())
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "Wishlist Item Removed", nil)
}
