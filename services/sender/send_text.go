package sender

import (
	"Hi-255-Core/utils"
	"bytes"
	"errors"
	"net/http"
)

func SendText(text string, remoteID string) error {
	remoteDevice := utils.GetDevice(remoteID)
	if remoteDevice == nil {
		return errors.New("device not exist")
	}
	addr := remoteDevice.DeviceHTTPURL

	client := &http.Client{}
	req, err := http.NewRequest("POST", addr+"/text", bytes.NewBufferString(text))
	if err != nil {
		utils.Err(err)
		return err
	}
	req.Header.Set("device", utils.Config.DeviceID)
	_, err = client.Do(req)
	if err != nil {
		utils.Err(err)
	}
	return err
}
