package email

import (
	"crypto/tls"
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
	d.TLSConfig = &tls.Config{InsecureSkipVerify: false}

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	return nil
}
