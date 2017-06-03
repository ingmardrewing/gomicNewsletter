package main

import (
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/ingmardrewing/gomicNewsletter/config"
	"github.com/ingmardrewing/gomicNewsletter/db"
	"github.com/ingmardrewing/gomicNewsletter/service"
)

func main() {
	db.Initialize()
	restful.Add(service.NewNewsletterService())

	port := "16443"

	log.Println("Reading crt and key data from files:")
	crt, key := config.GetTlsPaths()

	log.Println("Path to crt file: " + crt)
	log.Println("Path to key file: " + key)
	log.Println("Starting to serve via TLS on Port: " + port)

	err := http.ListenAndServeTLS(":"+port, crt, key, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
