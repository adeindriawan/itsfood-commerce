package utils

import (
	"bytes"
	"net/http"
	"encoding/json"
)

func AddItsmineOrder(itsmineOrder map[string]interface{}) (any, error) {
	var response any
	url := "http://dev.itsmine.id/api/itsfood/handle-new-order"
	jsonified, _ := json.Marshal(itsmineOrder)
	payload := bytes.NewReader(jsonified)
	client := &http.Client{}
	
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	errorDecodingResponse := json.NewDecoder(res.Body).Decode(&response)
	if errorDecodingResponse != nil {
		return nil, err
	}

	return response, nil
}