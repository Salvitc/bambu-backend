package mailing

import (
	"fmt"
	"net/smtp"
)

type Message struct {
	From    string `json:"from" bson:"from,omitempty"`
	To      string `json:"to" bson:"to,omitempty"`
	Subject string `json:"subject" bson:"subject,omitempty"`
	Body    string `json:"body" bson:"body,omitempty"`
}

func SendMail(message Message) error {
	from := "bambushop46@gmail.com"
	pass := "xcgi upyd qvym xmyc"
	to := "bambushop46@gmail.com"

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		message.From, message.To, message.Subject, message.Body)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	return err
}
