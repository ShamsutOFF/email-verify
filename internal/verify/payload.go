package verify

type EmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}
