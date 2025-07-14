package services

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"be-awarenix/models"

	"github.com/gophish/gomail"
)

func SendTestEmail(profile *models.SendingProfiles, recipientEmail, body, subject string) error {
	// Pastikan profil dan email tujuan valid
	if profile == nil || profile.SmtpFrom == "" || profile.Host == "" {
		return fmt.Errorf("invalid sending profile configuration")
	}
	if recipientEmail == "" {
		return fmt.Errorf("recipient email cannot be empty")
	}

	// from := profile.SmtpFrom
	from := profile.SmtpFrom
	password := profile.Password
	smtpHost := profile.Host
	smtpPort := "587" // **Menggunakan port 587 untuk STARTTLS**

	// Jika host memiliki port, pisahkan
	if strings.Contains(smtpHost, ":") {
		parts := strings.Split(smtpHost, ":")
		smtpHost = parts[0]
		smtpPort = parts[1]
	}

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// TLS config (tetap diperlukan untuk StartTLS)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false, // Set false untuk keamanan produksi, pastikan sertifikat valid
		ServerName:         smtpHost,
	}

	// --- BAGIAN YANG DIUBAH UNTUK PORT 587 (STARTTLS) ---
	// 1. Dial koneksi TCP biasa terlebih dahulu
	conn, err := net.Dial("tcp", smtpHost+":"+smtpPort)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}

	// 2. Buat klien SMTP dari koneksi TCP
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close() // Pastikan klien ditutup

	// 3. Lakukan STARTTLS untuk meng-upgrade koneksi ke TLS
	if err = client.StartTLS(tlsconfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}
	// --- AKHIR BAGIAN YANG DIUBAH ---

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate with SMTP server: %w", err)
	}

	// Setup message
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + recipientEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			mime + "\r\n" +
			body,
	)

	// Send email
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(recipientEmail); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	client.Quit() // Penting untuk mengakhiri sesi SMTP

	return nil
}

func SendEmail(profile models.SendingProfiles, toEmail, subject, body, landingPageBody string) error {
	// Untuk demo, kita akan print saja.
	// Di produksi, Anda akan menggunakan gomail atau net/smtp.

	// Contoh menggunakan gomail
	m := gomail.NewMessage()
	m.SetHeader("From", profile.SmtpFrom)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body) // Asumsi body adalah HTML

	// Tambahkan custom headers jika ada
	for _, header := range profile.EmailHeaders {
		m.SetHeader(header.Header, header.Value)
	}

	d := gomail.NewDialer(profile.Host, 587, profile.Username, profile.Password) // Port dan autentikasi
	// Anda mungkin perlu menyesuaikan port, TLSConfig, dll. tergantung pada SMTP server
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Hanya untuk testing

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	fmt.Printf("Email terkirim ke %s dari profil %s\n", toEmail, profile.Name)
	return nil
}

// ProcessEmailBody mengganti placeholder dalam body email
func ProcessEmailBody(emailBody string, member models.Member, campaignURL string, landingPageBody string) string {
	// Contoh sederhana penggantian placeholder.
	// Untuk yang lebih canggih, Anda bisa menggunakan text/template atau html/template.
	processedBody := strings.ReplaceAll(emailBody, "{{.Recipient.Name}}", member.Name)
	processedBody = strings.ReplaceAll(processedBody, "{{.Recipient.Email}}", member.Email)
	processedBody = strings.ReplaceAll(processedBody, "{{.Recipient.Position}}", member.Position)
	processedBody = strings.ReplaceAll(processedBody, "{{.Recipient.Company}}", member.Company)
	processedBody = strings.ReplaceAll(processedBody, "{{.Recipient.Country}}", member.Country)
	processedBody = strings.ReplaceAll(processedBody, "{{.URL}}", campaignURL)

	// Jika ada kebutuhan untuk menyuntikkan bagian dari landingPageBody ke email (jarang, tapi mungkin)
	processedBody = strings.ReplaceAll(processedBody, "{{.LandingPageBody}}", landingPageBody)

	// Placeholder untuk tracking image (1x1 pixel gif)
	// Anda akan membuat endpoint di router yang mengembalikan 1x1 gif dan mencatat kunjungan
	trackingPixelURL := fmt.Sprintf("%s/track/%d/%s", campaignURL, member.GroupID, member.Email) // Contoh URL tracking
	processedBody = strings.ReplaceAll(processedBody, "{{.Tracker}}", fmt.Sprintf("<img src=\"%s\" style=\"display:none;\"/>", trackingPixelURL))

	return processedBody
}
