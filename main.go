package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"log"
	gomail "gopkg.in/mail.v2"
	"github.com/joho/godotenv"
	
)

type ContactForm struct {
	From      string
	Subject string
	Content string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Please set up an .env file with a key of 'my mail' = name@domain.com")
	}
	MY_MAIL := os.Getenv("my_mail")

	receiveMail := func(w http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			jsonEncoder := json.NewEncoder(w)
			jsonEncoder.Encode("This is a POST endpoint to handle contact form submissions")
		} else if request.Method == "POST" {
			contactData := ContactForm{}
			body, err := io.ReadAll(request.Body)
			if err != nil {
				fmt.Print("Error in parsing json body", err)
				os.Exit(1)
			}

			err = json.Unmarshal(body, &contactData) // UnmarshaIlling here
			if err != nil {
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}

			fmt.Print(contactData)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			jsonEncoder := json.NewEncoder(w)
			jsonEncoder.Encode("Message: Email Received.")
			
			m := gomail.NewMessage()
			//
			m.SetHeader("From", MY_MAIL)
			m.SetHeader("To",MY_MAIL)
			m.SetHeader("Subject", "Personal Website Contact Me Form")
			m.SetBody("text/html", fmt.Sprintf("Message from : %v\n Subject: %v \nBody: %v",contactData.From, contactData.Subject,contactData.Content))
			d := gomail.NewDialer("smtp.gmail.com", 587, userName, password)
			if err := d.DialAndSend(m); err != nil{
				log.Fatal("Failed Sending email to personal email",err)
			}

			// send message to sender that I have received email
			m.SetHeader("From", MY_MAIL)
			m.SetHeader("To", contactData.From)
			m.SetHeader("Subject", contactData.Subject)
			m.SetBody("text/html", `Hello!<br>
				Thank you for reaching out to me through the contact form on my website!<br>
				This is an automated response to let you know that I have received the email and I will get back to you shortly.<br>
				Yours Sincerely,<br>
				Abhijit`)

			d = gomail.NewDialer("smtp.gmail.com", 587, userName, password)
			if err := d.DialAndSend(m); err != nil{
				log.Fatal("Failed sending update to client",err)
			}
		}
	}
	http.HandleFunc("/receive", receiveMail)

	err = http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
