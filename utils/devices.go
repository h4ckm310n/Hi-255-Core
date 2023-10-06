package utils

var Devices map[string]*Device

type Device struct {
	DeviceID      string
	DeviceName    string
	DeviceAddr    string
	DeviceHTTPURL string
	Platform      string
}

var IsDevicesUpdated = false

func init() {
	Devices = make(map[string]*Device)
}

func AddDevice(deviceID, deviceName, deviceAddr, platform string) {
	device := GetDevice(deviceID)
	if device != nil {
		if device.DeviceID == deviceID && device.DeviceName == deviceName && device.DeviceAddr == deviceAddr && device.Platform == Platform {
			return
		}
		device.DeviceID = deviceID
		device.DeviceName = deviceName
		device.DeviceAddr = deviceAddr
		device.DeviceHTTPURL = "http://" + deviceAddr + ":2551"
		device.Platform = platform
		IsDevicesUpdated = true
		return
	}
	Devices[deviceID] = &Device{
		DeviceID:      deviceID,
		DeviceName:    deviceName,
		DeviceAddr:    deviceAddr,
		DeviceHTTPURL: "http://" + deviceAddr + ":2551",
		Platform:      platform,
	}
	IsDevicesUpdated = true
}

func DeleteDevice(deviceID string) {
	delete(Devices, deviceID)
	IsDevicesUpdated = true
}

func GetDevice(deviceID string) *Device {
	device, exist := Devices[deviceID]
	if exist {
		return device
	}
	return nil
}
