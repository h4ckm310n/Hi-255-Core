package protocol

import (
	"Hi-255-Core/services/receiver"
	"Hi-255-Core/utils"
	"net/http"
)

var server *http.Server

func ListenHTTP() {
	http.HandleFunc("/greeting", receiver.ReceiveGreeting)
	http.HandleFunc("/fileinfo", receiver.ReceiveFileInfo)
	http.HandleFunc("/file", receiver.ReceiveFile)
	http.HandleFunc("/text", receiver.ReceiveText)
	server = &http.Server{Addr: ":2551"}
	server.ListenAndServe()
}

func StopHTTP() {
	server.Close()
	utils.Log("Stopped HTTP")
}
