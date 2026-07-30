package helpers

import (
	"bytes"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/dto"
	"html/template"
	"net/smtp"
	"strconv"
)

type AuthHelper struct {
	userHelper  *UserHelper
	tokenHelper *TokenHelper
	appConfig   *config.AppConfig
	mailConfig  *config.EmailConfig
}

func NewAuthHelper(userHelper *UserHelper, tokenHelper *TokenHelper, appConfig *config.AppConfig, mailConfig *config.EmailConfig) *AuthHelper {
	return &AuthHelper{
		userHelper:  userHelper,
		tokenHelper: tokenHelper,
		appConfig:   appConfig,
		mailConfig:  mailConfig,
	}
}

func (h *AuthHelper) SendEmail(payload *dto.SendMailMetaData) error {
	switch payload.Type {
	case dto.ResetPassword:
		return h.SendResetPasswordEmail(payload)
	case dto.EmailVerification:
		return h.SendVerificationEmail(payload)
	case dto.RegistrationInfo:
		return h.SendRegistrationInfo(payload)
	default:
		return fmt.Errorf("unknown email type: %s", payload.Type)
	}
}

func (h *AuthHelper) SendVerificationEmail(payload *dto.SendMailMetaData) error {
	url := fmt.Sprintf("https://%s/confirm?token=%s", h.appConfig.ClientOrigin, payload.Token)

	data := struct {
		VerificationURL string
		FullName        string
		Year            int
	}{
		VerificationURL: url,
		FullName:        payload.User.FullName,
		Year:            2025,
	}

	subject := "Verify Your Email Address"
	body, err := h.renderTemplate(emailVerificationTemplate, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return h.sendHTMLEmail(payload.To, subject, body)
}

func (h *AuthHelper) SendResetPasswordEmail(payload *dto.SendMailMetaData) error {
	url := fmt.Sprintf("https://%s/reset-password?token=%s", h.appConfig.ClientOrigin, payload.Token)

	data := struct {
		ResetURL string
		FullName string
		Year     int
	}{
		ResetURL: url,
		FullName: payload.User.FullName,
		Year:     2025,
	}

	subject := "Reset Your Password"
	body, err := h.renderTemplate(resetPasswordTemplate, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return h.sendHTMLEmail(payload.To, subject, body)
}

func (h *AuthHelper) SendRegistrationInfo(payload *dto.SendMailMetaData) error {
	data := struct {
		Email    string
		Password string
		FullName string
		LoginURL string
		Year     int
	}{
		Email:    payload.To,
		Password: payload.Password,
		FullName: payload.User.FullName,
		LoginURL: fmt.Sprintf("https://%s/signin", h.appConfig.ClientOrigin),
		Year:     2025,
	}

	subject := "Welcome to AGC Forge - Account Information"
	body, err := h.renderTemplate(registrationInfoTemplate, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return h.sendHTMLEmail(payload.To, subject, body)
}

func (h *AuthHelper) renderTemplate(tmpl string, data interface{}) (string, error) {
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (h *AuthHelper) sendHTMLEmail(to, subject, htmlBody string) error {
	message := []byte("Subject: " + subject + "\r\n" +
		"To: " + to + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		htmlBody + "\r\n")

	auth := smtp.PlainAuth(
		"",
		h.mailConfig.SMTPUsername,
		h.mailConfig.SMTPPassword,
		h.mailConfig.SMTPHost,
	)

	addr := fmt.Sprintf("%s:%s", h.mailConfig.SMTPHost, h.mailConfig.SMTPPort)
	err := smtp.SendMail(addr, auth, h.mailConfig.SMTPUsername, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (h *AuthHelper) ValidateSMTPConfig() error {
	if h.mailConfig == nil {
		return fmt.Errorf("SMTP configuration is nil")
	}

	requiredFields := []struct {
		value *string
		field string
	}{
		{&h.mailConfig.SMTPHost, "SMTP host"},
		{&h.mailConfig.SMTPUsername, "SMTP username"},
		{&h.mailConfig.SMTPPassword, "SMTP password"},
		{&h.mailConfig.SMTPPort, "SMTP port"},
	}

	for _, field := range requiredFields {
		if field.value == nil || *field.value == "" {
			return fmt.Errorf("%s is not configured", field.field)
		}
	}

	if _, err := strconv.Atoi(h.mailConfig.SMTPPort); err != nil {
		return fmt.Errorf("invalid SMTP port format: %v", err)
	}

	return nil
}

// Email Verification Template
const emailVerificationTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f7fa;">
    <table role="presentation" style="width: 100%; border-collapse: collapse; background-color: #f5f7fa;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 30px; text-align: center; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">Verify Your Email</h1>
                        </td>
                    </tr>
                    
                    <!-- Body -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                Hi <strong>{{.FullName}}</strong>,
                            </p>
                            <p style="margin: 0 0 30px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                Thank you for signing up! To complete your registration and start using your account, please verify your email address by clicking the button below.
                            </p>
                            
                            <!-- CTA Button -->
                            <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                <tr>
                                    <td align="center" style="padding: 20px 0;">
                                        <a href="{{.VerificationURL}}" style="display: inline-block; padding: 16px 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 4px 6px rgba(102, 126, 234, 0.4);">
                                            Verify Email Address
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="margin: 30px 0 20px; color: #718096; font-size: 14px; line-height: 1.6;">
                                Or copy and paste this link into your browser:
                            </p>
                            <p style="margin: 0; padding: 15px; background-color: #f7fafc; border-radius: 6px; word-break: break-all; font-size: 13px; color: #4a5568; border-left: 4px solid #667eea;">
                                {{.VerificationURL}}
                            </p>
                            
                            <p style="margin: 30px 0 0; color: #718096; font-size: 14px; line-height: 1.6;">
                                This link will expire in 24 hours. If you didn't create an account, please ignore this email.
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px; background-color: #f7fafc; border-radius: 0 0 12px 12px; text-align: center;">
                            <p style="margin: 0 0 10px; color: #a0aec0; font-size: 13px;">
                                © {{.Year}} AGC Forge. All rights reserved.
                            </p>
                            <p style="margin: 0; color: #a0aec0; font-size: 13px;">
                                Need help? Contact us at support@socialforge.io
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`

// Reset Password Template
const resetPasswordTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Password</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f7fa;">
    <table role="presentation" style="width: 100%; border-collapse: collapse; background-color: #f5f7fa;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 30px; text-align: center; background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">Reset Your Password</h1>
                        </td>
                    </tr>
                    
                    <!-- Body -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                Hi <strong>{{.FullName}}</strong>,
                            </p>
                            <p style="margin: 0 0 30px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                We received a request to reset your password. Click the button below to create a new password for your account.
                            </p>
                            
                            <!-- CTA Button -->
                            <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                <tr>
                                    <td align="center" style="padding: 20px 0;">
                                        <a href="{{.ResetURL}}" style="display: inline-block; padding: 16px 40px; background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 4px 6px rgba(245, 87, 108, 0.4);">
                                            Reset Password
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="margin: 30px 0 20px; color: #718096; font-size: 14px; line-height: 1.6;">
                                Or copy and paste this link into your browser:
                            </p>
                            <p style="margin: 0; padding: 15px; background-color: #f7fafc; border-radius: 6px; word-break: break-all; font-size: 13px; color: #4a5568; border-left: 4px solid #f5576c;">
                                {{.ResetURL}}
                            </p>
                            
                            <!-- Security Notice -->
                            <div style="margin-top: 30px; padding: 20px; background-color: #fff5f5; border-left: 4px solid #f5576c; border-radius: 6px;">
                                <p style="margin: 0 0 10px; color: #c53030; font-weight: 600; font-size: 14px;">
                                    🔒 Security Notice
                                </p>
                                <p style="margin: 0; color: #742a2a; font-size: 14px; line-height: 1.6;">
                                    This link will expire in 1 hour. If you didn't request a password reset, please ignore this email and your password will remain unchanged.
                                </p>
                            </div>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px; background-color: #f7fafc; border-radius: 0 0 12px 12px; text-align: center;">
                            <p style="margin: 0 0 10px; color: #a0aec0; font-size: 13px;">
                                © {{.Year}} AGC Forge. All rights reserved.
                            </p>
                            <p style="margin: 0; color: #a0aec0; font-size: 13px;">
                                Need help? Contact us at support@socialforge.io
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`

// Registration Info Template
const registrationInfoTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to AGC Forge</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f7fa;">
    <table role="presentation" style="width: 100%; border-collapse: collapse; background-color: #f5f7fa;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 30px; text-align: center; background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0 0 10px; color: #ffffff; font-size: 32px; font-weight: 600;">🎉 Welcome!</h1>
                            <p style="margin: 0; color: #ffffff; font-size: 16px; opacity: 0.95;">Your account has been created</p>
                        </td>
                    </tr>
                    
                    <!-- Body -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                Hi <strong>{{.FullName}}</strong>,
                            </p>
                            <p style="margin: 0 0 30px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                Welcome to <strong>AGC Forge</strong>! Your account has been successfully created. Below are your login credentials:
                            </p>
                            
                            <!-- Credentials Box -->
                            <div style="background: linear-gradient(135deg, #f6f8fb 0%, #e9ecef 100%); border-radius: 8px; padding: 25px; margin-bottom: 30px; border: 1px solid #dee2e6;">
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 12px 0; border-bottom: 1px solid #cbd5e0;">
                                            <p style="margin: 0; color: #718096; font-size: 13px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px;">Email Address</p>
                                            <p style="margin: 8px 0 0; color: #2d3748; font-size: 16px; font-weight: 500;">{{.Email}}</p>
                                        </td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 12px 0;">
                                            <p style="margin: 0; color: #718096; font-size: 13px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px;">Temporary Password</p>
                                            <p style="margin: 8px 0 0; color: #2d3748; font-size: 16px; font-weight: 500; font-family: 'Courier New', monospace; background-color: #ffffff; padding: 10px; border-radius: 4px; display: inline-block;">{{.Password}}</p>
                                        </td>
                                    </tr>
                                </table>
                            </div>
                            
                            <!-- CTA Button -->
                            <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                <tr>
                                    <td align="center" style="padding: 10px 0;">
                                        <a href="{{.LoginURL}}" style="display: inline-block; padding: 16px 40px; background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 4px 6px rgba(79, 172, 254, 0.4);">
                                            Login to Your Account
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <!-- Security Notice -->
                            <div style="margin-top: 30px; padding: 20px; background-color: #fffaf0; border-left: 4px solid #ed8936; border-radius: 6px;">
                                <p style="margin: 0 0 10px; color: #c05621; font-weight: 600; font-size: 14px;">
                                    🔐 Important Security Notice
                                </p>
                                <p style="margin: 0; color: #7c2d12; font-size: 14px; line-height: 1.6;">
                                    This is a temporary password. For your security, please change it immediately after your first login by going to Account Settings > Security.
                                </p>
                            </div>
                            
                            <p style="margin: 30px 0 0; color: #718096; font-size: 14px; line-height: 1.6;">
                                If you have any questions or need assistance, our support team is here to help!
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px; background-color: #f7fafc; border-radius: 0 0 12px 12px; text-align: center;">
                            <p style="margin: 0 0 10px; color: #a0aec0; font-size: 13px;">
                                © {{.Year}} AGC Forge. All rights reserved.
                            </p>
                            <p style="margin: 0; color: #a0aec0; font-size: 13px;">
                                Need help? Contact us at support@socialforge.io
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`
