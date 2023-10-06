package receiver

import (
	"Hi-255-Core/services/models"
	"Hi-255-Core/utils"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReceiveFile(w http.ResponseWriter, req *http.Request) {
	fileinfo := utils.GetSessionValue(req.Header.Get("session"))
	defer utils.DeleteSession(req.Header.Get("session"))

	// Checksum
	fileHash := sha256.New()
	io.Copy(fileHash, req.Body)
	if string(fileHash.Sum(nil)) != req.Header.Get("checksum") {
		w.Write([]byte("-1"))
		return
	}

	filepath := renameFile(utils.Config.DownloadPath, fileinfo.(models.FileInfo).Filename)
	file, _ := os.Create(filepath)
	defer file.Close()
	size, _ := io.Copy(file, req.Body)

	w.Write([]byte(strconv.FormatInt(size, 10)))

	// Set file original time
	if fileinfo.(models.FileInfo).FileTime != 0 && utils.Config.KeepFileTime {
		fileTime := time.Unix(fileinfo.(models.FileInfo).FileTime, 0)
		os.Chtimes(filepath, fileTime, fileTime)
	}
	utils.MessageEnqueue(1, time.Now().Unix(), req.Header.Get("device"), filepath)
}

func renameFile(path, filename string) string {
	_, err := os.Stat(path + filename)
	if os.IsNotExist(err) {
		return path + filename
	}

	extPos := strings.LastIndex(filename, ".")
	ext := ""
	basename := filename
	if !(extPos == -1 || extPos == len(filename)-1 || extPos == 0) {
		ext = filename[extPos:]
		basename = filename[:extPos]
	}

	n := 1
	newPath := ""
	for {
		newPath = fmt.Sprintf("%s%s (%d)%s", path, basename, n, ext)
		_, err = os.Stat(newPath)
		if os.IsNotExist(err) {
			break
		}
		n += 1
	}
	return newPath
}

func ReceiveFileInfo(w http.ResponseWriter, req *http.Request) {
	var info models.FileInfo
	json.NewDecoder(req.Body).Decode(&info)

	session := utils.GenerateSession()
	utils.SetSession(string(session), info)
	w.Write(session)
}
