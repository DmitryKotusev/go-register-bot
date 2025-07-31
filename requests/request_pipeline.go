package requests

import (
	"bot-main/models"
	modelerrors "bot-main/models/errors"
	"bot-main/requests/activeproceedings"
	"bot-main/requests/cookiesinit"
	"bot-main/requests/login"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func RequestPipeline(applicationData models.ApplicationData) error {
	// Creating custom transport, disabling HTTP/2.
	// We are cloning default transport and changing only one setting.
	// TODO: Consider using a custom transport to be able to auto uncompress gzip responses.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		NextProtos: []string{"http/1.1"},
	}
	// Disabling compression as we are doing it ourselves.
	transport.DisableCompression = true

	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("RequestPipeline error creating cookie jar: %v", err)
	}
	client := &http.Client{
		Jar:       jar,
		Transport: &DecompressingTransport{Transport: transport},
	}

	fmt.Println()
	fmt.Println("RequestPipeline started, initializing cookies...")
	err = cookiesinit.CookiesInit(client)
	if err != nil {
		fmt.Printf("RequestPipeline error during initializing cookies: %v", err)
		return err
	}

	time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

	fmt.Println()
	fmt.Println("RequestPipeline, trying to login...")
	sessionToken, err := login.Login(client, applicationData.LoginData)
	if err != nil {
		fmt.Printf("RequestPipeline error during login: %v", err)
		return err
	}
	fmt.Printf("Login request completed successfully, token: %s.\n", sessionToken)

	time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

	fmt.Println()
	fmt.Println("RequestPipeline, trying to get active proceedings...")
	activeProceedings, err := activeproceedings.GetActiveProceedings(client, sessionToken)
	if err != nil {
		fmt.Printf("RequestPipeline error during getting active proceedings: %v", err)
		return err
	}
	fmt.Println("Get active proceedings request completed successfully, proceedings:")
	printActiveProceedingsSlice(activeProceedings)
	if len(activeProceedings) <= applicationData.ProceedingsCheckIndex {
		fmt.Println("RequestPipeline, proceedings length and index incompatibility, returning error.")
		return modelerrors.ProceedingsCountError{
			Message: fmt.Sprintf("❌ RequestPipeline failed because proceedings count and index incompatibility: %d and %d.",
				len(activeProceedings),
				applicationData.ProceedingsCheckIndex),
		}
	}

	time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

	relevantProceeding := activeProceedings[len(activeProceedings)-1-applicationData.ProceedingsCheckIndex]
	fmt.Println()
	fmt.Printf("RequestPipeline, trying to get queues for reservation %s...\n", relevantProceeding.ProceedingsID)

	return nil
}

func printActiveProceedingsSlice(slice []models.ActiveProceeding) {
	data, err := json.MarshalIndent(slice, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data))
}
