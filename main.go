package main

import (
	// "fmt"
	"io"
	"log"
	"net/http"
)

func main(){
	
	receiveMail := func(w http.ResponseWriter, _ *http.Request){
		io.WriteString(w,"Hello from receiveMail");
	}
	http.HandleFunc("/receive", receiveMail)

	log.Fatal(http.ListenAndServe(":8080", nil))
}



