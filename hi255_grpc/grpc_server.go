package hi255_grpc

import (
	"Hi-255-Core/services/sender"
	"Hi-255-Core/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

type GRPCServer struct {
	UnimplementedServiceServer
}

var server *grpc.Server

func (s *GRPCServer) FetchRemoteDevices(req *Empty, stream Service_FetchRemoteDevicesServer) error {
	for {
		if !utils.IsDevicesUpdated {
			time.Sleep(10 * time.Second)
			continue
		}
		var resp []*RemoteDevicesResponse_RemoteDeviceItem
		for deviceID, device := range utils.Devices {
			resp = append(resp, &RemoteDevicesResponse_RemoteDeviceItem{
				Id:       deviceID,
				Name:     device.DeviceName,
				Address:  device.DeviceAddr,
				Platform: device.Platform,
			})
		}
		if err := stream.Send(&RemoteDevicesResponse{
			RemoteDevices: resp,
		}); err != nil {
			utils.Err("Fetch devices: stop unexpectedly")
			utils.Err(err)
			return err
		}
		utils.IsDevicesUpdated = false
		time.Sleep(10 * time.Second)
	}
	return nil
}

func (s *GRPCServer) FetchMessages(req *Empty, stream Service_FetchMessagesServer) error {
	for {
		if len(utils.MessageQueue) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		var resp []*MessagesResponse_MessageItem
		for {
			item := utils.MessageDequeue()
			if item == nil {
				break
			}
			resp = append(resp, &MessagesResponse_MessageItem{
				MessageType: item.MessageType,
				Timestamp:   item.Timestamp,
				RemoteId:    item.DeviceID,
				Content:     item.Content,
			})
		}
		if err := stream.Send(&MessagesResponse{
			Messages: resp,
		}); err != nil {
			utils.Err("Fetch messages: stopped unexpectedly")
			utils.Err(err)
			return err
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (s *GRPCServer) UpdateConfig(ctx context.Context, req *UpdateConfigRequest) (*Empty, error) {
	err := utils.UpdateConfig(req.DeviceId, req.DeviceName, req.DownloadPath, req.KeepFileTime)
	return &Empty{}, err
}

func (s *GRPCServer) SendGreeting(ctx context.Context, req *SendGreetingRequest) (*CommonResponse, error) {
	_, err := sender.SendGreeting(req.RemoteAddress)
	var status int32 = 0
	if err != nil {
		status = -1
	}
	return &CommonResponse{
		Status: status,
	}, err
}

func (s *GRPCServer) SendFile(ctx context.Context, req *SendFileRequest) (*CommonResponse, error) {
	err := sender.SendFile(req.FilePath, req.RemoteId)
	var status int32 = 0
	if err != nil {
		status = -1
	}
	return &CommonResponse{
		Status: status,
	}, err
}

func (s *GRPCServer) SendText(ctx context.Context, req *SendTextRequest) (*CommonResponse, error) {
	err := sender.SendText(req.Text, req.RemoteId)
	var status int32 = 0
	if err != nil {
		status = -1
	}
	return &CommonResponse{
		Status: status,
	}, err
}

func InitGRPCServer() {
	os.Remove(utils.GRPCSocketPath)
	conn, err := net.Listen("unix", utils.GRPCSocketPath)
	if err != nil {
		utils.Err(err)
		return
	}
	server = grpc.NewServer()
	RegisterServiceServer(server, &GRPCServer{})
	utils.Log("Start GRPC Server")
	err = server.Serve(conn)
	if err != nil {
		utils.Err(err)
		return
	}
}

func StopGRPC() {
	server.GracefulStop()
	utils.Log("Stopped GRPC")
}
