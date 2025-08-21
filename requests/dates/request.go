package dates

import (
	"bot-main/globalvars"
	"bot-main/models"
	modelerrors "bot-main/models/errors"
	"bot-main/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetReservationQueueDates(
	client *http.Client,
	sessionToken string,
	proceeding *models.DetailedProceedingData,
	reservationQueue models.ReservationQueue) ([]string, error) {
	if client == nil {
		return nil, fmt.Errorf("GetReservationQueueDates, HTTP client is nil")
	}
	if proceeding == nil {
		return nil, fmt.Errorf("GetReservationQueueDates, proceeding data is nil")
	}
	getReservationQueueDatesRequestUrl := fmt.Sprintf(globalvars.GetReservationQueueDatesRequestUrl, reservationQueue.ID)
	homePageCasesUrl := fmt.Sprintf(globalvars.HomePageCasesUrl, proceeding.ID)
	req, err := http.NewRequest("POST", getReservationQueueDatesRequestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueueDates request error creating request: %v", err)
	}
	// Set headers
	req.Header.Set("Referer", homePageCasesUrl)
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	utils.AttachDefaultRequestHeaders(req)

	fmt.Println("Sending GetReservationQueueDates request...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueueDates request error executing: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, modelerrors.UnauthorizedError{
			Message: fmt.Sprintf("❌ GetReservationQueueDates failed because of unauthorized status code: %s", resp.Status),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetReservationQueueDates request failed with status: %s", resp.Status)
	}

	fmt.Printf("GetReservationQueueDates response: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueueDates request error reading response body: %v", err)
	}

	var reservationQueueDates []string
	err = json.Unmarshal(body, &reservationQueueDates)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueueDates body JSON parcing error: %v", err)
	}

	for i, date := range reservationQueueDates {
		newDate := strings.Split(date, "T")[0]
		reservationQueueDates[i] = newDate
	}

	return reservationQueueDates, nil
}
