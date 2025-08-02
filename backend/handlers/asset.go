package handlers

import (
	"net/http"
	"time"
	"vendor-management/models"

	"github.com/gin-gonic/gin"
)

type CreateAssetRequest struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	SerialNumber string `json:"serialNumber"`
	VendorID     string `json:"vendor_id"`
}

type AssignAssetRequest struct {
	VendorID string `json:"vendorId" binding:"required"`
}

func CreateAsset(c *gin.Context) {
	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset := &models.Asset{
		ID:           generateID(),
		Name:         req.Name,
		Type:         req.Type,
		SerialNumber: req.SerialNumber,
		Status:       "available", // Default status
	}

	if req.VendorID != "" {
		vendor, exists := models.Vendors[req.VendorID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
			return
		}
		asset.AssignedTo = req.VendorID
		asset.Status = "assigned"
		vendor.Assets = append(vendor.Assets, *asset)
	}

	models.Assets[asset.ID] = asset
	c.JSON(http.StatusCreated, asset)
}

func ListAssets(c *gin.Context) {
	assets := make([]*models.Asset, 0)
	for _, asset := range models.Assets {
		assets = append(assets, asset)
	}
	c.JSON(http.StatusOK, assets)
}

func UpdateAsset(c *gin.Context) {
	id := c.Param("id")
	asset, exists := models.Assets[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset.Name = req.Name
	asset.Type = req.Type
	asset.SerialNumber = req.SerialNumber

	c.JSON(http.StatusOK, asset)
}

func AssignAsset(c *gin.Context) {
	id := c.Param("id")
	asset, exists := models.Assets[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	if asset.Status != "available" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Asset is not available"})
		return
	}

	var req AssignAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vendor, exists := models.Vendors[req.VendorID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}

	asset.AssignedTo = req.VendorID
	asset.AssignedAt = time.Now()
	asset.Status = "assigned"

	vendor.Assets = append(vendor.Assets, *asset)

	c.JSON(http.StatusOK, asset)
}

func ReturnAsset(c *gin.Context) {
	id := c.Param("id")
	asset, exists := models.Assets[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	if asset.Status != "assigned" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Asset is not assigned"})
		return
	}

	// Remove asset from vendor's assets
	if vendor, exists := models.Vendors[asset.AssignedTo]; exists {
		newAssets := make([]models.Asset, 0)
		for _, a := range vendor.Assets {
			if a.ID != asset.ID {
				newAssets = append(newAssets, a)
			}
		}
		vendor.Assets = newAssets
	}

	asset.ReturnedAt = time.Now()
	asset.Status = "available"
	asset.AssignedTo = ""

	c.JSON(http.StatusOK, asset)
}

func GetMyAssets(c *gin.Context) {
	userID, _ := c.Get("userId")
	user, exists := models.Users[userID.(string)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role != models.VendorRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only vendors can view their assets"})
		return
	}

	// Find vendor's assets
	for _, vendor := range models.Vendors {
		if vendor.UserID == user.ID {
			c.JSON(http.StatusOK, vendor.Assets)
			return
		}
	}

	c.JSON(http.StatusOK, []models.Asset{})
}
