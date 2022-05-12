package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type logRequestBody struct {
	Name          string `json:"name"`
	Url           string `json:"url"`
	Domain        string `json:"domain"`
	RemoteAddress string `json:"-"`
	path          string `json:"-"`
}

func NewLogRequestBody(path string, address string, domain string) *logRequestBody {
	return &logRequestBody{
		Name:          "pageview",
		Url:           fmt.Sprintf("https://%s%s", domain, path),
		Domain:        domain,
		RemoteAddress: address,
		path:          path,
	}
}

func LogRequestToPlausible(lrb *logRequestBody, statsApiUrl string) {

	if statsApiUrl != "" {

		postBody, err := json.Marshal(lrb)
		if err != nil {
			log.Println(err.Error())
		}

		responseBody := bytes.NewBuffer(postBody)

		//Leverage Go's HTTP Post function to make request

		request, err := http.NewRequest("POST", statsApiUrl, responseBody)
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		request.Header.Set("User-Agent", "API")
		request.Header.Set("X-Forwarded-For", lrb.RemoteAddress)

		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer response.Body.Close()

	}

}
