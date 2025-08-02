package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	AdminRole  Role = "admin"
	VendorRole Role = "vendor"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never sent in JSON responses
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type Vendor struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	CompanyName string     `json:"companyName"`
	JoiningDate time.Time  `json:"joiningDate"`
	EndDate     time.Time  `json:"endDate,omitempty"`
	Department  string     `json:"department"`
	ProjectName string     `json:"projectName"`
	Status      string     `json:"status"` // active, inactive
	Documents   []Document `json:"documents"`
	Assets      []Asset    `json:"assets"`
}

type Document struct {
	ID         string    `json:"id"`
	VendorID   string    `json:"vendorId"`
	Name       string    `json:"name"`
	Type       string    `json:"type"` // joining_letter, agreement, id_proof
	FilePath   string    `json:"filePath"`
	UploadedAt time.Time `json:"uploadedAt"`
}

type Asset struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"` // laptop, monitor, keyboard, etc.
	SerialNumber string    `json:"serialNumber"`
	AssignedTo   string    `json:"assignedTo"` // VendorID
	AssignedAt   time.Time `json:"assignedAt"`
	ReturnedAt   time.Time `json:"returnedAt,omitempty"`
	Status       string    `json:"status"` // assigned, available, maintenance
}

type Attendance struct {
	ID         string    `json:"id"`
	VendorID   string    `json:"vendorId"`
	Date       time.Time `json:"date"`
	LoginTime  time.Time `json:"loginTime"`
	LogoutTime time.Time `json:"logoutTime"`
	PresentDay float32   `json:"presentDay"` // 0, 0.5, 1.0 for absent, half-day, full-day
	Status     string    `json:"status"`     // Present, Absent, Leave
}

// In-memory storage (to be replaced with a real database later)
var (
	Users             = make(map[string]*User)
	Vendors           = make(map[string]*Vendor)
	Documents         = make(map[string]*Document)
	Assets            = make(map[string]*Asset)
	AttendanceRecords = make(map[string][]*Attendance) // map[vendorID][]Attendance
)

func init() {
	// Create default admin user if not present
	for _, u := range Users {
		if u.Role == AdminRole {
			return // Admin already exists
		}
	}
	adminID := generateID()
	adminUser := &User{
		ID:        adminID,
		Name:      "Admin",
		Email:     "admin@company.com",
		Role:      AdminRole,
		CreatedAt: time.Now(),
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	adminUser.Password = string(hashedPassword)
	Users[adminID] = adminUser
}

// generateID generates a random ID (reused from handlers package)
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
