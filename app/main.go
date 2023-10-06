package main

import (
	"Hi-255-Core/hi255_grpc"
	"Hi-255-Core/utils"
	"Hi-255-Core/utils/protocol"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configFilePath := flag.String("c", "", "")
	flag.Parse()
	if *configFilePath == "" {
		utils.NewConfig()
	} else {
		utils.LoadConfig(*configFilePath)
	}

	go handleSignals()
	go hi255_grpc.InitGRPCServer()
	go protocol.ListenUDP()
	protocol.ListenHTTP()

}

func handleSignals() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	utils.Log("Stopping")
	protocol.StopMulticast()
	protocol.StopHTTP()
	hi255_grpc.StopGRPC()
	os.Exit(0)
}
