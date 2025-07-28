package models

type LoginData struct {
	Email    string
	Password string
}

type LoginPayload struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	ExpiryMinutes int    `json:"expiryMinutes"`
}

type LoginResponse struct {
	IsAuthSuccessful bool    `json:"isAuthSuccessful"`
	ErrorMessage     any     `json:"errorMessage"`
	Token            string  `json:"token"`
	Code             *string `json:"optionalField"`
}
