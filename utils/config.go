package utils

import (
	"encoding/json"
	"github.com/google/uuid"
	"os"
	"runtime"
)

type ConfigType struct {
	DeviceID     string `json:"device_id"`
	DeviceName   string `json:"device_name"`
	DownloadPath string `json:"download_path"`
	KeepFileTime bool   `json:"keep_file_time"`
}

var Config = ConfigType{}
var ConfigFilePath string
var Platform string
var GRPCSocketPath string
var homeDir string
var hi255Path string

func init() {
	Platform = runtime.GOOS
	homeDir, _ = os.UserHomeDir()
	hi255Path = homeDir + "/.hi255/"
	_, err := os.Stat(hi255Path)
	if os.IsNotExist(err) {
		os.MkdirAll(hi255Path, 0755)
	}
	GRPCSocketPath = hi255Path + "hi255.sock"

}

func LoadConfig(path string) bool {
	ConfigFilePath = path
	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	err = json.Unmarshal(jsonBytes, &Config)
	return err == nil
}

func UpdateConfig(DeviceID string, DeviceName string, DownloadPath string, KeepFileTime bool) error {
	Config.DeviceID = DeviceID
	Config.DeviceName = DeviceName
	Config.DownloadPath = DownloadPath
	Config.KeepFileTime = KeepFileTime
	jsonBytes, err := json.Marshal(Config)
	if err != nil {
		return err
	}
	err = os.WriteFile(ConfigFilePath, jsonBytes, 0644)
	return err
}

func NewConfig() bool {
	Config.DownloadPath = homeDir + "/Hi255Downloads/"
	ConfigFilePath = hi255Path + "config.json"
	_, err := os.Stat(ConfigFilePath)

	// Config file exists
	if (err == nil || os.IsExist(err)) && LoadConfig(ConfigFilePath) {
		return true
	}

	Config.DeviceID = uuid.NewString()
	Config.DeviceName, _ = os.Hostname()
	Config.KeepFileTime = true

	_, err = os.Stat(hi255Path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(hi255Path, 0755)
		if err != nil {
			return false
		}
	}

	_, err = os.Stat(Config.DownloadPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(Config.DownloadPath, 0755)
		if err != nil {
			return false
		}
	}

	jsonBytes, err := json.Marshal(Config)
	if err != nil {
		return false
	}
	err = os.WriteFile(ConfigFilePath, jsonBytes, 0644)
	return err == nil
}
