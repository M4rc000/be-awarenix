package controllers

import (
	"be-awarenix/config"
	"be-awarenix/models"
	"be-awarenix/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type loginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Status   string `json:"status"`
}

func AuthLogin(c *gin.Context) {
	var input loginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := services.ComparePassword(user.PasswordHash, input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, exp, err := services.GenerateJWT(user.ID, user.Email, input.Status)
	if err != nil {
		log.Printf("JWT generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	user.LastLogin = time.Now().In(services.JakartaLocation)
	if err := config.DB.Save(&user).Error; err != nil {
		log.Printf("Failed to update last_login: %v", err)
	}

	userdata := map[string]interface{}{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"position": user.Position,
		"role":     user.Role,
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": userdata, "expires_at": exp})
}

func AuthLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func GetUserSession(c *gin.Context) {
	// Ambil dari context
	uidRaw, _ := c.Get("userID")
	emailRaw, _ := c.Get("userEmail")
	userID := uidRaw.(uint)
	email := emailRaw.(string)

	// Optional: fetch full user record
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		// fallback: kirim minimal data
		c.JSON(http.StatusOK, gin.H{
			"id":    userID,
			"email": email,
		})
		return
	}

	// Kirim hanya field yang diperlukan
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"position": user.Position,
	})
}
