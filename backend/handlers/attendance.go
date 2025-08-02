package handlers

import (
	"net/http"
	"time"
	"vendor-management/models"
	"vendor-management/utils"

	"github.com/gin-gonic/gin"
)

func ListAttendance(c *gin.Context) {
	// Optional date range filters
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	var start, end time.Time
	var err error
	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}
	}
	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}
	}

	result := make(map[string][]models.Attendance)
	for vendorID, attendances := range models.AttendanceRecords {
		filtered := make([]models.Attendance, 0)
		for _, a := range attendances {
			if startDate != "" && a.Date.Before(start) {
				continue
			}
			if endDate != "" && a.Date.After(end) {
				continue
			}
			filtered = append(filtered, *a)
		}
		if len(filtered) > 0 {
			result[vendorID] = filtered
		}
	}

	c.JSON(http.StatusOK, result)
}

func GetVendorAttendance(c *gin.Context) {
	vendorID := c.Param("vendorId")

	// Verify vendor exists
	if _, exists := models.Vendors[vendorID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}

	// Get attendance records
	attendance := models.AttendanceRecords[vendorID]
	if attendance == nil {
		attendance = make([]*models.Attendance, 0)
	}
	c.JSON(http.StatusOK, attendance)
}

func GetMyAttendance(c *gin.Context) {
	userID, _ := c.Get("userId")
	user, exists := models.Users[userID.(string)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role != models.VendorRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only vendors can view their attendance"})
		return
	}

	// Find vendor's attendance
	for _, vendor := range models.Vendors {
		if vendor.UserID == user.ID {
			attendance := models.AttendanceRecords[vendor.ID]
			if attendance == nil {
				attendance = make([]*models.Attendance, 0)
			}
			c.JSON(http.StatusOK, attendance)
			return
		}
	}

	c.JSON(http.StatusOK, []models.Attendance{})
}

// ManuallyUpdateAttendance triggers the attendance update for testing purposes
func ManuallyUpdateAttendance(c *gin.Context) {
	utils.UpdateAttendance()
	c.JSON(http.StatusOK, gin.H{"message": "Attendance updated successfully"})
}
