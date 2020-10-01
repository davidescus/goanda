package goanda

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func unmarshalJson(body []byte, data interface{}) error {
	return json.Unmarshal(body, &data)
}

func createUrl(host string, endpoint string) string {
	var buffer bytes.Buffer
	// Generate the auth header
	buffer.WriteString(host)
	buffer.WriteString(endpoint)

	url := buffer.String()
	return url
}

func makeRequest(c *OandaConnection, endpoint string, client http.Client, req *http.Request) ([]byte, error) {
	req.Header.Set("User-Agent", c.headers.agent)
	req.Header.Set("Authorization", c.headers.auth)
	req.Header.Set("Content-Type", c.headers.contentType)

	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	bodyString := string(body[:])
	if strings.Contains(bodyString, "errorMessage") {
		return []byte{}, errors.New(fmt.Sprintf("OANDA API Error: %s On route: %s", bodyString, endpoint))
	}

	return body, nil
}
