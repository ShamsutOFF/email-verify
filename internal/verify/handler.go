package verify

import (
	"email-verify/configs"
	"email-verify/internal/repository"
	"email-verify/pkg/hash"
	"email-verify/pkg/req"
	"email-verify/pkg/res"
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
	repo := repository.NewFileVerificationRepo("./verifications.json")
	router.HandleFunc("POST /send", handler.Send(repo))
	router.HandleFunc("GET /verify/{hash}", handler.Verify())
}

func (handler *VerifyHandler) Send(repo repository.VerificationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Send")
		body, err := req.HandleBody[EmailRequest](&w, r)
		if err != nil {
			return
		}
		log.Println(body)
		hashValue := hash.Generate()
		log.Println("Generated hash:", hashValue)

		err = repo.Create(body.Email, hashValue)
		if err != nil {
			log.Println("error create record in repo:", err)
			return
		}

		res.JsonResp(w, map[string]string{
			"email": body.Email,
			"hash":  hashValue,
		}, 200)
	}
}

func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Verify")
	}
}
