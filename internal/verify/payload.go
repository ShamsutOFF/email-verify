package verify

type EmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyResponse struct {
	Verified bool `json:"verified"`
}
