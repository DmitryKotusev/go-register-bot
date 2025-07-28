package main

import (
	"bot-main/requests"
	"bot-main/utils"
	"fmt"
)

func main() {
	fmt.Println("Starting the bot, press Ctrl+C to stop it at any time.")
	fmt.Println("Reading input data...")
	loginData := utils.ReadRequiredLoginData()
	requests.RequestPipeline(loginData)
}
