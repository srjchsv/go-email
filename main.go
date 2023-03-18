package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"strings"
)

type Email struct {
	From     string
	Password string
	To       string
	Subject  string
	Body     string
	Provider string
}

func (e *Email) Send(auth smtp.Auth) error {
	body, err := base64.StdEncoding.DecodeString(e.Body)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
		"%s\r\n", e.From, e.To, e.Subject, body)

	err = smtp.SendMail(fmt.Sprintf("%s:587", e.Provider), auth, e.From, []string{e.To}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl := template.Must(template.ParseFiles("index.html"))
			tmpl.Execute(w, nil)
		} else if r.Method == "POST" {
			from := r.FormValue("from")
			password := r.FormValue("password")
			to := strings.Split(r.FormValue("to"), " ")
			subject := r.FormValue("subject")
			body := r.FormValue("body")
			provider := r.FormValue("provider")

			auth := smtp.PlainAuth("", from, password, provider)

			for _, recipient := range to {
				email := &Email{
					From:     from,
					Password: password,
					To:       strings.TrimSpace(recipient),
					Subject:  subject,
					Body:     base64.StdEncoding.EncodeToString([]byte(body)),
					Provider: provider,
				}

				err := email.Send(auth)
				if err != nil {
					log.Printf("Error sending email to %s: %v", recipient, err)
				} else {
					log.Printf("Email sent to %s", recipient)
				}
			}
			fmt.Fprintln(w, "Email(s) sent!")
		}
	})
	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
