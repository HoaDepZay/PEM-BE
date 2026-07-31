package email

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

// SendVerificationEmail sends an email with the verification link using gomail
func SendVerificationEmail(toEmail string, verifyLink string) error {
	fromEmail := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASS")

	if fromEmail == "" || password == "" {
		return fmt.Errorf("email credentials are not set in environment variables")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Xác thực tài khoản Visual Finance")

	// Email body (HTML format)
	body := fmt.Sprintf(`
		<h2>Chào mừng bạn đến với Visual Finance!</h2>
		<p>Cảm ơn bạn đã đăng ký tài khoản. Vui lòng bấm vào nút bên dưới để xác thực địa chỉ email của bạn:</p>
		<a href="%s" style="display:inline-block;padding:10px 20px;color:#fff;background-color:#007BFF;text-decoration:none;border-radius:5px;">Xác thực tài khoản</a>
		<p>Hoặc copy đường link này dán vào trình duyệt: <br> %s</p>
		<p>Nếu bạn không yêu cầu đăng ký, vui lòng bỏ qua email này.</p>
	`, verifyLink, verifyLink)

	m.SetBody("text/html", body)

	// Gmail SMTP settings
	d := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	return nil
}

// SendOTPEmail sends an email with the 6-digit OTP code for password reset
func SendOTPEmail(toEmail string, otpCode string) error {
	fromEmail := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASS")

	if fromEmail == "" || password == "" {
		return fmt.Errorf("email credentials are not set in environment variables")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Mã xác nhận quên mật khẩu Visual Finance")

	body := fmt.Sprintf(`
		<h2>Quên mật khẩu Visual Finance</h2>
		<p>Bạn vừa yêu cầu đặt lại mật khẩu. Mã OTP của bạn là:</p>
		<h1 style="color:#007BFF; letter-spacing: 5px;">%s</h1>
		<p>Mã này có hiệu lực trong vòng 5 phút. Vui lòng không chia sẻ mã này cho bất kỳ ai.</p>
		<p>Nếu bạn không yêu cầu, vui lòng bỏ qua email này.</p>
	`, otpCode)

	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	return nil
}
