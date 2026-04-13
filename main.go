package main

import (
	"log"
	"shortener/handlers"
	"shortener/repository"
	"shortener/router"
	"shortener/service"
)

func main() {
	repo := repository.NewShortLinkRepo()
	svc := service.NewShortLinkService(repo)
	h := handlers.NewShortLinkHandler(svc)

	r := router.Setup(h)

	log.Println("starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
