package main

import (
	"log"
	"time"

	"github.com/akhyar02/bookings/internal/models"
	"gopkg.in/gomail.v2"
)

func listenForMail() {
	go func() {
		server := gomail.NewDialer("localhost", 587, "", "")
		var s gomail.SendCloser
		open := false
		for {
			select {
			case m, ok := <-app.MailChan:
				if !ok {
					return
				}
				sendMail(m, server, &s, &open)
			case <-time.After(30 * time.Second):
				if !open {
					continue
				}
				s.Close()
				open = false
			}
		}
	}()
}

func sendMail(m models.MailData, server *gomail.Dialer, s *gomail.SendCloser, open *bool) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", m.To...)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", string(m.Content))

	if !*open {
		var err error
		*s, err = server.Dial()
		if err != nil {
			errorLog.Println(err)
			return
		} else {
			*open = true
		}
	}

	if err := gomail.Send(*s, msg); err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent")
	}

}
