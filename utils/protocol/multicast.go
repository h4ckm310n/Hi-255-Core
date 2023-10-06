package protocol

import (
	"Hi-255-Core/utils"
	"encoding/json"
	"net"
	"time"
)

var (
	multicastConn     *net.UDPConn
	multicastUDPAddr  *net.UDPAddr
	multicastBodyJson []byte
	udpClient         *net.UDPConn
	exitSignal        bool
)

const (
	MULTICAST_ADDR = "224.0.25.5:2550"
	HEAD           = "0x10a9fc70042"

	TYPE_HELLO = "1"
	TYPE_REPLY = "2"
	TYPE_BYE   = "3"
	LEN_HEAD   = len(HEAD)
)

func ListenUDP() {
	exitSignal = false
	multicastUDPAddr, _ = net.ResolveUDPAddr("udp", MULTICAST_ADDR)
	multicastConn, _ = net.ListenMulticastUDP("udp", nil, multicastUDPAddr)
	defer multicastConn.Close()

	initMulticastBody()
	go sendMultiCast()

	for {
		data := make([]byte, 1024)
		n, remoteAddr, err := multicastConn.ReadFromUDP(data)
		if err != nil {
			if exitSignal {
				break
			}
			utils.Err(err)
			continue
		}
		if n < LEN_HEAD || string(data[:LEN_HEAD]) != HEAD {
			continue
		}
		go handleUDPConn(n, data, remoteAddr)
	}
}

func handleUDPConn(size int, data []byte, remoteAddr *net.UDPAddr) {
	typeMsg := string(data[LEN_HEAD : LEN_HEAD+1])
	switch typeMsg {
	case TYPE_BYE:
		// Remote termination
		deviceID := string(data[LEN_HEAD+1 : size])
		utils.DeleteDevice(deviceID)
		break
	case TYPE_HELLO:
		// Remote startup
		var remoteDevice map[string]string
		json.Unmarshal(data[LEN_HEAD+1:size], &remoteDevice)
		if remoteDevice["DeviceID"] != utils.Config.DeviceID && utils.GetDevice(remoteDevice["DeviceID"]) == nil {
			utils.Log(remoteDevice)
			utils.AddDevice(remoteDevice["DeviceID"], remoteDevice["DeviceName"], remoteAddr.IP.String(), remoteDevice["Platform"])
			multicastConn.WriteToUDP(append([]byte(HEAD+TYPE_REPLY), multicastBodyJson...), remoteAddr)
		}
		break
	case TYPE_REPLY:
		// Remote reply
		var remoteDevice map[string]string
		json.Unmarshal(data[LEN_HEAD+1:size], &remoteDevice)
		utils.AddDevice(remoteDevice["DeviceID"], remoteDevice["DeviceName"], remoteAddr.IP.String(), remoteDevice["Platform"])
		break
	default:
		break
	}
}

func sendMultiCast() {
	body := append([]byte(HEAD+TYPE_HELLO), multicastBodyJson...)
	var err error
	udpClient, err = net.DialUDP("udp", nil, multicastUDPAddr)
	if err != nil {
		utils.Err(err)
		return
	}

	listenReply := func() {
		for {
			data := make([]byte, 1024)
			n, remoteAddr, err := udpClient.ReadFromUDP(data)
			if err != nil {
				if exitSignal {
					break
				}
				utils.Err(err)
				continue
			}
			if n < LEN_HEAD+1 || string(data[:LEN_HEAD+1]) != HEAD+TYPE_REPLY {
				continue
			}
			var remoteDevice map[string]string
			json.Unmarshal(data[LEN_HEAD+1:], &remoteDevice)
			utils.AddDevice(remoteDevice["DeviceID"], remoteDevice["DeviceName"], remoteAddr.IP.String(), remoteDevice["Platform"])
		}
	}

	go listenReply()
	for {
		_, err = udpClient.Write(body)
		if err != nil {
			if exitSignal {
				break
			}
			utils.Err(err)
			continue
		}
		time.Sleep(5 * time.Second)
	}
}

func initMulticastBody() {
	multicastBody := map[string]string{
		"DeviceID":   utils.Config.DeviceID,
		"DeviceName": utils.Config.DeviceName,
		"Platform":   utils.Platform,
	}
	multicastBodyJson, _ = json.Marshal(multicastBody)
}

func StopMulticast() {
	exitSignal = true
	udpClient.WriteTo([]byte(HEAD+TYPE_BYE+utils.Config.DeviceID), multicastUDPAddr)
	udpClient.Close()
	utils.Log("Stopped Multicast")
}
