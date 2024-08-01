package server

import (
	"log"
	"net/http"
)

func Serve() {
	//change directory location
	f := http.Dir("index.html")
	fs := http.FileServer(f)
	http.Handle("/", fs)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
