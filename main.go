package main

import (
	"container-manager/handler"
	"container-manager/p2p"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// setup p2p service
	p2pService, err := p2p.NewP2PService()
	if err != nil {
		logrus.Fatalf("Failed to create P2P service: %v", err)
	}
	p2pService.Start()

	jrpcHandler := rpc.NewServer()
	jrpcHandler.RegisterCodec(json.NewCodec(), "application/json")
	err = jrpcHandler.RegisterService(new(handler.ContainerService), "")
	if err != nil {
		logrus.Fatalf("Failed to register service: %v", err)
	}
	http.Handle("/jrpc", jrpcHandler)

	logrus.Info("Starting jrpc server on port 8080")
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}
