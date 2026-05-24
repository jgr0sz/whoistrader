package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

//Generic function for processing GET requests from various CS2 APIs.
func GetAPI(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Malformed request: %s", err)
		return nil, err
	}

	//Auth tokens, verification, etc...
	for k, v := range headers {
		request.Header.Add(k, v)
	}

	//Reading response input in
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	//Error response handling to prevent silent marshalling errors.
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d", response.StatusCode)
	}
	return body, nil
}
