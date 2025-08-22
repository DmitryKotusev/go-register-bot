package main

import (
	"bot-main/requests"
	"bot-main/utils"
	"fmt"
)

func main() {
	// Look at ability to use https://github.com/fatih/color
	fmt.Println("Starting the bot, press Ctrl+C to stop it at any time.")
	fmt.Println("Reading input data...")
	utils.RegisterCommandLineArgs()
	applicationData := utils.ReadRequiredApplicationData()
	requests.RequestPipeline(applicationData)
}
