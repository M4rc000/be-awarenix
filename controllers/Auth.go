package controllers

import (
	"be-awarenix/config"
	"be-awarenix/models"
	"be-awarenix/services"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthLogin(c *gin.Context) {
	var input models.LoginInput
	var fullUserData models.FullUserLoginData
	var userResp models.UserLoginResponse

	if err := c.ShouldBindJSON(&input); err != nil {
		services.LogActivity(config.DB, c, "Login", "Auth", "", nil, input, "failed", "Invalid input: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	// Coba cari user dengan data role
	err := config.DB.Table("users").
		Select(`users.*, roles.name AS role_name`).
		Joins(`LEFT JOIN roles ON roles.id = users.role`).
		Where("users.email = ?", input.Email).
		First(&userResp).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Log ketika user tidak ditemukan, recordID kosong atau bisa pakai email jika mau
			services.LogActivity(config.DB, c, "Login", "Auth", "", nil, input, "failed", "Account haven't registered yet")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Account haven't registered yet",
				"error":   "User not found",
			})
		} else {
			// Log error database lain
			log.Printf("Database error during login: %v", err)
			services.LogActivity(config.DB, c, "Login", "Auth", "", nil, input, "failed", "Failed to process login: "+err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to process login",
				"error":   err.Error(),
			})
		}
		return
	}

	// Mengambil hash password secara terpisah jika diperlukan, namun First(&userResp) seharusnya sudah cukup
	// Jika models.UserLoginResponse tidak memiliki PasswordHash, maka bagian ini diperlukan.
	var userWithHash models.User
	err = config.DB.Where("email = ?", input.Email).First(&userWithHash).Error
	if err != nil {
		log.Printf("Error fetching user hash: %v", err)
		services.LogActivity(config.DB, c, "Login", "Auth", fmt.Sprintf("%v", userResp.ID), nil, input, "failed", "Failed to retrieve user data for hash: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve user data",
			"error":   err.Error(),
		})
		return
	}

	// Jika userResp dan fullUserData adalah model yang berbeda dan memerlukan query terpisah, ini bisa tetap ada.
	// Jika userResp sudah cukup, baris ini bisa dihilangkan atau diganti dengan penugasan langsung.
	err = config.DB.Table("users").
		Select(`users.*, roles.name AS role_name`).
		Joins(`LEFT JOIN roles ON roles.id = users.role`).
		Where("users.email = ?", input.Email).
		First(&fullUserData).Error // Ini akan mengisi fullUserData

	if err != nil {
		log.Printf("Database error during login (fullUserData fetch): %v", err)
		services.LogActivity(config.DB, c, "Login", "Auth", fmt.Sprintf("%v", userResp.ID), nil, input, "failed", "Failed to process login (full data fetch): "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to process login",
			"error":   err.Error(),
		})
		return
	}

	fullUserData.ID = userResp.ID
	fullUserData.Email = userResp.Email
	fullUserData.Name = userResp.Name
	fullUserData.Position = userResp.Position
	fullUserData.Role = userResp.Role
	fullUserData.RoleName = userResp.RoleName
	fullUserData.Company = userResp.Company
	fullUserData.Country = userResp.Country
	fullUserData.IsActive = userResp.IsActive
	fullUserData.PasswordHash = userWithHash.PasswordHash

	if fullUserData.IsActive == 0 {
		services.LogActivity(config.DB, c, "Login", "Auth", fmt.Sprintf("%v", fullUserData.ID), nil, input, "failed", "Account is not active")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Account is not active",
			"error":   "Account is inactive",
		})
		return
	}

	if err := services.ComparePassword(fullUserData.PasswordHash, input.Password); err != nil {
		services.LogActivity(config.DB, c, "Login", "Auth", fmt.Sprintf("%v", fullUserData.ID), nil, input, "failed", "Invalid credentials")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid credentials",
			"error":   "Password mismatch",
		})
		return
	}

	token, exp, err := services.GenerateJWT(fullUserData.ID, fullUserData.Email, input.Status)
	if err != nil {
		services.LogActivity(config.DB, c, "Login", "Auth", fmt.Sprintf("%v", fullUserData.ID), nil, input, "failed", "Could not create token: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Could not create token",
			"error":   err.Error(),
		})
		return
	}

	fullUserData.LastLogin = time.Now()
	if err := config.DB.Save(&fullUserData.User).Error; err != nil {
		log.Printf("Failed to update last_login: %v", err)
		services.LogActivity(config.DB, c, "Login", "Auth", fmt.Sprintf("%v", fullUserData.ID), nil, input, "warning", "Failed to update last_login: "+err.Error())
	}

	// Siapkan data untuk response
	userdata := map[string]interface{}{
		"id":         fullUserData.ID,
		"name":       fullUserData.Name,
		"email":      fullUserData.Email,
		"position":   fullUserData.Position,
		"role":       fullUserData.Role,
		"role_name":  fullUserData.RoleName,
		"company":    fullUserData.Company,
		"country":    fullUserData.Country,
		"last_login": fullUserData.LastLogin,
	}

	userid := int(fullUserData.ID)
	services.LogActivity(config.DB, c, "Login", "Auth", strconv.Itoa(userid), nil, userdata, "success", "Login successful")
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Login successful",
		"token":      token,
		"user":       userdata,
		"expires_at": exp,
	})
}

func AuthLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func GetUserSession(c *gin.Context) {
	var input models.GetUserSession

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Success": false,
			"Message": "Invalid request body",
			"Error":   err.Error(),
		})
		return
	}

	var user models.User
	if err := config.DB.Select("id", "name", "email", "position", "role").First(&user, input.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Success": false,
			"Message": "User not found",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Success": true,
		"Message": "User session retrieved successfully",
		"Data": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"position": user.Position,
			"role":     user.Role,
		},
	})
}
