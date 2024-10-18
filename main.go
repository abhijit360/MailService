package main

import (
	"fmt"
	"errors"
	"net/http"
	"encoding/json"
)

func main(){
	
	receiveMail := func(w http.ResponseWriter, request *http.Request){
		if request.Method == "GET"{
			w.WriteHeader(http.StatusBadRequest);
			w.Header().Set("Content-Type", "application/json")
			jsonEncoder := json.NewEncoder(w)
			jsonEncoder.Encode("This is a POST endpoint to handle contact form submissions")
		}else if request.Method == "POST"{
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			jsonEncoder := json.NewEncoder(w)
			jsonEncoder.Encode("Returning json input?")
		}	
	}
	http.HandleFunc("/receive", receiveMail)

	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed){
		fmt.Printf("server closed\n")
	}else if err != nil {
		fmt.Printf("error: %v\n",err)
	}
}



