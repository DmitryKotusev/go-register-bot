package proceeding

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

func GetProceedingData(client *http.Client,
	sessionToken string,
	proceeding models.ActiveProceeding,
) (*models.DetailedProceedingData, error) {
	if client == nil {
		return nil, fmt.Errorf("GetProceedingData, HTTP client is nil")
	}
	getProceedingRequestUrl := fmt.Sprintf(globalvars.GetProceedingRequestUrl, proceeding.ProceedingsID)
	req, err := http.NewRequest("GET", getProceedingRequestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("GetProceedingData request error creating request: %v", err)
	}
	// Set headers
	req.Header.Set("Referer", globalvars.HomePageUrl)
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	utils.AttachDefaultRequestHeaders(req)

	fmt.Println("Sending GetProceedingData request...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetProceedingData request error executing: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, modelerrors.UnauthorizedError{
			Message: fmt.Sprintf("❌ GetProceedingData failed because of unauthorized status code: %s", resp.Status),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetProceedingData request failed with status: %s", resp.Status)
	}

	fmt.Printf("GetProceedingData response: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetProceedingData request error reading response body: %v", err)
	}

	var proceedingData models.DetailedProceedingData
	err = json.Unmarshal(body, &proceedingData)
	if err != nil {
		return nil, fmt.Errorf("GetProceedingData body JSON parcing error: %v", err)
	}

	return &proceedingData, nil
}
