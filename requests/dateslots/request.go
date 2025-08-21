package dateslots

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

func GetReservationQueueDateSlots(
	client *http.Client,
	sessionToken string,
	proceeding *models.DetailedProceedingData,
	reservationQueue models.ReservationQueue,
	simpleDate string) ([]models.Slot, error) {
	if client == nil {
		return nil, fmt.Errorf("GetReservationQueueDateSlots, HTTP client is nil")
	}
	if proceeding == nil {
		return nil, fmt.Errorf("GetReservationQueueDateSlots, proceeding data is nil")
	}
	getReservationQueueDateSlotsRequestUrl := fmt.Sprintf(globalvars.GetReservationQueueDateSlotsRequestUrl, reservationQueue.ID, simpleDate)
	homePageCasesUrl := fmt.Sprintf(globalvars.HomePageCasesUrl, proceeding.ID)
	req, err := http.NewRequest("POST", getReservationQueueDateSlotsRequestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueueDateSlots request error creating request: %v", err)
	}
	// Set headers
	req.Header.Set("Referer", homePageCasesUrl)
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	utils.AttachDefaultRequestHeaders(req)

	fmt.Println("Sending GetReservationQueueDateSlots request...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueueDateSlots request error executing: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, modelerrors.UnauthorizedError{
			Message: fmt.Sprintf("❌ GetReservationQueueDateSlots failed because of unauthorized status code: %s", resp.Status),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetReservationQueueDateSlots request failed with status: %s", resp.Status)
	}

	fmt.Printf("GetReservationQueueDateSlots response: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueueDateSlots request error reading response body: %v", err)
	}

	var reservationQueueDateSlots []models.Slot
	err = json.Unmarshal(body, &reservationQueueDateSlots)
	if err != nil {
		return nil, fmt.Errorf("GetReservationQueueDateSlots body JSON parcing error: %v", err)
	}

	return reservationQueueDateSlots, nil
}
