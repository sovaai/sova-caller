package api

import (
	"fmt"
	"sova-caller-backend/utils"
	"time"

	"github.com/go-resty/resty/v2"
)

/* TYPES */
type RequestConfig struct {
	URL         string `json:"url"`
	METHOD      string `json:"method,omitempty"`
	HTTP_METHOD string `json:"http_method"`
	PAYLOAD     []byte `json:"payload"`
	TOKEN       string `json:"token,omitempty"`
}

/* TYPES */

var api_client = resty.New()

func CreateRequestConfig(data []byte, URL string, API_METHOD string, HTTP_METHOD string) RequestConfig {
	config := RequestConfig{}

	config.PAYLOAD = data
	config.URL = URL
	config.METHOD = API_METHOD
	config.HTTP_METHOD = HTTP_METHOD

	return config
}

func Request(config RequestConfig) ([]byte, error) {
	defer utils.TimeTrack(time.Now(), "REQUEST")

	var URL = utils.If(config.METHOD == "", config.URL, fmt.Sprintf("%s/%s", config.URL, config.METHOD))

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	if config.TOKEN != "" {
		headers["Authorization"] = config.TOKEN
	}

	res, err := api_client.R().
		SetHeaders(headers).
		SetBody(config.PAYLOAD).
		Execute(config.HTTP_METHOD, URL)

	if err != nil {
		return nil, err
	}

	result := res.Body()

	return result, nil
}
