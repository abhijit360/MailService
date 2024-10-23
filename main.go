package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	gomail "gopkg.in/mail.v2"
)

type ContactForm struct {
	From    string `json:"from,omitempty"`
	Subject string `json:"subject,omitempty"`
	Content string `json:"content,omitempty"`
}

type Message struct {
	Status string `json:"status"`
	Body string `json:"body"`
}

func rateLimiter(eventRate float64, burstRate int ,function func (w http.ResponseWriter,r *http.Request)) http.Handler{
	limiter := rate.NewLimiter(rate.Limit(eventRate),burstRate) // rate of 
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !limiter.Allow(){
			errorMessage := Message{
				Status: "Request Failed",
				Body: "You are being rate limited.",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&errorMessage)
			return
		}else{
			function(w,r)
		}
	})
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Please set up an .env file with a key of 'my mail' = name@domain.com and 'my_password' = your email password")
	}
	MY_MAIL := os.Getenv("my_mail")
	PASSWORD := os.Getenv("my_password")

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
			d := gomail.NewDialer("smtp.gmail.com", 587, MY_MAIL,PASSWORD)
			if err := d.DialAndSend(m); err != nil{
				log.Fatal("Failed Sending email to personal email",err)
			}

			// send message to sender that I have received email
			m.SetHeader("From", MY_MAIL)
			m.SetHeader("To", contactData.From)
			m.SetHeader("Subject", contactData.Subject)
			m.SetBody("text/html", `Hello There,<br>
				<br>
				Thank you for reaching out to me through the contact form on my website!<br>
				This is an automated response to let you know that I have received the email and I will get back to you shortly.<br>
				<br>
				Yours Sincerely,<br>
				Abhijit`)

			d = gomail.NewDialer("smtp.gmail.com", 587, MY_MAIL, PASSWORD)
			if err := d.DialAndSend(m); err != nil{
				log.Fatal("Failed sending update to client",err)
			}
		}
	}
	http.Handle("/receive", rateLimiter(1,2,receiveMail))

	err = http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
