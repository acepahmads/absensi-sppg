package service

import (
	"absensi-sppg/internal/repository"
	"absensi-sppg/pkg/utils"
	"context"
	"errors"
	"fmt"
	"os"

	"absensi-sppg/internal/model"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return UserService{
		repo: repo,
	}
}

func (s *UserService) RegisterUser(userp model.RegisterAccount) (string, error) {
	// Check if email already exists
	existingUser, err := s.repo.FindByEmail(userp.Email)
	if err != nil {
		fmt.Println("Error finding user by email:", err)
		return "", err
	}
	if existingUser != "" {
		fmt.Println("Email already exists")
		return "", errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(userp.Password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return "", err
	}

	// Create the user
	userp.Password = hashedPassword
	return s.repo.Create(&userp)
}

func (s *UserService) GetUserInfoByID(id string) (*model.UserAccount, error) {
	return s.repo.FindUserInfoByID(id)
}

func (s *UserService) Registered(tenantID int) ([]string, error) {
	return s.repo.Registered(tenantID)
}

func (s *UserService) GetLeaders(tenantID int) ([]string, error) {
	return s.repo.GetLeaders(tenantID)
}

func (s *UserService) GetTenants(ctx context.Context) ([]model.Tenant, error) {
	return s.repo.GetTenants(ctx)
}

func (s *UserService) GetAllUserKaryawan(ctx context.Context) ([]model.UserKaryawan, error) {
	return s.repo.GetAllUserKaryawan(ctx)
}

func (s *UserService) CreateUserKaryawan(ctx context.Context, uk *model.UserKaryawan) error {
	return s.repo.CreateUserKaryawan(ctx, uk)
}

func (s *UserService) UpdateUserKaryawan(ctx context.Context, uk *model.UserKaryawan) error {
	return s.repo.UpdateUserKaryawan(ctx, uk)
}

func (s *UserService) DeleteUserKaryawan(ctx context.Context, id int) error {
	return s.repo.DeleteUserKaryawan(ctx, id)
}

func (s *UserService) GetLeadersList(ctx context.Context) ([]model.KaryawanLeader, error) {
	return s.repo.GetLeadersList(ctx)
}

// Leader CRUD
func (s *UserService) GetAllLeaders(ctx context.Context) ([]model.KaryawanLeader, error) {
	return s.repo.GetAllLeaders(ctx)
}

func (s *UserService) CreateLeader(ctx context.Context, leader *model.KaryawanLeader) error {
	return s.repo.CreateLeader(ctx, leader)
}

func (s *UserService) UpdateLeader(ctx context.Context, leader *model.KaryawanLeader) error {
	return s.repo.UpdateLeader(ctx, leader)
}

func (s *UserService) DeleteLeader(ctx context.Context, id int) error {
	return s.repo.DeleteLeader(ctx, id)
}

// User Account CRUD
func (s *UserService) GetAllUserAccounts(ctx context.Context) ([]model.UserAccountCRUD, error) {
	return s.repo.GetAllUserAccounts(ctx)
}

func (s *UserService) CreateUserAccount(ctx context.Context, ua *model.UserAccountCRUD) error {
	if ua.Password != "" {
		hashed, err := utils.HashPassword(ua.Password)
		if err != nil {
			return err
		}
		ua.Password = hashed
	}
	return s.repo.CreateUserAccount(ctx, ua)
}

func (s *UserService) UpdateUserAccount(ctx context.Context, ua *model.UserAccountCRUD) error {
	if ua.Password != "" {
		hashed, err := utils.HashPassword(ua.Password)
		if err != nil {
			return err
		}
		ua.Password = hashed
	}
	return s.repo.UpdateUserAccount(ctx, ua)
}

func (s *UserService) DeleteUserAccount(ctx context.Context, id string) error {
	return s.repo.DeleteUserAccount(ctx, id)
}

func (s *UserService) RegisterTenant(req model.RegisterTenantRequest) error {
	hashed, err := utils.HashPassword(req.AdminPassword)
	if err != nil {
		return err
	}
	req.AdminPassword = hashed
	return s.repo.CreateTenant(&req)
}

func (s *UserService) ForgotPasswordRequest(email string) error {
	existing, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}
	if existing == "" {
		return errors.New("email tidak terdaftar")
	}

	// Resolve the base APP URL
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		port := os.Getenv("APP_PORT")
		if port == "" {
			port = "8080"
		}
		appURL = fmt.Sprintf("http://localhost:%s", port)
	}
	resetLink := fmt.Sprintf("%s/reset_password?email=%s", appURL, email)

	// SMTP settings
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpSender := os.Getenv("SMTP_SENDER")

	// If SMTP is not properly configured, log to console as simulation fallback
	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || smtpUser == "your_email@gmail.com" {
		fmt.Printf("[WARN] SMTP is not configured or uses default values. Falling back to simulation link.\n")
		fmt.Printf("[EMAIL SIMULATION] Reset link sent to %s: %s\n", email, resetLink)
		return nil
	}

	// Construct HTML email body
	subject := "Reset Kata Sandi Akun Absensi SPPG Anda"
	htmlBody := fmt.Sprintf(`
		<div style="font-family: 'Inter', Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 30px; border: 1px solid #e2e8f0; border-radius: 12px; background-color: #ffffff; color: #1e293b;">
			<div style="text-align: center; margin-bottom: 24px;">
				<img src="https://absen.cbinstrument.com/static/LOGO-CBI.png" alt="Logo CBI" style="max-height: 60px; object-fit: contain;">
				<h2 style="color: #0f172a; font-size: 22px; margin-top: 15px; font-weight: 700;">Pemulihan Kata Sandi</h2>
			</div>
			<p style="font-size: 15px; line-height: 1.6; color: #475569;">Halo,</p>
			<p style="font-size: 15px; line-height: 1.6; color: #475569;">Kami menerima permintaan untuk mengatur ulang kata sandi akun Absensi SPPG Anda. Silakan klik tombol di bawah ini untuk mereset kata sandi Anda:</p>
			
			<div style="text-align: center; margin: 30px 0;">
				<a href="%s" style="display: inline-block; padding: 14px 28px; background: linear-gradient(135deg, #0052cc 0%%, #0066ff 100%%); color: #ffffff; text-decoration: none; border-radius: 9999px; font-weight: 600; font-size: 15px; box-shadow: 0 4px 15px rgba(0, 102, 255, 0.25);">Atur Ulang Kata Sandi</a>
			</div>
			
			<p style="font-size: 14px; line-height: 1.6; color: #64748b;">Tautan ini berlaku selama 1 jam. Jika Anda tidak merasa melakukan permintaan ini, silakan abaikan email ini secara aman.</p>
			<hr style="border: none; border-top: 1px solid #e2e8f0; margin: 24px 0;">
			<p style="font-size: 12px; text-align: center; color: #94a3b8;">&copy; 2026 CBI Absensi SPPG. All rights reserved.</p>
		</div>
	`, resetLink)

	// Send the real email
	err = utils.SendEmail(smtpHost, smtpPort, smtpUser, smtpPass, smtpSender, email, subject, htmlBody)
	if err != nil {
		fmt.Printf("[ERROR] Failed to send real reset email: %v\n", err)
		return errors.New("gagal mengirim email pemulihan, silakan hubungi admin")
	}

	fmt.Printf("[EMAIL SUCCESS] Reset email successfully sent to %s\n", email)
	return nil
}

func (s *UserService) ResetPassword(email, newPassword string) error {
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.repo.ResetPassword(email, hashed)
}

func (s *UserService) GetTenantIDByDeviceSN(ctx context.Context, sn string) (int, error) {
	return s.repo.GetTenantIDByDeviceSN(ctx, sn)
}


