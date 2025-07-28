package requests

import (
	"bot-main/models"
	"bot-main/requests/cookiesinit"
	"bot-main/requests/login"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func RequestPipeline(loginData models.LoginData) error {
	// Creating custom transport, disabling HTTP/2.
	// We are cloning default transport and changing only one setting.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		NextProtos: []string{"http/1.1"},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("RequestPipeline error creating cookie jar: %v", err)
	}
	client := &http.Client{
		Jar:       jar,
		Transport: transport,
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
	sessionToken, err := login.Login(client, loginData)
	if err != nil {
		fmt.Printf("RequestPipeline error during login: %v", err)
		return err
	}
	fmt.Printf("Login request completed successfully, token: %s.\n", sessionToken)

	return nil
}
