// controllers/campaigns.go
package controllers

import (
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

// Catatan: Semua helper response SendSuccessResponse, SendErrorResponse, SendValidationErrorResponse
// akan dihapus dari penggunaan dalam file ini, dan diganti dengan c.JSON langsung.

// RegisterCampaign handles the creation of a new campaign
func RegisterCampaign(c *gin.Context) {
	var input models.CampaignRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := services.ParseValidationErrors(err)
		if validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Validasi gagal",
				"fields":  validationErrors,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Parsing tanggal dari string ke time.Time
	launchDate, err := time.Parse(time.RFC3339, input.LaunchDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format Launch Date tidak valid. Gunakan format RFC3339.",
			"fields":  map[string]string{"launch_date": "Format tanggal tidak valid"},
		})
		return
	}

	var sendEmailBy *time.Time
	if input.SendEmailBy != nil && *input.SendEmailBy != "" {
		parsedSendEmailBy, err := time.Parse(time.RFC3339, *input.SendEmailBy)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Format Send Email By tidak valid. Gunakan format RFC3339.",
				"fields":  map[string]string{"send_email_by": "Format tanggal tidak valid"},
			})
			return
		}
		sendEmailBy = &parsedSendEmailBy
	} else {
		sendEmailBy = nil
	}

	// Dapatkan instance DB dari context
	db := c.MustGet("db").(*gorm.DB)

	// Verifikasi keberadaan Group, EmailTemplate, LandingPage, SendingProfile
	var group models.Group
	if err := db.First(&group, input.GroupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Group ID tidak ditemukan",
			"fields":  map[string]string{"group_id": "Group tidak ada"},
		})
		return
	}

	var emailTemplate models.EmailTemplate
	if err := db.First(&emailTemplate, input.EmailTemplateID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Email Template ID tidak ditemukan",
			"fields":  map[string]string{"email_template_id": "Template email tidak ada"},
		})
		return
	}

	var landingPage models.LandingPage
	if err := db.First(&landingPage, input.LandingPageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Landing Page ID tidak ditemukan",
			"fields":  map[string]string{"landing_page_id": "Landing page tidak ada"},
		})
		return
	}

	var sendingProfile models.SendingProfiles
	if err := db.First(&sendingProfile, input.SendingProfileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Sending Profile ID tidak ditemukan",
			"fields":  map[string]string{"sending_profile_id": "Profil pengiriman tidak ada"},
		})
		return
	}

	campaign := models.Campaign{
		Name:             input.Name,
		LaunchDate:       launchDate,
		SendEmailBy:      sendEmailBy,
		GroupID:          input.GroupID,
		EmailTemplateID:  input.EmailTemplateID,
		LandingPageID:    input.LandingPageID,
		SendingProfileID: input.SendingProfileID,
		URL:              input.URL,
		CreatedBy:        input.CreatedBy,
		CreatedAt:        time.Now(),
		Status:           "draft", // Status awal kampanye
	}

	if err := db.Create(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat kampanye: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kampanye berhasil didaftarkan",
		"data": models.CampaignResponse{
			ID:               campaign.ID,
			Name:             campaign.Name,
			LaunchDate:       campaign.LaunchDate,
			SendEmailBy:      campaign.SendEmailBy,
			GroupID:          campaign.GroupID,
			EmailTemplateID:  campaign.EmailTemplateID,
			LandingPageID:    campaign.LandingPageID,
			SendingProfileID: campaign.SendingProfileID,
			URL:              campaign.URL,
			CreatedBy:        campaign.CreatedBy,
			CreatedAt:        campaign.CreatedAt,
			UpdatedAt:        campaign.UpdatedAt,
			Status:           campaign.Status,
			CompletedDate:    campaign.CompletedDate,
		},
	})
}

// GetCampaigns retrieves all campaigns
func GetCampaigns(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var campaigns []models.Campaign
	if err := db.Find(&campaigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil kampanye: " + err.Error(),
		})
		return
	}

	var campaignResponses []models.CampaignResponse
	for _, campaign := range campaigns {
		campaignResponses = append(campaignResponses, models.CampaignResponse{
			ID:               campaign.ID,
			Name:             campaign.Name,
			LaunchDate:       campaign.LaunchDate,
			SendEmailBy:      campaign.SendEmailBy,
			GroupID:          campaign.GroupID,
			EmailTemplateID:  campaign.EmailTemplateID,
			LandingPageID:    campaign.LandingPageID,
			SendingProfileID: campaign.SendingProfileID,
			URL:              campaign.URL,
			CreatedBy:        campaign.CreatedBy,
			CreatedAt:        campaign.CreatedAt,
			UpdatedAt:        campaign.UpdatedAt,
			Status:           campaign.Status,
			CompletedDate:    campaign.CompletedDate,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kampanye berhasil diambil",
		"data":    campaignResponses,
	})
}

// GetCampaignDetail retrieves a single campaign by ID
func GetCampaignDetail(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var campaign models.Campaign
	if err := db.First(&campaign, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Kampanye tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil detail kampanye: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail kampanye berhasil diambil",
		"data": models.CampaignResponse{
			ID:               campaign.ID,
			Name:             campaign.Name,
			LaunchDate:       campaign.LaunchDate,
			SendEmailBy:      campaign.SendEmailBy,
			GroupID:          campaign.GroupID,
			EmailTemplateID:  campaign.EmailTemplateID,
			LandingPageID:    campaign.LandingPageID,
			SendingProfileID: campaign.SendingProfileID,
			URL:              campaign.URL,
			CreatedBy:        campaign.CreatedBy,
			CreatedAt:        campaign.CreatedAt,
			UpdatedAt:        campaign.UpdatedAt,
			Status:           campaign.Status,
			CompletedDate:    campaign.CompletedDate,
		},
	})
}

// UpdateCampaign updates an existing campaign
func UpdateCampaign(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var existingCampaign models.Campaign
	if err := db.First(&existingCampaign, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Kampanye tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menemukan kampanye: " + err.Error(),
		})
		return
	}

	var input models.CampaignRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := services.ParseValidationErrors(err)
		if validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Validasi gagal",
				"fields":  validationErrors,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Parsing tanggal dari string ke time.Time
	launchDate, err := time.Parse(time.RFC3339, input.LaunchDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format Launch Date tidak valid. Gunakan format RFC3339.",
			"fields":  map[string]string{"launch_date": "Format tanggal tidak valid"},
		})
		return
	}

	var sendEmailBy *time.Time
	if input.SendEmailBy != nil && *input.SendEmailBy != "" {
		parsedSendEmailBy, err := time.Parse(time.RFC3339, *input.SendEmailBy)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Format Send Email By tidak valid. Gunakan format RFC3339.",
				"fields":  map[string]string{"send_email_by": "Format tanggal tidak valid"},
			})
			return
		}
		sendEmailBy = &parsedSendEmailBy
	} else {
		sendEmailBy = nil
	}

	// Verifikasi keberadaan Group, EmailTemplate, LandingPage, SendingProfile
	var group models.Group
	if err := db.First(&group, input.GroupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Group ID tidak ditemukan",
			"fields":  map[string]string{"group_id": "Group tidak ada"},
		})
		return
	}

	var emailTemplate models.EmailTemplate
	if err := db.First(&emailTemplate, input.EmailTemplateID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Email Template ID tidak ditemukan",
			"fields":  map[string]string{"email_template_id": "Template email tidak ada"},
		})
		return
	}

	var landingPage models.LandingPage
	if err := db.First(&landingPage, input.LandingPageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Landing Page ID tidak ditemukan",
			"fields":  map[string]string{"landing_page_id": "Landing page tidak ada"},
		})
		return
	}

	var sendingProfile models.SendingProfiles
	if err := db.First(&sendingProfile, input.SendingProfileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Sending Profile ID tidak ditemukan",
			"fields":  map[string]string{"sending_profile_id": "Profil pengiriman tidak ada"},
		})
		return
	}

	// Update fields
	existingCampaign.Name = input.Name
	existingCampaign.LaunchDate = launchDate
	existingCampaign.SendEmailBy = sendEmailBy
	existingCampaign.GroupID = input.GroupID
	existingCampaign.EmailTemplateID = input.EmailTemplateID
	existingCampaign.LandingPageID = input.LandingPageID
	existingCampaign.SendingProfileID = input.SendingProfileID
	existingCampaign.URL = input.URL
	existingCampaign.UpdatedAt = time.Now()
	// CreatedBy tidak diubah saat update, UpdatedBy bisa ditambahkan jika ada di struct

	if err := db.Save(&existingCampaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memperbarui kampanye: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kampanye berhasil diperbarui",
		"data": models.CampaignResponse{
			ID:               existingCampaign.ID,
			Name:             existingCampaign.Name,
			LaunchDate:       existingCampaign.LaunchDate,
			SendEmailBy:      existingCampaign.SendEmailBy,
			GroupID:          existingCampaign.GroupID,
			EmailTemplateID:  existingCampaign.EmailTemplateID,
			LandingPageID:    existingCampaign.LandingPageID,
			SendingProfileID: existingCampaign.SendingProfileID,
			URL:              existingCampaign.URL,
			CreatedBy:        existingCampaign.CreatedBy,
			CreatedAt:        existingCampaign.CreatedAt,
			UpdatedAt:        existingCampaign.UpdatedAt,
			Status:           existingCampaign.Status,
			CompletedDate:    existingCampaign.CompletedDate,
		},
	})
}

// DeleteCampaign deletes a campaign by ID
func DeleteCampaign(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var campaign models.Campaign
	if err := db.First(&campaign, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Kampanye tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menemukan kampanye: " + err.Error(),
		})
		return
	}

	// Tidak boleh menghapus kampanye yang sedang berjalan/selesai
	if campaign.Status == "in_progress" || campaign.Status == "completed" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Tidak dapat menghapus kampanye yang sedang berjalan atau sudah selesai",
		})
		return
	}

	if err := db.Delete(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus kampanye: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kampanye berhasil dihapus",
		"data":    nil, // Data null jika tidak ada yang dikembalikan
	})
}

func LaunchCampaign(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID kampanye tidak valid",
		})
		return
	}

	// Panggil fungsi inti peluncuran kampanye
	campaign, err := PerformCampaignLaunch(db, uint(id))
	if err != nil {
		// Cek jenis error untuk memberikan respons yang lebih spesifik
		if err == gorm.ErrRecordNotFound || err.Error() == "campaign not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Kampanye tidak ditemukan",
			})
		} else if err.Error() == "campaign not in launchable status" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Kampanye tidak dapat diluncurkan karena statusnya tidak 'draft' atau 'stopped'",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal meluncurkan kampanye: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kampanye berhasil diluncurkan",
		"data": models.CampaignResponse{
			ID:               campaign.ID,
			Name:             campaign.Name,
			LaunchDate:       campaign.LaunchDate,
			SendEmailBy:      campaign.SendEmailBy,
			GroupID:          campaign.GroupID,
			EmailTemplateID:  campaign.EmailTemplateID,
			LandingPageID:    campaign.LandingPageID,
			SendingProfileID: campaign.SendingProfileID,
			URL:              campaign.URL,
			CreatedBy:        campaign.CreatedBy,
			CreatedAt:        campaign.CreatedAt,
			UpdatedAt:        campaign.UpdatedAt,
			Status:           campaign.Status,
			CompletedDate:    campaign.CompletedDate,
		},
	})
}

// PerformCampaignLaunch adalah logika inti untuk meluncurkan kampanye, dapat dipanggil secara internal atau dari HTTP handler
func PerformCampaignLaunch(db *gorm.DB, campaignID uint) (*models.Campaign, error) {
	var campaign models.Campaign
	if err := db.First(&campaign, campaignID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("campaign not found")
		}
		return nil, fmt.Errorf("failed to find campaign: %w", err)
	}

	// Hanya kampanye dengan status 'draft' atau 'stopped' yang bisa diluncurkan
	if campaign.Status != "draft" && campaign.Status != "stopped" {
		return nil, fmt.Errorf("campaign not in launchable status")
	}

	// Ubah status kampanye menjadi 'in_progress'
	campaign.Status = "in_progress"
	if err := db.Save(&campaign).Error; err != nil {
		return nil, fmt.Errorf("failed to update campaign status to in_progress: %w", err)
	}

	// --- Logika Peluncuran Kampanye Sebenarnya ---
	// Ini adalah bagian yang akan dijalankan di goroutine terpisah
	go StartEmailSendingProcess(db, campaign)

	return &campaign, nil
}

// StartEmailSendingProcess memulai proses pengiriman email untuk kampanye yang diberikan
func StartEmailSendingProcess(db *gorm.DB, campaign models.Campaign) {
	// Ambil detail Email Template, Landing Page, Sending Profile, Group dan Members
	var emailTemplate models.EmailTemplate
	if err := db.First(&emailTemplate, campaign.EmailTemplateID).Error; err != nil {
		fmt.Printf("Error fetching email template %d for campaign %d: %v\n", campaign.EmailTemplateID, campaign.ID, err)
		return // Hentikan proses jika template tidak ditemukan
	}

	var landingPage models.LandingPage
	if err := db.First(&landingPage, campaign.LandingPageID).Error; err != nil {
		fmt.Printf("Error fetching landing page %d for campaign %d: %v\n", campaign.LandingPageID, campaign.ID, err)
		return // Hentikan proses jika landing page tidak ditemukan
	}

	var sendingProfile models.SendingProfiles
	if err := db.Preload("EmailHeaders").First(&sendingProfile, campaign.SendingProfileID).Error; err != nil {
		fmt.Printf("Error fetching sending profile %d for campaign %d: %v\n", campaign.SendingProfileID, campaign.ID, err)
		return // Hentikan proses jika sending profile tidak ditemukan
	}

	var group models.Group
	if err := db.Preload("Members").First(&group, campaign.GroupID).Error; err != nil {
		fmt.Printf("Error fetching group %d for campaign %d: %v\n", campaign.GroupID, campaign.ID, err)
		return // Hentikan proses jika group tidak ditemukan
	}

	fmt.Printf("Meluncurkan kampanye '%s' untuk %d anggota grup '%s'...\n", campaign.Name, len(group.Members), group.Name)
	for _, member := range group.Members {
		emailBody := emailTemplate.Body
		processedEmailBody := services.ProcessEmailBody(emailBody, member, campaign.URL, landingPage.Body)

		err := services.SendEmail(sendingProfile, member.Email, emailTemplate.Subject, processedEmailBody, landingPage.Body)
		if err != nil {
			fmt.Printf("Gagal mengirim email ke %s: %v\n", member.Email, err)
			// TODO: Log kegagalan, update status member di database jika diperlukan
		} else {
			fmt.Printf("Email berhasil dikirim ke %s\n", member.Email)
			// TODO: Update status member di database (misal: "email_sent")
		}

		// Tambahkan delay jika SendEmailBy ditentukan untuk distribusi email
		if campaign.SendEmailBy != nil && campaign.LaunchDate.Before(*campaign.SendEmailBy) && len(group.Members) > 1 {
			duration := campaign.SendEmailBy.Sub(campaign.LaunchDate)
			if duration > 0 {
				delayPerEmail := duration / time.Duration(len(group.Members))
				time.Sleep(delayPerEmail)
			}
		}
	}
	fmt.Printf("Proses pengiriman email untuk kampanye '%s' selesai.\n", campaign.Name)

	// Opsi: Setelah semua email dikirim, secara otomatis tandai kampanye sebagai 'completed'
	// campaign.Status = "completed"
	// now := time.Now()
	// campaign.CompletedDate = &now
	// if err := db.Save(&campaign).Error; err != nil {
	//     log.Printf("Failed to mark campaign %d as completed after sending: %v", campaign.ID, err)
	// }
}

// CompleteCampaign handles marking a campaign as completed
func CompleteCampaign(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var campaign models.Campaign
	if err := db.First(&campaign, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Kampanye tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menemukan kampanye: " + err.Error(),
		})
		return
	}

	if campaign.Status == "completed" {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Kampanye sudah selesai",
			"data": models.CampaignResponse{
				ID:               campaign.ID,
				Name:             campaign.Name,
				LaunchDate:       campaign.LaunchDate,
				SendEmailBy:      campaign.SendEmailBy,
				GroupID:          campaign.GroupID,
				EmailTemplateID:  campaign.EmailTemplateID,
				LandingPageID:    campaign.LandingPageID,
				SendingProfileID: campaign.SendingProfileID,
				URL:              campaign.URL,
				CreatedBy:        campaign.CreatedBy,
				CreatedAt:        campaign.CreatedAt,
				UpdatedAt:        campaign.UpdatedAt,
				Status:           campaign.Status,
				CompletedDate:    campaign.CompletedDate,
			},
		})
		return
	}

	campaign.Status = "completed"
	now := time.Now()
	campaign.CompletedDate = &now
	if err := db.Save(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyelesaikan kampanye: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kampanye berhasil diselesaikan",
		"data": models.CampaignResponse{
			ID:               campaign.ID,
			Name:             campaign.Name,
			LaunchDate:       campaign.LaunchDate,
			SendEmailBy:      campaign.SendEmailBy,
			GroupID:          campaign.GroupID,
			EmailTemplateID:  campaign.EmailTemplateID,
			LandingPageID:    campaign.LandingPageID,
			SendingProfileID: campaign.SendingProfileID,
			URL:              campaign.URL,
			CreatedBy:        campaign.CreatedBy,
			CreatedAt:        campaign.CreatedAt,
			UpdatedAt:        campaign.UpdatedAt,
			Status:           campaign.Status,
			CompletedDate:    campaign.CompletedDate,
		},
	})
}

func CheckAndLaunchCampaigns(db *gorm.DB) {
	var campaigns []models.Campaign
	// Cari kampanye yang statusnya 'draft' atau 'stopped' (jika ingin bisa dijadwalkan ulang setelah dihentikan)
	// DAN launch_date-nya sudah tiba atau terlewati
	result := db.Where("(status = ? OR status = ?) AND launch_date <= ?", "draft", "stopped", time.Now()).Find(&campaigns)
	if result.Error != nil {
		// Jangan gunakan log.Fatal di sini karena akan menghentikan aplikasi. Gunakan log.Printf.
		log.Printf("Error checking campaigns for launch: %v", result.Error)
		return
	}

	if len(campaigns) > 0 {
		log.Printf("Found %d campaign(s) ready for launch.", len(campaigns))
	}

	for _, campaign := range campaigns {
		log.Printf("Attempting to launch campaign '%s' (ID: %d) with LaunchDate: %s...", campaign.Name, campaign.ID, campaign.LaunchDate.Format(time.RFC3339))

		// Panggil logika inti peluncuran kampanye yang sudah di-refactor
		_, err := PerformCampaignLaunch(db, uint(campaign.ID))
		if err != nil {
			log.Printf("Failed to launch campaign %d ('%s') via scheduler: %v", campaign.ID, campaign.Name, err)
			// Anda mungkin ingin mencatat error ini ke activity logs atau sistem monitoring
		} else {
			log.Printf("Successfully launched campaign %d ('%s') via scheduler.", campaign.ID, campaign.Name)
		}
	}
}
