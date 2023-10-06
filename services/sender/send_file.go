package sender

import (
	"Hi-255-Core/services/models"
	"Hi-255-Core/utils"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func SendFile(filepath string, remoteID string) error {
	remoteDevice := utils.GetDevice(remoteID)
	if remoteDevice == nil {
		return errors.New("device not exist")
	}
	addr := remoteDevice.DeviceHTTPURL

	client := &http.Client{}
	session, err := sendFileInfo(filepath, addr, client)
	if err != nil {
		return err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	// Checksum
	fileHash := sha256.New()
	if _, err = io.Copy(fileHash, file); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", addr+"/file", file)
	if err != nil {
		return err
	}
	req.Header.Set("session", session)
	req.Header.Set("device", utils.Config.DeviceID)
	req.Header.Set("hash", string(fileHash.Sum(nil)))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	size, err := io.ReadAll(resp.Body)
	if string(size) == "-1" {
		return errors.New("checksum failed")
	}
	return err
}

func sendFileInfo(filepath string, addr string, client *http.Client) (string, error) {
	info, err := os.Stat(filepath)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	fileTime := info.ModTime().Unix()
	size := info.Size()
	filename := info.Name()
	data := models.FileInfo{
		Filename: filename,
		FileTime: fileTime,
		Size:     size,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", addr+"/fileinfo", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	session, err := io.ReadAll(resp.Body)
	return string(session), err
}
