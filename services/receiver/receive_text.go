package receiver

import (
	"Hi-255-Core/utils"
	"io"
	"net/http"
	"time"
)

func ReceiveText(w http.ResponseWriter, req *http.Request) {
	textBytes, _ := io.ReadAll(req.Body)
	text := string(textBytes)
	utils.MessageEnqueue(2, time.Now().Unix(), req.Header.Get("device"), text)
}
