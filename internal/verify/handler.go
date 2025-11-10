package verify

import (
	"email-verify/configs"
	"email-verify/internal/repository"
	"email-verify/pkg/email"
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
	emailSender := email.NewEmailSender(deps.Config)

	router.HandleFunc("POST /send", handler.Send(repo, emailSender))
	router.HandleFunc("GET /verify/{hash}", handler.Verify(repo))
}

func (handler *VerifyHandler) Send(repo repository.VerificationRepository, emailSender *email.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Send")
		body, err := req.HandleBody[EmailRequest](&w, r)
		if err != nil {
			return
		}
		log.Println("Received email:", body.Email)

		hashValue, err := hash.Generate()
		if err != nil {
			log.Println("error generating hash:", err)
			res.JsonResp(w, map[string]string{"error": "Failed to generate hash"}, 500)
			return
		}
		log.Println("Generated hash:", hashValue)

		// Сохраняем в репозиторий
		err = repo.Create(body.Email, hashValue)
		if err != nil {
			log.Println("error create record in repo:", err)
			res.JsonResp(w, map[string]string{"error": "Failed to save verification"}, 500)
			return
		}

		// Отправляем email
		err = emailSender.SendVerificationEmail(body.Email, hashValue)
		if err != nil {
			log.Println("error sending email:", err)

			// Если не удалось отправить email, удаляем запись из репозитория
			deleteErr := repo.Delete(hashValue)
			if deleteErr != nil {
				log.Println("error cleaning up failed verification:", deleteErr)
			}

			res.JsonResp(w, map[string]string{"error": "Failed to send verification email"}, 500)
			return
		}

		res.JsonResp(w, map[string]string{
			"message": "Verification email sent successfully",
			"email":   body.Email,
		}, 200)
	}
}

func (handler *VerifyHandler) Verify(repo repository.VerificationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Verify")

		// Получаем хеш из URL параметра
		hash := r.PathValue("hash")
		if hash == "" {
			res.JsonResp(w, map[string]string{"error": "Hash is required"}, 400)
			return
		}

		log.Println("Verifying hash:", hash)

		// Ищем верификацию по хешу
		verification, err := repo.FindByHash(hash)
		if err != nil {
			log.Println("error finding verification:", err)
			res.JsonResp(w, map[string]string{"error": "Internal server error"}, 500)
			return
		}

		if verification == nil {
			// Хеш не найден
			res.JsonResp(w, map[string]bool{"verified": false}, 200)
			return
		}

		// Хеш найден - удаляем запись и возвращаем true
		err = repo.Delete(hash)
		if err != nil {
			log.Println("error deleting verification:", err)
			// Но все равно возвращаем true, так как верификация успешна
		}

		log.Printf("Email %s verified successfully", verification.Email)
		res.JsonResp(w, VerifyResponse{
			Verified: true,
		}, 200)
	}
}
