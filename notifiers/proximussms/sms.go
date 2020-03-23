package proximussms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mhereman/cryptotrader/interfaces"
	"github.com/mhereman/cryptotrader/notifiers"

	"github.com/mhereman/cryptotrader/logger"
)

const notifierName = "proximus-sms"

func init() {
	notifiers.RegisterNotifier(notifierName, createProximusSms)
}

type ProximusSMS struct {
	apiToken    string
	destination string
}

type message struct {
	Message      string   `json:"message"`
	Binary       bool     `json:"binary"`
	Destinations []string `json:"destinations"`
}

type errorResponse struct {
	Fault faultInfo `json:"fault"`
}

type faultInfo struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:description"`
}

type successResponse struct {
	ResourceURL  string         `json:"resourceURL"`
	DeliveryInfo []deliveryInfo `json:"deliveryInfo"`
}

type deliveryInfo struct {
	Address        string `json:"address"`
	DeliveryStatus string `json:"deliveryStatus"`
}

func New(ctx context.Context, config map[string]string) (sms *ProximusSMS, err error) {
	var apiToken, destination string
	var ok bool

	if apiToken, ok = config["apiToken"]; !ok {
		err = fmt.Errorf("ProximusSMS config error: 'apiToken' entry not found")
		return
	}

	if destination, ok = config["destination"]; !ok {
		err = fmt.Errorf("Proximus config error: 'destination' entry not found")
		return
	}

	sms = new(ProximusSMS)
	sms.apiToken = apiToken
	sms.destination = destination
	return
}

func createProximusSms(ctx context.Context, config map[string]string) (notifier interfaces.INotifier, err error) {
	notifier, err = New(ctx, config)
	return
}

func (sms ProximusSMS) Name() string {
	return notifierName
}

func (sms ProximusSMS) Notify(ctx context.Context, data []byte) (err error) {
	var smsMessage message
	var client *http.Client
	var url string
	var req *http.Request
	var resp *http.Response
	var jsonData []byte
	var okResp successResponse
	var errResp errorResponse
	var deliveryStatus string

	smsMessage = message{
		Message:      string(data),
		Binary:       false,
		Destinations: []string{sms.destination},
	}
	if jsonData, err = json.Marshal(smsMessage); err != nil {
		logger.Errorf("ProximusSMS::Notify Error %v\n", err)
		return
	}

	client = &http.Client{}
	url = "https://api.enco.io/sms/1.0.0/sms/outboundmessages"
	if req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData)); err != nil {
		logger.Errorf("ProximusSMS::Notify Error: %v\n", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sms.apiToken))

	if resp, err = client.Do(req); err != nil {
		logger.Errorf("ProximusSMS::Notify Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		if err = json.Unmarshal(body, &errResp); err != nil {
			logger.Errorf("ProximusSMS::Notify Error %v\n", err)
			return
		}

		err = fmt.Errorf("API Error: Code: %d Message: %s Details: %s", errResp.Fault.Code, errResp.Fault.Message, errResp.Fault.Description)
		logger.Errorf("ProximusSMS::Notify Error %v\n", err)
		return
	}

	if err = json.Unmarshal(body, &okResp); err != nil {
		logger.Errorf("ProximusSMS::Notify Error %v\n", err)
		return
	}

	if len(okResp.DeliveryInfo) == 0 {
		err = fmt.Errorf("API Error No delivery information")
		logger.Errorf("ProximusSMS::Notify Error %v\n", err)
		return
	}

	deliveryStatus = okResp.DeliveryInfo[0].DeliveryStatus
	if deliveryStatus == "AddressInvalid" || deliveryStatus == "DeliveryImpossible" {
		logger.Warningf("ProsimusSMS::Notify Warning API Warning Delivery: %s\n", deliveryStatus)
	}
	return
}
