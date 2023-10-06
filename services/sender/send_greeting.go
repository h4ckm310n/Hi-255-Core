package sender

import (
	"Hi-255-Core/services/models"
	"Hi-255-Core/utils"
	"bytes"
	"encoding/json"
	"net/http"
)

func SendGreeting(addr string) (*models.Greeting, error) {
	greeting := models.Greeting{
		DeviceID:   utils.Config.DeviceID,
		DeviceName: utils.Config.DeviceName,
		Platform:   utils.Platform,
	}
	jsonData, err := json.Marshal(greeting)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://"+addr+":2551"+"/greeting", bytes.NewBuffer(jsonData))
	if err != nil {
		utils.Err(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		utils.Err(err)
		return nil, err
	}

	var respGreeting models.Greeting
	err = json.NewDecoder(resp.Body).Decode(&respGreeting)
	if err != nil {
		return nil, err
	}
	utils.AddDevice(respGreeting.DeviceID, respGreeting.DeviceName, addr, respGreeting.Platform)
	return &respGreeting, nil
}
