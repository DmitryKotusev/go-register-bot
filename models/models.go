package models

import "time"

type ApplicationData struct {
	LoginData             LoginData
	ProceedingsCheckIndex int
}

type LoginData struct {
	Email    string
	Password string
}

type LoginPayload struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	ExpiryMinutes int    `json:"expiryMinutes"`
}

type ReservePayload struct {
	ProceedingID string `json:"proceedingId"`
	SlotID       int64  `json:"slotId"`
	Name         string `json:"name"`
	LastName     string `json:"lastName"`
	DateOfBirth  string `json:"dateOfBirth"` // ISO8601
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

type ReservationQueue struct {
	Localization string `json:"localization"`
	Prefix       string `json:"prefix"`
	ID           string `json:"id"`
	Polish       string `json:"polish"`
	English      string `json:"english"`
	Russian      string `json:"russian"`
	Ukrainian    string `json:"ukrainian"`
}

type Slot struct {
	ID    int    `json:"id"`
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// Proceeding detailed data models
type DetailedProceedingData struct {
	ID                  string       `json:"id"`
	Signature           *string      `json:"signature"`
	Circumstance        Translation  `json:"circumstance"`
	Type                Translation  `json:"type"`
	Status              string       `json:"status"`
	Description         string       `json:"description"`
	DecisionDate        *time.Time   `json:"decisionDate"`
	EditDate            time.Time    `json:"editDate"`
	CreationDate        time.Time    `json:"creationDate"`
	Person              Person       `json:"person"`
	ProxyPersons        []any        `json:"proxyPersons"`
	TimelineEvents      []Event      `json:"timelineEvents"`
	RelatedDocs         []RelatedDoc `json:"relatedDocs"`
	CaseFilesAccessInfo *string      `json:"caseFilesAccessInfo"`
	RelatedProceedings  []any        `json:"relatedProceedings"`
	CanMakeAppointment  bool         `json:"canMakeAppointment"`
	CircumstanceText    string       `json:"circumstanceText"`
}

type Translation struct {
	ID        string  `json:"id"`
	Polish    *string `json:"polish"`
	English   *string `json:"english"`
	Russian   *string `json:"russian"`
	Ukrainian *string `json:"ukrainian"`
	GroupBy   *string `json:"groupBy,omitempty"`
	OrderBy   *string `json:"orderBy,omitempty"`
	Active    *bool   `json:"active,omitempty"`
}

type Person struct {
	ResidenceAddress       Address     `json:"residenceAddress"`
	PostalAddress          Address     `json:"postalAddress"`
	PhoneNumber            string      `json:"phoneNumber"`
	Email                  string      `json:"email"`
	IdentityDocumentNumber string      `json:"identityDocumentDocumentNumber"`
	IdentityDocumentType   Translation `json:"identityDocumentType"`
	DateOfBirth            string      `json:"dateOfBirth"`
	ID                     string      `json:"id"`
	Surname                string      `json:"surname"`
	FirstName              string      `json:"firstName"`
	SecondName             string      `json:"secondName"`
}

type Address struct {
	Province                  string `json:"province"`
	County                    string `json:"county"`
	Community                 string `json:"community"`
	Locality                  string `json:"locality"`
	ZipCode                   string `json:"zipCode"`
	Street                    string `json:"street"`
	HouseNumber               string `json:"houseNumber"`
	ApartmentNumberWithPrefix string `json:"apartmentNumberWithPrefix"`
	ApartmentNumber           string `json:"apartmentNumber"`
}

type Event struct {
	// Can be any of the following:
	// "AppointmentMade", "Created"
	EventType string      `json:"eventType"`
	Date      time.Time   `json:"date"`
	Name      Translation `json:"name"`
	Author    string      `json:"author"`
}

type RelatedDoc struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"createdDate"`
}

//////////////////////////////////
