package main

import (
	"bot-main/requests"
	"bot-main/utils"
	"fmt"
)

func main() {
	fmt.Println("Starting the bot, press Ctrl+C to stop it at any time.")
	fmt.Println("Reading input data...")
	utils.RegisterCommandLineArgs()
	applicationData := utils.ReadRequiredApplicationData()
	requests.RequestPipeline(applicationData)
}
