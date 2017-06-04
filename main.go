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

	wsContainer := restful.NewContainer()
	wsContainer.Add(service.NewNewsletterService())

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{""},
		AllowedHeaders: []string{"Content-Type", "X-Request-With, Content-Type, Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      wsContainer}

	wsContainer.Filter(cors.Filter)
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	port := "16443"

	log.Println("Reading crt and key data from files:")
	crt, key := config.GetTlsPaths()

	log.Println("Path to crt file: " + crt)
	log.Println("Path to key file: " + key)
	log.Println("Starting to serve via TLS on Port: " + port)

	server := &http.Server{Addr: ":" + port, Handler: wsContainer}
	err := server.ListenAndServeTLS(crt, key)
	if err != nil {
		log.Fatal(err.Error())
	}
}
