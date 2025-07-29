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
	Code             *string `json:"code"`
}

type ActiveProceeding struct {
	Signature            *string         `json:"signature"`
	ProceedingsID        string          `json:"proceedingsId"`
	Type                 ProceedingsType `json:"proceedingsType"`
	DeadlineDecisionDate *string         `json:"deadlineDecisionDate"`
	Status               int             `json:"status"`
	SubmitDate           *string         `json:"submitDate"`
	ForeignerFullName    string          `json:"foreignerFullName"`
}

type ProceedingsType struct {
	ID        string `json:"id"`
	Polish    string `json:"polish"`
	English   string `json:"english"`
	Russian   string `json:"russian"`
	Ukrainian string `json:"ukrainian"`
	GroupBy   string `json:"groupBy"`
	OrderBy   string `json:"orderBy"`
	Active    bool   `json:"active"`
}
