package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"vendor-management/handlers"
	"vendor-management/middleware"
	"vendor-management/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	r := gin.Default()

	// Auth routes
	r.POST("/api/auth/signup", handlers.HandleSignup)
	r.POST("/api/auth/login", handlers.HandleLogin)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		admin := api.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/vendors", handlers.CreateVendor)
			admin.GET("/vendors", handlers.ListVendors)
		}

		api.GET("/profile", handlers.GetProfile)
		api.GET("/my-attendance", handlers.GetMyAttendance)
	}

	return r
}

func TestVendorManagementFlow(t *testing.T) {
	router := setupTestRouter()

	// 1. Login as admin user (default admin is created by models.init)
	adminLogin := map[string]interface{}{
		"email":    "admin@company.com",
		"password": "admin",
	}
	adminLoginJSON, _ := json.Marshal(adminLogin)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(adminLoginJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var adminResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &adminResponse)
	assert.NoError(t, err)
	adminToken := adminResponse["token"].(string)

	// 2. Create vendor user
	vendorSignup := map[string]interface{}{
		"name":     "Vendor User",
		"email":    "vendor@example.com",
		"password": "vendor123",
		"role":     "vendor",
	}
	vendorJSON, _ := json.Marshal(vendorSignup)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(vendorJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var vendorResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &vendorResponse)
	assert.NoError(t, err)
	vendorToken := vendorResponse["token"].(string)

	// 3. Create vendor profile (as admin)
	vendor := map[string]interface{}{
		"companyName": "Test Vendor Company",
		"joiningDate": time.Now().Format(time.RFC3339),
		"department":  "IT",
		"projectName": "Test Project",
	}
	vendorCreateJSON, _ := json.Marshal(vendor)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/admin/vendors", bytes.NewBuffer(vendorCreateJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 4. List vendors (as admin)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/admin/vendors", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var vendors []models.Vendor
	err = json.Unmarshal(w.Body.Bytes(), &vendors)
	assert.NoError(t, err)
	assert.Len(t, vendors, 1)
	assert.Equal(t, "Test Vendor Company", vendors[0].CompanyName)

	// 5. Get profile (as vendor)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+vendorToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var profileResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &profileResponse)
	assert.NoError(t, err)
	user := profileResponse["user"].(map[string]interface{})
	assert.Equal(t, "Vendor User", user["name"])
	assert.Equal(t, "vendor", user["role"])

	// 6. Get attendance (as vendor)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/my-attendance", nil)
	req.Header.Set("Authorization", "Bearer "+vendorToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// Optionally, check attendance structure
	var attendance []models.Attendance
	err = json.Unmarshal(w.Body.Bytes(), &attendance)
	assert.NoError(t, err)
}
