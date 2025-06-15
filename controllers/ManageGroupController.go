package controllers

import (
	"be-awarenix/config"
	"be-awarenix/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetGroups(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	// Calculate offset
	offset := (page - 1) * pageSize

	// Build query
	query := config.DB.Model(&models.Group{})

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
			"Message": "Failed to count groups",
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

	// Get groups with pagination
	var groups []models.Group
	if err := query.Order(orderClause).Offset(offset).Limit(pageSize).Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Success": false,
			"Message": "Failed to fetch groups",
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
		"Message": "Groups retrieved successfully",
		"Data":    groups,
		"Total":   total,
	})
}
