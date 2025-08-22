package reserve

import (
	"bot-main/globalvars"
	"bot-main/models"
	modelerrors "bot-main/models/errors"
	"bot-main/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func ReserveDateSlot(
	client *http.Client,
	sessionToken string,
	proceeding *models.DetailedProceedingData,
	reservationQueue models.ReservationQueue,
	dateSlot models.Slot) error {
	if client == nil {
		return fmt.Errorf("ReserveDateSlot, HTTP client is nil")
	}
	if proceeding == nil {
		return fmt.Errorf("ReserveDateSlot, proceeding data is nil")
	}
	payload := models.ReservePayload{
		ProceedingID: proceeding.ID,
		SlotID:       int64(dateSlot.ID),
		Name:         proceeding.Person.FirstName,
		LastName:     proceeding.Person.Surname,
		DateOfBirth:  proceeding.Person.DateOfBirth,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ReserveDateSlot request error encoding JSON: %v", err)
	}
	reserveAppointmentRequestUrl := fmt.Sprintf(globalvars.ReserveAppointmentRequestUrl, reservationQueue.ID)
	homePageCasesUrl := fmt.Sprintf(globalvars.HomePageCasesUrl, proceeding.ID)
	req, err := http.NewRequest("POST", reserveAppointmentRequestUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("ReserveDateSlot request error creating request: %v", err)
	}
	// Set headers
	req.Header.Set("Referer", homePageCasesUrl)
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	utils.AttachDefaultRequestHeaders(req)

	fmt.Println("Sending ReserveDateSlot request...")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ReserveDateSlot request error executing: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return modelerrors.UnauthorizedError{
			Message: fmt.Sprintf("❌ ReserveDateSlot failed because of unauthorized status code: %s", resp.Status),
		}
	} else if resp.StatusCode == http.StatusForbidden {
		return modelerrors.ForbiddenError{
			Message: fmt.Sprintf("❌ ReserveDateSlot failed because of forbidden status code: %s, probably needs cookies update", resp.Status),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ReserveDateSlot request for %s failed with status: %s", dateSlot.Date, resp.Status)
	}

	fmt.Printf("ReserveDateSlot %s response: %s\n", dateSlot.Date, resp.Status)

	return nil
}
