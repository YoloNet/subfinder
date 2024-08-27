package main

import (
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	_ "github.com/projectdiscovery/fdmax/autofdmax"
	"github.com/projectdiscovery/gologger"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"fmt"
)

// Function to read file content
func readFromFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil // Trim whitespace for clean reading
}

// Function to send a message to the Telegram bot
func sendTelegramMessage(botToken string, chatID string, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	// Parse the command line flags and read config files
	options := runner.ParseOptions()

	newRunner, err := runner.NewRunner(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %s\n", err)
	}

	// Run the enumeration
	err = newRunner.RunEnumeration()
	if err != nil {
		gologger.Fatal().Msgf("Could not run enumeration: %s\n", err)
	}

	// After running the enumeration, gather results and send them to Telegram
	results := newRunner.GetResults() // Assuming GetResults() gives the output

	// Read the Telegram bot token and chat ID from files
	botTokenFilePath := "/root/yoloautosc/tele_bot_token"
	chatIDFilePath := "/root/yoloautosc/tele_chat_id"

	telegramBotToken, err := readFromFile(botTokenFilePath)
	if err != nil {
		gologger.Error().Msgf("Failed to read Telegram bot token from file: %s\n", err)
		return
	}

	telegramChatID, err := readFromFile(chatIDFilePath)
	if err != nil {
		gologger.Error().Msgf("Failed to read Telegram chat ID from file: %s\n", err)
		return
	}

	// Format the message to send
	message := fmt.Sprintf("Subfinder results:\n%s", results)

	// Send the results to the Telegram bot
	err = sendTelegramMessage(telegramBotToken, telegramChatID, message)
	if err != nil {
		gologger.Error().Msgf("Failed to send Telegram message: %s\n", err)
	} else {
		gologger.Info().Msg("Results sent successfully to Telegram")
	}
}
