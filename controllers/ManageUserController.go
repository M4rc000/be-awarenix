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
