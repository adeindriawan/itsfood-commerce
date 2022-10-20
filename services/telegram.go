package services

import (
	"bytes"
	"encoding/json"
	"os"
	"log"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	Token  string
	ChatId string
)

func getUrl() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", Token)
}

func SendTelegramToGroup(text string) (bool, error) {
	// Global variables
	var err error
	var response *http.Response

	ChatId = os.Getenv("TELEGRAM_CHAT_ID")
	Token = os.Getenv("TELEGRAM_BOT_TOKEN")

	// Send the message
	url := fmt.Sprintf("%s/sendMessage", getUrl())
	body, _ := json.Marshal(map[string]string{
		"chat_id": ChatId,
		"text":    text,
	})
	response, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return false, err
	}

	// Close the request at the end
	defer response.Body.Close()

	// Body
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	// Log
	log.Print("Message: ", text, " was sent")
	log.Print("Response JSON: ", string(body))

	// Return
	return true, nil
}

func SendTelegramToVendor(text string, chatId string) (bool, error) {
	// Global variables
	var err error
	var response *http.Response

	ChatId = chatId
	Token = os.Getenv("TELEGRAM_BOT_TOKEN")

	// Send the message
	url := fmt.Sprintf("%s/sendMessage", getUrl())
	body, _ := json.Marshal(map[string]string{
		"chat_id": ChatId,
		"text":    text,
	})
	response, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return false, err
	}

	// Close the request at the end
	defer response.Body.Close()

	// Body
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	// Log
	log.Print("Message: ", text, " was sent")
	log.Print("Response JSON: ", string(body))

	// Return
	return true, nil
}