package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// WhatsAppClient handles sending WhatsApp messages via an API
type WhatsAppClient struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

// NewWhatsAppClient creates a new instance of WhatsAppClient
func NewWhatsAppClient(accountSID, authToken, fromNumber string) *WhatsAppClient {
	return &WhatsAppClient{
		AccountSID: accountSID,
		AuthToken:  authToken,
		FromNumber: fromNumber,
	}
}

// SendMessage sends a WhatsApp message to the specified number
func (client *WhatsAppClient) SendMessage(to string, message string) error {
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + client.AccountSID + "/Messages.json"
	msgData := url.Values{}
	msgData.Set("To", "whatsapp:"+to)
	msgData.Set("From", "whatsapp:"+client.FromNumber)
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(client.AccountSID, client.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	clientHTTP := &http.Client{}
	resp, _ := clientHTTP.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	} else {
		return fmt.Errorf("Failed to send WhatsApp message: %s", resp.Status)
	}
}
