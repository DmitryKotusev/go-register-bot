package activeproceedings

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

func GetActiveProceedings(client *http.Client, sessionToken string) ([]models.ActiveProceeding, error) {
	if client == nil {
		return nil, fmt.Errorf("GetActiveProceedings, HTTP client is nil")
	}
	getActiveProceedingsRequestUrl := globalvars.GetActiveProceedingsRequestUrl
	req, err := http.NewRequest("GET", getActiveProceedingsRequestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("GetActiveProceedings request error creating request: %v", err)
	}
	// Set headers
	req.Header.Set("Referer", globalvars.HomePageUrl)
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	utils.AttachDefaultRequestHeaders(req)

	fmt.Println("Sending GetActiveProceedings request...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetActiveProceedings request error executing: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, modelerrors.UnauthorizedError{
			Message: fmt.Sprintf("❌ GetActiveProceedings failed because of unauthorized status code: %s", resp.Status),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetActiveProceedings request failed with status: %s", resp.Status)
	}

	fmt.Printf("GetActiveProceedings response: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetActiveProceedings request error reading response body: %v", err)
	}

	var activeProceedings []models.ActiveProceeding
	err = json.Unmarshal(body, &activeProceedings)
	if err != nil {
		return nil, fmt.Errorf("GetActiveProceedings body JSON parcing error: %v", err)
	}

	return activeProceedings, nil
}
