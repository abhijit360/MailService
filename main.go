package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ContactForm struct {
	From    string
	To      string
	Content string
}

func main() {
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
		}
	}
	http.HandleFunc("/receive", receiveMail)

	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
