package cookiesinit

import (
	"bot-main/globalvars"
	"fmt"
	"log"
	"net/http"
)

func CookiesInit(client *http.Client) error {
	if client == nil {
		return fmt.Errorf("HTTP client is nil")
	}

	fmt.Printf("CookiesInit, sending GET-request to %s to get cookie...\n", globalvars.LoginPageUrl)

	// Creating request
	preReq, err := http.NewRequest("GET", globalvars.LoginPageUrl, nil)
	if err != nil {
		log.Fatalf("Error creating GET-request: %v", err)
	}
	// Setting headers similar to real browser
	attachHeaders(preReq)
	preResp, err := client.Do(preReq)
	if err != nil {
		return fmt.Errorf("CookiesInit request error executing: %v", err)
	}
	defer preResp.Body.Close()
	fmt.Println("✅ CookiesInit, cookies initialized successfully!")
	return nil
}

func attachHeaders(req *http.Request) {
	req.Header.Set("User-Agent", globalvars.DefaultUserAgent)
	req.Header.Set("Accept", globalvars.HtmlAcceptHeader)
	req.Header.Set("Connection", globalvars.KeepAliveHeader)
	req.Header.Set("Accept-Encoding", globalvars.AcceptEncodingHeader)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", globalvars.Origin)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
}
