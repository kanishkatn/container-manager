package main

import (
	"container-manager/handler"
	"container-manager/job"
	"container-manager/p2p"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	jobQueue := job.NewQueue(10)
	jobQueue.Run(2)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := "0.0.0.0:" + port

	// setup p2p service
	logrus.Infof("Starting P2P service")
	p2pService, err := p2p.NewP2PService(jobQueue)
	if err != nil {
		logrus.Fatalf("Failed to create P2P service: %v", err)
	}
	p2pService.Start()

	jrpcHandler := rpc.NewServer()
	jrpcHandler.RegisterCodec(json.NewCodec(), "application/json")
	err = jrpcHandler.RegisterService(handler.NewContainerService(jobQueue, p2pService), "")
	if err != nil {
		logrus.Fatalf("Failed to register service: %v", err)
	}
	http.Handle("/jrpc", jrpcHandler)

	logrus.Infof("JRPC server listening on port %s", port)
	logrus.Fatal(http.ListenAndServe(address, nil))
}
