package httprequest

import (
	"bytes"
	"encoding/json"
	"io"

	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

func BodyDumpHandlerfunc(url, req string, res []byte) {

	log.Printf("URL: %s, Request: %s, Response: %v", url, req, string(res[:]))

}

func PostJson(url string, param interface{}, timeout time.Duration, token string) (body []byte, statusCode int, err error) {

	byteR, _ := json.Marshal(param)
	netClient := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(byteR))
	if err != nil {
		err = errors.Wrapf(err, "[pkg.httprequest] PostJson.NewRequest url:%s request:%+v", url, param)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-access-token", token)

	resp, err := netClient.Do(req)
	if err != nil {
		err = errors.Wrapf(err, "[pkg.httprequest] PostJson.Do url:%s request:%+v", url, param)
		statusCode = http.StatusBadGateway
		return
	}

	defer resp.Body.Close()

	statusCode = resp.StatusCode
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrapf(err, "[pkg.httprequest] PostJson.ReadAll url:%s request:%+v statusCode:%v", url, param, resp.StatusCode)
	}

	return
}

// SendGetRequest makes an authenticated GET request and returns the response.
func SendGetRequest(url string, param interface{}, timeout time.Duration, token string) (response *http.Response, statusCode int, err error) {

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "[pkg.httprequest] Failed to create request. url:%s", url)
	}

	// Set the authorization header
	req.Header.Set("x-access-token", token)

	// Make the request
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway, errors.Wrapf(err, "[pkg.httprequest] Request failed. url:%s", url)
	}

	// Return response even if status is not 200, so caller can handle it
	return resp, resp.StatusCode, nil
}
