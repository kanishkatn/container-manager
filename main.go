package main

import (
	"container-manager/handler"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	jrpcHandler := rpc.NewServer()
	jrpcHandler.RegisterCodec(json.NewCodec(), "application/json")
	err := jrpcHandler.RegisterService(new(handler.ContainerService), "")
	if err != nil {
		logrus.Fatalf("Failed to register service: %v", err)
	}
	http.Handle("/jrpc", jrpcHandler)

	logrus.Info("Starting server on port 8080")
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}
