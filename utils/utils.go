package utils

import (
	"bot-main/globalvars"
	"bot-main/models"
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func RegisterCommandLineArgs() {
	flag.StringVar(&globalvars.Email, "email", "", "Login email for enter")
	flag.StringVar(&globalvars.Password, "password", "", "Password for enter")
	flag.IntVar(&globalvars.ProceedingsCheckIndex, "proceedings-check-index", 0, "Proceedings check index for enter from end (by default 0)")
	flag.Parse()
}

func ReadRequiredApplicationData() models.ApplicationData {
	return models.ApplicationData{
		LoginData:             ReadRequiredLoginData(),
		ProceedingsCheckIndex: globalvars.ProceedingsCheckIndex,
	}
}

func ReadRequiredLoginData() models.LoginData {
	// Checking if data was entered. If not, request it.
	if globalvars.Email == "" {
		globalvars.Email = ReadStringFromConsole("Enter email: ")
	}

	if globalvars.Password == "" {
		globalvars.Password = ReadStringFromConsole("Enter password: ")
	}

	// Printing the entered data for check
	fmt.Println("\n---")
	fmt.Printf("✅ Login data saved.\n")
	fmt.Printf("Email: %s\n", globalvars.Email)
	fmt.Printf("Password: %s [Length: %d]\n", globalvars.Password, len(globalvars.Password))
	fmt.Println("---")

	return models.LoginData{
		Email:    globalvars.Email,
		Password: globalvars.Password,
	}
}

func ReadStringFromConsole(message string) string {
	fmt.Print(message)
	var result string
	reader := bufio.NewReader(os.Stdin)
	// Reading a line from standard input
	result, _ = reader.ReadString('\n')
	// Remove any trailing newline characters
	result = strings.TrimSpace(result)
	return result
}

func AttachDefaultRequestHeaders(req *http.Request) {
	req.Header.Set("Content-Type", globalvars.ApplicationJson)
	req.Header.Set("User-Agent", globalvars.DefaultUserAgent)
	req.Header.Set("Accept-Encoding", globalvars.AcceptEncodingHeader)
	req.Header.Set("Origin", globalvars.Origin)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8,ru-RU;q=0.7")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Ch-Ua", `"Not A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("priority", "u=1, i")
}
