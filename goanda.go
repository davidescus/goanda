package goanda

import (
	"bytes"
	"net/http"
	"time"
)

type Headers struct {
	contentType    string
	agent          string
	DatetimeFormat string
	auth           string
}

type Connection interface {
	Get(endpoint string) ([]byte, error)
	Post(endpoint string, data []byte) ([]byte, error)
	Put(endpoint string, data []byte) ([]byte, error)
	GetOrderDetails(instrument string, units string) OrderDetails
	GetAccountSummary() AccountSummary
	CreateOrder(body OrderPayload) OrderResponse
}

type OandaConnection struct {
	hostname       string
	port           int
	ssl            bool
	token          string
	accountID      string
	DatetimeFormat string
	headers        *Headers
}

const OANDA_AGENT string = "v20-golang/0.0.1"

func NewConnection(accountID string, token string, live bool) *OandaConnection {
	hostname := ""
	// should we use the live API?
	if live {
		hostname = "https://api-fxtrade.oanda.com/v3"
	} else {
		hostname = "https://api-fxpractice.oanda.com/v3"
	}

	var buffer bytes.Buffer
	// Generate the auth header
	buffer.WriteString("Bearer ")
	buffer.WriteString(token)

	authHeader := buffer.String()
	// Create headers for oanda to be used in requests
	headers := &Headers{
		contentType:    "application/json",
		agent:          OANDA_AGENT,
		DatetimeFormat: "RFC3339",
		auth:           authHeader,
	}
	// Create the connection object
	connection := &OandaConnection{
		hostname:  hostname,
		port:      443,
		ssl:       true,
		token:     token,
		headers:   headers,
		accountID: accountID,
	}

	return connection
}

// TODO refactor methods get, post, put, they looks similar

// TODO: include params as a second option
func (c *OandaConnection) Get(endpoint string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 5, // 5 sec timeout
	}

	url := createUrl(c.hostname, endpoint)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, err
	}

	body, err := makeRequest(c, endpoint, client, req)
	if err != nil {
		return []byte{}, err
	}

	return body, err
}

func (c *OandaConnection) Post(endpoint string, data []byte) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 5, // 5 sec timeout
	}

	url := createUrl(c.hostname, endpoint)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return []byte{}, err
	}

	body, err := makeRequest(c, endpoint, client, req)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func (c *OandaConnection) Put(endpoint string, data []byte) ([]byte, error) {
	client := http.Client{Timeout: time.Second * 5}
	url := createUrl(c.hostname, endpoint)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return []byte{}, err
	}

	body, err := makeRequest(c, endpoint, client, req)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
