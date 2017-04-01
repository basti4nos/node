package ipify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/cihub/seelog"
)

const IPIFY_API_URL = "https://api.ipify.org/"
const IPIFY_API_CLIENT = "goclient-v0.1"
const IPIFY_API_LOG_PREFIX = "[ipify.api] "

func NewClient() *client {
	return &client{
		httpClient: http.Client{},
	}
}

type client struct {
	httpClient http.Client
}

func (client *client) GetIp() (string, error) {
	var ipResponse IpResponse

	request, err := http.NewRequest("GET", IPIFY_API_URL+"/?format=json", nil)
	request.Header.Set("User-Agent", IPIFY_API_CLIENT)
	request.Header.Set("Accept", "application/json")
	if err != nil {
		log.Critical(IPIFY_API_LOG_PREFIX, err)
		return "", err
	}

	err = client.doRequest(request, &ipResponse)
	if err != nil {
		return "", err
	}

	log.Info(IPIFY_API_LOG_PREFIX, "IP detected: ", ipResponse.IP)
	return ipResponse.IP, nil
}

func (client *client) doRequest(request *http.Request, responseDto interface{}) error {
	response, err := client.httpClient.Do(request)
	defer response.Body.Close()
	if err != nil {
		log.Error(IPIFY_API_LOG_PREFIX, err)
		return err
	}

	err = parseResponseError(response)
	if err != nil {
		log.Error(IPIFY_API_LOG_PREFIX, err)
		return err
	}

	return parseResponseJson(response, &responseDto)
}

func parseResponseJson(response *http.Response, dto interface{}) error {
	responseJson, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseJson, dto)
	if err != nil {
		return err
	}

	return nil
}

func parseResponseError(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Server response invalid: %s", response.Status))
	}

	return nil
}
