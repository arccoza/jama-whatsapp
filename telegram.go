package main

import (
	"os"
	"fmt"
	"github.com/arccoza/go-tdlib"
)


func main() {
	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")

	// Create new instance of client
	client := tdlib.NewClient(tdlib.Config{
		APIID:               os.Getenv("TELEGRAM_API_ID"),
		APIHash:             os.Getenv("TELEGRAM_API_HASH"),
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  false,
		UseFileDatabase:     false,
		UseChatInfoDatabase: false,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
	})

	for {
		fmt.Println("a")
		currentState, _ := client.Authorize()
		fmt.Println("b", currentState)
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Println("1")
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			state, err := client.SendPhoneNumber(number)
			fmt.Println("SendPhoneNumber state: ", state)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Println("2")
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			state, err := client.SendAuthCode(code)
			fmt.Println("SendAuthCode state: ", state)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Println("3")
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			fmt.Println("4")
			fmt.Println("Authorization Ready! Let's rock")
			break
		}
	}

	// Main loop
	for update := range client.GetRawUpdatesChannel(10) {
		fmt.Println("5")
		// Show all updates
		fmt.Println(update.Data)
		fmt.Print("\n\n")
	}

}
