package isyscol

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const URL = "http://api.messaging-service.com"

type SmsMessage struct {
	To   []string `json:"to"`
	Text string   `json:"text"`
	From string   `json:"from,omitempty"`
}

type SendSmsResponse struct {
	BulkId   *string                 `json:"bulkId"`
	Messages []SendSmsResponseDetail `json:"messages"`
}

type SendSmsResponseDetail struct {
	To        string `json:"to"`
	SmsCount  uint   `json:"smsCount"`
	MessageId string `json:"messageId"`
	Status    Status `json:"status"`
}

type Status struct {
	Id          *uint   `json:"id"`
	GroupId     *uint   `json:"groupId"`
	GroupName   *string `json:"groupName"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Action      *string `json:"action"`
}

type Sender struct {
	Username string
	Password string
}

func (s Sender) authString() string {
	authString := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", s.Username, s.Password)))
	return "Basic " + authString
}

func (s Sender) SendSms(message SmsMessage) (*SendSmsResponse, error) {
	var response SendSmsResponse
	var err error

	url := URL + "/sms/1/text/single"
	p, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(p))
	if err != nil {
		return nil, err
	}

	req.Header["Content-type"] = []string{"application/json"}
	req.Header["Accept"] = []string{"application/json"}
    req.Header["Authorization"] = []string{s.authString()}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

    defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
