package middlewares

import (
	"fmt"

	"gopkg.in/mail.v2"
)

func SendEmail(verifyCode string, email string,message string) bool {
	m := mail.NewMessage()
	m.SetHeader("From", "ristian.rehi@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verification Code")
	m.SetBody("text/plain", message)

	// contoh smtp
	d := mail.NewDialer("smtp.gmail.com", 587, "ristian.rehi@gmail.com", "ljfp ghgp guev eped")

	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println("Send email error:", err)
		return false
	}

	return true
}
