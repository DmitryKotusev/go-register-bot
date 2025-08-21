package reservationqueues

import (
	"bot-main/globalvars"
	"bot-main/models"
	modelerrors "bot-main/models/errors"
	"bot-main/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetReservationQueues(client *http.Client, sessionToken string, proceeding *models.DetailedProceedingData) ([]models.ReservationQueue, error) {
	if client == nil {
		return nil, fmt.Errorf("GetReservationQueues, HTTP client is nil")
	}
	if proceeding == nil {
		return nil, fmt.Errorf("GetReservationQueues, proceeding data is nil")
	}
	getProceedingReservationQueuesRequestUrl := fmt.Sprintf(globalvars.GetProceedingReservationQueuesRequestUrl, proceeding.ID)
	homePageCasesUrl := fmt.Sprintf(globalvars.HomePageCasesUrl, proceeding.ID)
	req, err := http.NewRequest("GET", getProceedingReservationQueuesRequestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueues request error creating request: %v", err)
	}
	// Set headers
	req.Header.Set("Referer", homePageCasesUrl)
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	utils.AttachDefaultRequestHeaders(req)

	fmt.Println("Sending GetReservationQueues request...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueues request error executing: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, modelerrors.UnauthorizedError{
			Message: fmt.Sprintf("❌ GetReservationQueues failed because of unauthorized status code: %s", resp.Status),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetReservationQueues request failed with status: %s", resp.Status)
	}

	fmt.Printf("GetReservationQueues response: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueues request error reading response body: %v", err)
	}

	var reservationQueues []models.ReservationQueue
	err = json.Unmarshal(body, &reservationQueues)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueues body JSON parcing error: %v", err)
	}

	return reservationQueues, nil
}
