package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"vendor-management/models"

	"github.com/gin-gonic/gin"
)

const uploadDir = "uploads"

func init() {
	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
}

func UploadDocument(c *gin.Context) {
	vendorID := c.PostForm("vendorId")
	docType := c.PostForm("type")

	vendor, exists := models.Vendors[vendorID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s%s", vendorID, generateID(), ext)
	filepath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	doc := &models.Document{
		ID:         generateID(),
		VendorID:   vendorID,
		Name:       file.Filename,
		Type:       docType,
		FilePath:   filepath,
		UploadedAt: time.Now(),
	}

	models.Documents[doc.ID] = doc
	vendor.Documents = append(vendor.Documents, *doc)

	c.JSON(http.StatusCreated, doc)
}

func ListDocuments(c *gin.Context) {
	vendorID := c.Query("vendorId")
	if vendorID != "" {
		vendor, exists := models.Vendors[vendorID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
			return
		}
		c.JSON(http.StatusOK, vendor.Documents)
		return
	}

	// Return all documents if no vendor ID specified
	docs := make([]*models.Document, 0)
	for _, doc := range models.Documents {
		docs = append(docs, doc)
	}
	c.JSON(http.StatusOK, docs)
}

func GetDocument(c *gin.Context) {
	id := c.Param("id")
	doc, exists := models.Documents[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Check if file exists
	if _, err := os.Stat(doc.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(doc.FilePath)
}

func DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	doc, exists := models.Documents[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Remove document from vendor's documents
	if vendor, exists := models.Vendors[doc.VendorID]; exists {
		newDocs := make([]models.Document, 0)
		for _, d := range vendor.Documents {
			if d.ID != doc.ID {
				newDocs = append(newDocs, d)
			}
		}
		vendor.Documents = newDocs
	}

	// Delete file from disk
	if err := os.Remove(doc.FilePath); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	delete(models.Documents, id)
	c.Status(http.StatusNoContent)
}

func GetMyDocuments(c *gin.Context) {
	userID, _ := c.Get("userId")
	user, exists := models.Users[userID.(string)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role != models.VendorRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only vendors can view their documents"})
		return
	}

	// Find vendor's documents
	for _, vendor := range models.Vendors {
		if vendor.UserID == user.ID {
			c.JSON(http.StatusOK, vendor.Documents)
			return
		}
	}

	c.JSON(http.StatusOK, []models.Document{})
}
