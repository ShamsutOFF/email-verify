package main

import (
	"email-verify/configs"
	"email-verify/internal/verify"
	"log"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	router := http.NewServeMux()
	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{Config: conf})

	server := http.Server{
		Addr:    ":7777",
		Handler: router,
	}

	log.Println("Server listening on port 7777")
	err := server.ListenAndServe()
	if err != nil {
		log.Println("Server listening error")
		return
	}
}
