package server

import (
	"log"
	"net/http"
)

func Serve() {
	// directory where index.html is located
	dir := "D:/SiteVisualization/"

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
