package verify

import (
	"email-verify/configs"
	"log"
	"net/http"
)

type VerifyHandlerDeps struct {
	Config *configs.Config
}

type VerifyHandler struct {
	Config *configs.Config
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
	handler := &VerifyHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST send", handler.Send())
	router.HandleFunc("GET verify/{hash}", handler.Verify())
}

func (handler VerifyHandler) Send() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Send")
	}
}

func (handler VerifyHandler) Verify() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Verify")
	}
}
