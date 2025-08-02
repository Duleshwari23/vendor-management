package handlers

import (
	"net/http"
	"time"
	"vendor-management/models"

	"github.com/gin-gonic/gin"
)

type CreateVendorRequest struct {
	CompanyName string `json:"companyName" binding:"required"`
	JoiningDate string `json:"joiningDate" binding:"required"`
	EndDate     string `json:"endDate"`
	Department  string `json:"department" binding:"required"`
	ProjectName string `json:"projectName" binding:"required"`
}

func CreateVendor(c *gin.Context) {
	var req CreateVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	joiningDate, err := time.Parse("2006-01-02", req.JoiningDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid joining date format"})
		return
	}

	var endDate time.Time
	if req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}
	}

	vendor := &models.Vendor{
		ID:          generateID(),
		CompanyName: req.CompanyName,
		JoiningDate: joiningDate,
		EndDate:     endDate,
		Department:  req.Department,
		ProjectName: req.ProjectName,
		Status:      "active",
		Documents:   make([]models.Document, 0),
		Assets:      make([]models.Asset, 0),
	}

	models.Vendors[vendor.ID] = vendor
	c.JSON(http.StatusCreated, vendor)
}

func ListVendors(c *gin.Context) {
	vendors := make([]*models.Vendor, 0)
	for _, vendor := range models.Vendors {
		vendors = append(vendors, vendor)
	}
	c.JSON(http.StatusOK, vendors)
}

func GetVendor(c *gin.Context) {
	id := c.Param("id")
	vendor, exists := models.Vendors[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}

	// Find associated assets
	vendorAssets := make([]models.Asset, 0)
	for _, asset := range models.Assets {
		if asset.AssignedTo == id {
			vendorAssets = append(vendorAssets, *asset)
		}
	}

	// Find associated attendance
	vendorAttendance := make([]*models.Attendance, 0)
	if records, ok := models.AttendanceRecords[id]; ok {
		vendorAttendance = records
	}

	c.JSON(http.StatusOK, gin.H{
		"vendor":     vendor,
		"assets":     vendorAssets,
		"attendance": vendorAttendance,
	})
}

func UpdateVendor(c *gin.Context) {
	id := c.Param("id")
	vendor, exists := models.Vendors[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}

	var req CreateVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	joiningDate, err := time.Parse("2006-01-02", req.JoiningDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid joining date format"})
		return
	}

	var endDate time.Time
	if req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}
	}

	vendor.CompanyName = req.CompanyName
	vendor.JoiningDate = joiningDate
	vendor.EndDate = endDate
	vendor.Department = req.Department
	vendor.ProjectName = req.ProjectName

	c.JSON(http.StatusOK, vendor)
}

func GetProfile(c *gin.Context) {
	userID, _ := c.Get("userId")
	user, exists := models.Users[userID.(string)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// For vendors, include their vendor details
	if user.Role == models.VendorRole {
		for _, vendor := range models.Vendors {
			if vendor.UserID == user.ID {
				c.JSON(http.StatusOK, gin.H{
					"user":   user,
					"vendor": vendor,
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
