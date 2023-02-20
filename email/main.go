package main

import (
	"fmt"
	"net/smtp"
)

func main() {
	fmt.Println("golang email sending")
	email()
}

func email() {
	// sender data
	from := "toewailinbdl@proton.me"
	pass := "Toewailin2016"

	// receiver data
	toEmail := "dasapi9538@ibansko.com"
	to := []string{toEmail}

	// smtp
	host := "smtp.gmail.com"
	// port := "587"
	// port := "465"
	port := "1025"
	address := host + ":" + port

	//message
	subject := "Subject: Our Golang Email Testing\n"
	body := "our first email"
	message := []byte(subject + body)
	auth := smtp.PlainAuth("", from, pass, host)
	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	fmt.Println("go check your email")
}
