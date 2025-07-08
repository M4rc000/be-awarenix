package config

import (
	"log"
	"time"

	"be-awarenix/models"
	"be-awarenix/services"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	usersToSeed := []models.User{
		{
			Name:         "Marco Antonio",
			Email:        "marco@gmail.com",
			PasswordHash: "123456",
			Position:     "Security Consultant",
			IsActive:     1,
			Role:         2,
		},
		{
			Name:         "Jane Doe",
			Email:        "jane.doe@example.com",
			PasswordHash: "password_jane",
			Position:     "Developer",
			IsActive:     1,
			Role:         3,
		},
	}

	// 2. Lakukan perulangan untuk setiap user dalam slice di atas.
	for _, userData := range usersToSeed {
		// 3. Cek apakah user dengan email tersebut sudah ada di database.
		var existingUser models.User
		err := db.Where("email = ?", userData.Email).First(&existingUser).Error

		// Jika user tidak ditemukan (ini yang kita harapkan), maka buat user baru.
		if err != nil && err == gorm.ErrRecordNotFound {
			log.Printf("Seeding user '%s'...", userData.Name)

			// Hash password mentah dari data slice
			hashedPassword, hashErr := services.HashPassword(userData.PasswordHash)
			if hashErr != nil {
				log.Fatalf("Failed to hash password for seeder: %v", hashErr)
			}
			// Update field PasswordHash dengan hasil hash
			userData.PasswordHash = hashedPassword

			// Buat user di database
			if err := db.Create(&userData).Error; err != nil {
				log.Fatalf("Failed to seed user '%s': %v", userData.Name, err)
			}
			log.Printf("User '%s' seeded successfully.", userData.Name)

		} else if err != nil {
			// Handle error lain selain record not found
			log.Fatalf("Error checking for user '%s': %v", userData.Name, err)
		} else {
			// Jika tidak ada error, berarti user sudah ada
			log.Printf("User with email '%s' already exists. Seeder skipped.", userData.Email)
		}
	}
}

func SeedRoles(db *gorm.DB) {
	rolesToSeed := []models.Role{
		{
			Name:      "Super Admin",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			Name:      "Admin",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			Name:      "Engineer",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
	}

	for _, roleData := range rolesToSeed {
		var existingRole models.Role
		err := db.Where("name = ?", roleData.Name).First(&existingRole).Error

		if err != nil && err == gorm.ErrRecordNotFound {
			log.Printf("Seeding role '%s'...", roleData.Name)

			// Buat role di database
			if err := db.Create(&roleData).Error; err != nil {
				log.Fatalf("Failed to seed role '%s': %v", roleData.Name, err)
			}
			log.Printf("Role '%s' seeded successfully.", roleData.Name)

		} else if err != nil {
			log.Fatalf("Error checking for role '%s': %v", roleData.Name, err)
		} else {
			log.Printf("Role with name '%s' already exists. Seeder skipped.", roleData.Name)
		}
	}
}

func SeedMenus(db *gorm.DB) {
	rolesToMenu := []models.Menu{
		{
			Name:      "Admin",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			Name:      "Phishing Simulation",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			Name:      "User",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			Name:      "Logging & Monitoring",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
	}

	for _, menuData := range rolesToMenu {
		var existingRole models.Menu
		err := db.Where("name = ?", menuData.Name).First(&existingRole).Error

		if err != nil && err == gorm.ErrRecordNotFound {
			log.Printf("Seeding role '%s'...", menuData.Name)

			// Buat role di database
			if err := db.Create(&menuData).Error; err != nil {
				log.Fatalf("Failed to seed role '%s': %v", menuData.Name, err)
			}
			log.Printf("Role '%s' seeded successfully.", menuData.Name)

		} else if err != nil {
			log.Fatalf("Error checking for role '%s': %v", menuData.Name, err)
		} else {
			log.Printf("Role with name '%s' already exists. Seeder skipped.", menuData.Name)
		}
	}
}

func SeedSubmenus(db *gorm.DB) {
	rolesToMenu := []models.Submenu{
		{
			MenuID:    1,
			Name:      "Dashboard",
			Icon:      "GridIcon",
			Url:       "/admin/dashboard",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    1,
			Name:      "Campaign",
			Icon:      "CalenderIcon",
			Url:       "/admin/campaign",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    1,
			Name:      "Role Management",
			Icon:      "FaUserCog",
			Url:       "/admin/role-management",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    1,
			Name:      "Role Management",
			Icon:      "FaUserCog",
			Url:       "/admin/role-management",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    1,
			Name:      "User Management",
			Icon:      "User",
			Url:       "/admin/user-management",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    1,
			Name:      "Groups & Members",
			Icon:      "GroupIcon",
			Url:       "/admin/groups-members",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    1,
			Name:      "Email Templates",
			Icon:      "MailIcon",
			Url:       "/admin/email-templates",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    1,
			Name:      "Landing Pages",
			Icon:      "TableIcon",
			Url:       "/admin/landing-pages",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    1,
			Name:      "Sending Profiles",
			Icon:      "UserIcon",
			Url:       "/admin/sending-profiles",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    2,
			Name:      "Phishing Emails",
			Icon:      "MdOutlineAttachEmail",
			Url:       "/phishing-simulation/phishing-emails",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    2,
			Name:      "Phishing Websites",
			Icon:      "CgWebsite",
			Url:       "/phishing-simulation/phishing-websites",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    2,
			Name:      "Training Modules",
			Icon:      "IoIosBookmarks",
			Url:       "/phishing-simulation/training-modules",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    3,
			Name:      "Account Settings",
			Icon:      "IoSettingsOutline",
			Url:       "/user/account-settings",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    3,
			Name:      "Subscription & Billing",
			Icon:      "FaMoneyCheckAlt",
			Url:       "/user/subscription-billing",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    4,
			Name:      "Logging Activity",
			Icon:      "IoFootstepsOutline",
			Url:       "/logging-monitoring/logging-activity",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
		{
			MenuID:    4,
			Name:      "Environtment Check",
			Icon:      "DiEnvato",
			Url:       "/logging-monitoring/environment-check",
			IsActive:  1,
			CreatedAt: time.Now(),
			CreatedBy: 0,
			UpdatedAt: time.Now(),
			UpdatedBy: 0,
		},
	}

	for _, submenuData := range rolesToMenu {
		var existingRole models.Submenu
		err := db.Where("name = ?", submenuData.Name).First(&existingRole).Error

		if err != nil && err == gorm.ErrRecordNotFound {
			log.Printf("Seeding role '%s'...", submenuData.Name)

			// Buat role di database
			if err := db.Create(&submenuData).Error; err != nil {
				log.Fatalf("Failed to seed role '%s': %v", submenuData.Name, err)
			}
			log.Printf("Role '%s' seeded successfully.", submenuData.Name)

		} else if err != nil {
			log.Fatalf("Error checking for role '%s': %v", submenuData.Name, err)
		} else {
			log.Printf("Role with name '%s' already exists. Seeder skipped.", submenuData.Name)
		}
	}
}

func SeedEmailTemplates(db *gorm.DB) {
	emailTemplates := []models.EmailTemplate{
		{
			Name:           "Google Meet Invitation",
			EnvelopeSender: "noreply@yourdomain.com",
			Subject:        "Invitation: General Meeting @ {{.MeetingTime}}",
			Body: `<!DOCTYPE html>
				<html lang="en">
				<head>
					<meta charset="UTF-8">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
					<title>Google Calendar Invitation</title>
					<style>
						body {
							font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
							background-color: #f5f5f5;
							margin: 0;
							padding: 20px;
							font-size: 14px;
							line-height: 1.4;
						}

						.container {
							max-width: 600px;
							margin: 0 auto;
							background-color: white;
							padding: 24px;
							border-radius: 8px;
							box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
						}

						.section {
							margin-bottom: 24px;
						}

						.section-title {
							font-weight: 600;
							color: #333;
							margin-bottom: 8px;
							font-size: 14px;
						}

						.when-info {
							color: #666;
							font-size: 13px;
							line-height: 1.5;
						}

						.guest-info {
							margin-bottom: 4px;
						}

						.guest-name {
							color: #333;
							font-weight: 500;
						}

						.guest-role {
							color: #666;
							font-size: 13px;
						}

						.guest-email {
							color: #666;
							font-size: 13px;
						}

						.view-all-link {
							color: #1a73e8;
							text-decoration: none;
							font-size: 13px;
							cursor: pointer;
						}

						.view-all-link:hover {
							text-decoration: underline;
						}

						.reply-section {
							margin-top: 16px;
						}

						.reply-text {
							color: #666;
							font-size: 13px;
							margin-bottom: 12px;
						}

						.reply-buttons {
							display: flex;
							gap: 8px;
							flex-wrap: wrap;
						}

						.reply-button {
							padding: 8px 16px;
							border: 1px solid #dadce0;
							background-color: #f8f9fa;
							color: #3c4043;
							border-radius: 4px;
							font-size: 13px;
							cursor: pointer;
							transition: all 0.2s ease;
						}

						.reply-button:hover {
							background-color: #e8eaed;
							border-color: #c4c7ca;
						}

						.join-button {
							background-color: #1a73e8;
							color: white;
							border: none;
							padding: 12px 24px;
							border-radius: 4px;
							font-size: 14px;
							font-weight: 500;
							cursor: pointer;
							float: right;
							margin-top: -40px;
							transition: background-color 0.2s ease;
						}

						.join-button:hover {
							background-color: #1557b0;
						}

						.meeting-link-section {
							margin-top: 20px;
							clear: both;
						}

						.meeting-link-title {
							color: #666;
							font-size: 13px;
							margin-bottom: 4px;
						}

						.meeting-link {
							color: #666;
							font-size: 13px;
						}

						.footer {
							margin-top: 40px;
							padding-top: 20px;
							border-top: 1px solid #e0e0e0;
							font-size: 12px;
							color: #666;
							line-height: 1.5;
						}

						.footer-text {
							margin-bottom: 12px;
						}

						.footer-link {
							color: #1a73e8;
							text-decoration: none;
						}

						.footer-link:hover {
							text-decoration: underline;
						}

						.google-calendar-link {
							color: #1a73e8;
							text-decoration: none;
						}

						.google-calendar-link:hover {
							text-decoration: underline;
						}
					</style>
				</head>
				<body>
					<div class="container">
						<div class="section">
							<div class="section-title">When</div>
							<div class="when-info">Wednesday Jul 9, 2025 · 1pm – 2pm (Western Indonesia Time - Jakarta)</div>
						</div>

						<div class="section">
							<div class="section-title">Guests</div>
							<div class="guest-info">
								<div class="guest-name">Marco Antonio <span class="guest-role">- organizer</span></div>
								<div class="guest-email">marcoantoniomadigaskar90@gmail.com</div>
							</div>
							<a href="#" class="view-all-link">View all guest info</a>
						</div>

						<div class="reply-section">
							<div class="reply-text">Reply for marcoantoniomadigaskar90@gmail.com</div>
							<div class="reply-buttons">
								<button class="reply-button">Yes</button>
								<button class="reply-button">No</button>
								<button class="reply-button">Maybe</button>
								<button class="reply-button">More options</button>
							</div>
						</div>

						<button class="join-button">Join with Google Meet</button>

						<div class="meeting-link-section">
							<div class="meeting-link-title">Meeting link</div>
							<div class="meeting-link">meet.google.com/znd-xjey-isd</div>
						</div>

						<div class="footer">
							<div class="footer-text">
								Invitation from <a href="#" class="google-calendar-link">Google Calendar</a>
							</div>
							<div class="footer-text">
								You are receiving this email because you are subscribed to calendar notifications. To stop receiving these emails, go to <a href="#" class="footer-link">Calendar settings</a>, select this calendar, and change "Other notifications".
							</div>
							<div class="footer-text">
								Forwarding this invitation could allow any recipient to send a response to the organizer, be added to the guest list, invite others regardless of their own invitation status, or modify your RSVP. <a href="#" class="footer-link">Learn more</a>
							</div>
						</div>
					</div>
				</body>
				</html>`,
			TrackerImage: 1,
			CreatedAt:    time.Now(),
			CreatedBy:    0,
		},
	}

	for _, emailTemplateData := range emailTemplates {
		var existingEmailTemplate models.EmailTemplate
		err := db.Where("name = ?", emailTemplateData.Name).First(&existingEmailTemplate).Error

		if err != nil && err == gorm.ErrRecordNotFound {
			log.Printf("Seeding email template '%s'...", emailTemplateData.Name)

			// Buat email template di database
			if err := db.Create(&emailTemplateData).Error; err != nil {
				log.Fatalf("Failed to seed email template '%s': %v", emailTemplateData.Name, err)
			}
			log.Printf("Email template '%s' seeded successfully.", emailTemplateData.Name)

		} else if err != nil {
			log.Fatalf("Error checking for email template '%s': %v", emailTemplateData.Name, err)
		} else {
			log.Printf("Email template with name '%s' already exists. Seeder skipped.", emailTemplateData.Name)
		}
	}
}
