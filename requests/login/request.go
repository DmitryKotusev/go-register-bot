package login

import (
	"bot-main/globalvars"
	"bot-main/models"
	modelerrors "bot-main/models/errors"
	"bot-main/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Login(client *http.Client, loginData models.LoginData) (string, error) {
	if client == nil {
		return "", fmt.Errorf("Login, HTTP client is nil")
	}
	loginURL := globalvars.LoginRequestUrl
	userEmail := loginData.Email
	userPassword := loginData.Password
	payload := models.LoginPayload{
		Email:         userEmail,
		Password:      userPassword,
		ExpiryMinutes: 0,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("Login request error encoding JSON: %v", err)
	}

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("Login request error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Referer", globalvars.LoginPageUrl)
	utils.AttachDefaultRequestHeaders(req)

	fmt.Println("Sending login request...")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Login request error executing: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && !(http.StatusBadRequest <= resp.StatusCode && resp.StatusCode < 500) {
		return "", fmt.Errorf("Login request failed with status: %s", resp.Status)
	}

	fmt.Printf("Login response: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Login request error reading response body: %v", err)
	}

	var loginResp models.LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return "", fmt.Errorf("Login body JSON parcing error: %v", err)
	}

	if loginResp.IsAuthSuccessful {
		fmt.Printf("✅ Successful login!\n")
	} else {
		if loginResp.Code != nil {
			return "", modelerrors.InvalidCredentailsError{
				Message: fmt.Sprintf("❌ Login failed because of wrong credentials, code %s", *loginResp.Code),
			}
		}
		return "", fmt.Errorf("❌ Login failed: %v", loginResp.ErrorMessage)
	}

	return loginResp.Token, nil
}
