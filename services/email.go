package services

import (
	"gopkg.in/gomail.v2"
	"os"
	"log"
	"strconv"
	"fmt"
)

func SendMail(mailTo string, mailSubject string, mailBody string) bool {
	host := os.Getenv("SMTP_HOST")
	portVar := os.Getenv("SMTP_PORT")
	port, _ := strconv.Atoi(portVar)
	email := os.Getenv("AUTH_EMAIL")
	password := os.Getenv("AUTH_PASSWORD")

	msg := gomail.NewMessage()
	msg.SetHeader("From", "<no-reply@itsfood.my.id>")
	msg.SetHeader("To", mailTo)
	msg.SetHeader("Subject", mailSubject)
	msg.SetBody("text/html", mailBody)

	dialer := gomail.NewDialer(host, port, email, password)

	errSendingEmail := dialer.DialAndSend(msg)
	if errSendingEmail != nil {
		log.Fatal(errSendingEmail.Error())
		fmt.Println(port)
		fmt.Println(host)
		return false
	}
	
	return true
}