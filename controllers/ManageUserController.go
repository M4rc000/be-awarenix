package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"be-awarenix/config"
	"be-awarenix/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// userInput merepresentasikan JSON yang dikirim FE
type userInput struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Position string `json:"position" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserResponse untuk response ke frontend (tanpa password)
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Position  string    `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func RegisterUser(c *gin.Context) {
	var input userInput

	// Bind dan validasi input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": err.Error(),
		})
		return
	}

	// Cek apakah email sudah digunakan
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Email already exists",
			"message": "User with this email already registered",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Password hashing failed",
			"message": "Failed to process password",
		})
		return
	}

	// Buat user baru
	newUser := models.User{
		Name:         input.Name,
		Email:        input.Email,
		Position:     input.Position,
		PasswordHash: string(hashedPassword),
	}

	// Simpan ke database
	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"message": "Failed to create user",
		})
		return
	}

	// Siapkan response (tanpa password)
	userResponse := UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Position:  newUser.Position,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	// Response sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    userResponse,
	})
}

func GetUsers(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	// Calculate offset
	offset := (page - 1) * pageSize

	// Build query
	query := config.DB.Model(&models.User{})

	// Add search conditions
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(position) LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Success": false,
			"Message": "Failed to count users",
			"Error":   err.Error(),
		})
		return
	}

	// Add sorting
	orderClause := sortBy
	if sortOrder == "desc" {
		orderClause += " DESC"
	} else {
		orderClause += " ASC"
	}

	// Get users with pagination
	var users []models.User
	if err := query.Order(orderClause).Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Success": false,
			"Message": "Failed to fetch users",
			"Error":   err.Error(),
		})
		return
	}

	// Calculate pagination info
	// totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	// hasNextPage := page < totalPages
	// hasPreviousPage := page > 1

	// pagination := PaginationInfo{
	// 	CurrentPage:     page,
	// 	PageSize:        pageSize,
	// 	TotalItems:      total,
	// 	TotalPages:      totalPages,
	// 	HasNextPage:     hasNextPage,
	// 	HasPreviousPage: hasPreviousPage,
	// }

	// response := gin{
	// 	Users:      users,
	// 	Pagination: pagination,
	// }

	c.JSON(http.StatusOK, gin.H{
		"Success": true,
		"Message": "Users retrieved successfully",
		"Data":    users,
	})
}

func DeleteUser(c *gin.Context) {
	// Get user ID from URL parameter
	userID := c.Param("id")

	// Validate user ID
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID format",
			"error":   "User ID must be a valid number",
		})
		return
	}

	// Get current user from JWT token (from middleware)
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized access",
			"error":   "User session not found",
		})
		return
	}

	// Type assertion to get user data
	user := currentUser.(*models.User)

	// Prevent user from deleting themselves
	if user.ID == uint(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Cannot delete your own account",
			"error":   "Self-deletion is not allowed",
		})
		return
	}

	// Check if user to be deleted exists
	var userToDelete models.User
	if err := config.DB.First(&userToDelete, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found",
				"error":   "The specified user does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
			"error":   err.Error(),
		})
		return
	}

	// Start database transaction for safe deletion
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to start transaction",
			"error":   tx.Error.Error(),
		})
		return
	}

	// Hard Delete user (permanently remove from database)
	if err := tx.Unscoped().Delete(&userToDelete).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
		return
	}

	// Optional: Delete related data (sessions, logs, etc.) - HARD DELETE
	// Example: Delete user sessions permanently
	if err := tx.Unscoped().Where("user_id = ?", id).Delete(&models.UserSession{}).Error; err != nil {
		// Log error but don't fail the deletion
		// You might want to handle this differently based on your needs
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to commit transaction",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
		"data": gin.H{
			"deleted_user": gin.H{
				"id":       userToDelete.ID,
				"name":     userToDelete.Name,
				"email":    userToDelete.Email,
				"position": userToDelete.Position,
			},
		},
	})
}
