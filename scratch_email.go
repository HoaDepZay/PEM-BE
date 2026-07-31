package main

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

func main() {
	fromEmail := "servernodejs26@gmail.com"
	password := "nunl rcgu pmlw fsfu"
	toEmail := "dangquanghoa206@gmail.com"

	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Test Email")
	m.SetBody("text/plain", "This is a test email")

	d := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, password)

	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("could not send email: %v", err)
	}
	fmt.Println("Email sent successfully!")
}
