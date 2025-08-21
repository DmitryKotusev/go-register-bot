package requests

import (
	"bot-main/models"
	modelerrors "bot-main/models/errors"
	"bot-main/requests/activeproceedings"
	"bot-main/requests/cookiesinit"
	"bot-main/requests/dates"
	"bot-main/requests/dateslots"
	"bot-main/requests/login"
	"bot-main/requests/proceeding"
	"bot-main/requests/reservationqueues"
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

	//////////////////////////////////////////////////////
	time.Sleep(time.Duration(rand.Float32()) * time.Second)

	fmt.Println()
	fmt.Println("RequestPipeline, trying to login...")
	sessionToken, err := login.Login(client, applicationData.LoginData)
	if err != nil {
		fmt.Printf("RequestPipeline error during login: %v", err)
		return err
	}
	fmt.Printf("Login request completed successfully, token: %s.\n", sessionToken)

	//////////////////////////////////////////////////////
	time.Sleep(time.Duration(rand.Float32()) * time.Second)

	fmt.Println()
	fmt.Println("RequestPipeline, trying to get active proceedings...")
	activeProceedings, err := activeproceedings.GetActiveProceedings(client, sessionToken)
	if err != nil {
		fmt.Printf("RequestPipeline error during getting active proceedings: %v", err)
		return err
	}
	fmt.Println("Get active proceedings request completed successfully, proceedings:")
	printData(activeProceedings)
	if len(activeProceedings) <= applicationData.ProceedingsCheckIndex {
		fmt.Println("RequestPipeline, proceedings length and index incompatibility, returning error.")
		return modelerrors.ProceedingsCountError{
			Message: fmt.Sprintf("❌ RequestPipeline failed because proceedings count and index incompatibility: %d and %d.",
				len(activeProceedings),
				applicationData.ProceedingsCheckIndex),
		}
	}

	//////////////////////////////////////////////////////
	time.Sleep(time.Duration(rand.Float32()) * time.Second)
	relevantProceeding := activeProceedings[applicationData.ProceedingsCheckIndex]

	fmt.Println()
	fmt.Printf("RequestPipeline, trying to get detailed info about proceeding %s...\n", relevantProceeding.ProceedingsID)
	proceedingData, err := proceeding.GetProceedingData(client, sessionToken, relevantProceeding)
	if err != nil {
		fmt.Printf("RequestPipeline error during getting detailed proceeding data: %v", err)
		return err
	}
	fmt.Printf("Get detailed proceeding data for %s completed successfully, data:\n", relevantProceeding.ProceedingsID)
	printData(proceedingData)

	//////////////////////////////////////////////////////
	time.Sleep(time.Duration(rand.Float32()) * time.Second)

	fmt.Println()
	fmt.Printf("RequestPipeline, trying to get queues for reservation for proceeding %s...\n", proceedingData.ID)
	reservationQueues, err := reservationqueues.GetReservationQueues(client, sessionToken, proceedingData)
	if err != nil {
		fmt.Printf("RequestPipeline error during getting reservation queues: %v", err)
		return err
	}
	fmt.Printf("Get reservation queues for %s completed successfully, queues:\n", relevantProceeding.ProceedingsID)
	printData(reservationQueues)

	//////////////////////////////////////////////////////
	time.Sleep(time.Duration(rand.Float32()) * time.Second)

	relevantQueue := reservationQueues[0]
	fmt.Println()
	fmt.Printf("RequestPipeline, trying to get dates for query %s...\n", relevantQueue.Localization)
	queueDates, err := dates.GetReservationQueueDates(client, sessionToken, proceedingData, relevantQueue)
	if err != nil {
		fmt.Printf("RequestPipeline error during getting queue dates: %v", err)
		return err
	}
	fmt.Printf("Get queue dates for %s completed successfully, dates:\n", relevantQueue.ID)
	printData(queueDates)

	//////////////////////////////////////////////////////
	time.Sleep(time.Duration(rand.Float32()) * time.Second)

	queueDate := queueDates[0]
	fmt.Println()
	fmt.Printf("RequestPipeline, trying to get date slots for date %s at %s...\n", queueDate, relevantQueue.Localization)
	queueDateSlots, err := dateslots.GetReservationQueueDateSlots(client, sessionToken, proceedingData, relevantQueue, queueDate)
	if err != nil {
		fmt.Printf("RequestPipeline error during getting queue date slots: %v", err)
		return err
	}
	fmt.Printf("Get queue date slots for %s completed successfully, date slots:\n", relevantQueue.Localization)
	printData(queueDateSlots)

	return nil
}

func printData(input any) {
	data, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data))
}
