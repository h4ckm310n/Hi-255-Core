package receiver

import (
	"Hi-255-Core/services/models"
	"Hi-255-Core/utils"
	"encoding/json"
	"net"
	"net/http"
)

func ReceiveGreeting(w http.ResponseWriter, req *http.Request) {
	var greeting models.Greeting
	json.NewDecoder(req.Body).Decode(&greeting)
	respGreeting := models.Greeting{
		DeviceID:   utils.Config.DeviceID,
		DeviceName: utils.Config.DeviceName,
		Platform:   utils.Platform,
	}
	data, err := json.Marshal(respGreeting)
	if err != nil {
		return
	}
	_, err = w.Write(data)
	if err != nil {
		return
	}
	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	utils.Log(req.RemoteAddr, greeting.DeviceName)
	utils.Log(req.Header.Get("X-Real-IP"), req.Header.Get("X-Forwarded-For"))
	utils.AddDevice(greeting.DeviceID, greeting.DeviceName, ip, greeting.Platform)
}
