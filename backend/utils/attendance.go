package utils

import (
	crand "crypto/rand"
	"encoding/hex"
	rand "math/rand"
	"time"
	"vendor-management/models"
)

// GenerateRandomAttendance generates random attendance data for the last 5 days
func GenerateRandomAttendance(vendorID string) []*models.Attendance {
	attendance := make([]*models.Attendance, 0)
	now := time.Now()

	// Generate attendance for last 5 days
	for i := 1; i <= 5; i++ {
		date := now.AddDate(0, 0, -i)
		status := generateRandomStatus()

		// Generate random login time between 8:00 AM and 10:00 AM
		loginTime := time.Date(date.Year(), date.Month(), date.Day(), 8+rand.Intn(3), rand.Intn(60), 0, 0, time.Local)

		// Generate random logout time between 5:00 PM and 7:00 PM
		logoutTime := time.Date(date.Year(), date.Month(), date.Day(), 17+rand.Intn(3), rand.Intn(60), 0, 0, time.Local)

		// Calculate duration in hours
		duration := logoutTime.Sub(loginTime).Hours()

		attendance = append(attendance, &models.Attendance{
			ID:         generateID(),
			VendorID:   vendorID,
			Date:       date,
			LoginTime:  loginTime,
			LogoutTime: logoutTime,
			Duration:   float32(duration),
			Status:     status,
		})
	}

	return attendance
}

// generateRandomStatus returns a random attendance status (0, 0.5, or 1.0)
func generateRandomStatus() float32 {
	statuses := []float32{0.0, 0.5, 1.0}
	weights := []int{10, 20, 70} // 10% absent, 20% half-day, 70% full-day

	randNum := rand.Intn(100)
	sum := 0
	for i, weight := range weights {
		sum += weight
		if randNum < sum {
			return statuses[i]
		}
	}
	return 1.0 // Default to full day if something goes wrong
}

// generateID generates a random ID (reused from handlers package)
func generateID() string {
	bytes := make([]byte, 16)
	crand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// UpdateAttendance updates attendance for all vendors
func UpdateAttendance() {
	for _, vendor := range models.Vendors {
		if vendor.Status == "active" {
			// Get existing attendance records for vendor
			existingAttendance := models.AttendanceRecords[vendor.ID]
			if existingAttendance == nil {
				existingAttendance = make([]*models.Attendance, 0)
			}

			// Generate new attendance records
			newAttendance := GenerateRandomAttendance(vendor.ID)

			// Remove any existing records for the same dates
			filteredAttendance := make([]*models.Attendance, 0)
			for _, existing := range existingAttendance {
				isOld := false
				for _, new := range newAttendance {
					if existing.Date.Equal(new.Date) {
						isOld = true
						break
					}
				}
				if !isOld {
					filteredAttendance = append(filteredAttendance, existing)
				}
			}

			// Combine filtered existing records with new records
			models.AttendanceRecords[vendor.ID] = append(filteredAttendance, newAttendance...)
		}
	}
}
