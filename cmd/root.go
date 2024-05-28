package cmd

import (
	cfg "container-manager/config"
	"container-manager/handler"
	"container-manager/services"
	"fmt"
	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initialize the default config
var config = cfg.DefaultConfig()

// rootCmd represents the base container-manager command
var rootCmd = &cobra.Command{
	Use:   "container-manager",
	Short: "container-manager is a tool to manage containers",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}

		// set log level
		level, err := logrus.ParseLevel(config.LogLevel)
		if err != nil {
			return fmt.Errorf("failed to parse log level: %w", err)
		}
		logrus.SetLevel(level)

		if err := runNode(); err != nil {
			return fmt.Errorf("failed to run node: %w", err)
		}

		return nil
	},
}

// init initializes the flags for the root command
func init() {
	rootCmd.Flags().StringVar(&config.LogLevel, "log-level", config.LogLevel, "log level")
	rootCmd.Flags().IntVar(&config.QueueSize, "queue-size", config.QueueSize, "the size of the job queue")
	rootCmd.Flags().IntVar(&config.WorkerCount, "worker-count", config.WorkerCount, "the number of workers to run")
	rootCmd.Flags().StringVar(&config.ListenAddress, "listen-address", config.ListenAddress, "the address to listen on")
	rootCmd.Flags().StringVar(&config.Port, "port", config.Port, "the port to listen on")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("failed to execute command: %v", err)
	}
}

// runNode runs the container manager node
func runNode() error {
	ds, err := services.NewDockerService()
	if err != nil {
		return fmt.Errorf("failed to create docker service: %w", err)
	}

	jobQueue := services.NewQueue(config.QueueSize, ds)
	jobQueue.Run(config.WorkerCount)

	// setup p2p service
	logrus.Infof("Starting P2P service")
	p2pService, err := services.NewP2PService(jobQueue)
	if err != nil {
		return fmt.Errorf("failed to create P2P service: %w", err)
	}
	p2pService.Start()

	// setup jrpc handler
	jrpcHandler := rpc.NewServer()
	jrpcHandler.RegisterCodec(json.NewCodec(), "application/json")
	err = jrpcHandler.RegisterService(handler.NewContainerService(jobQueue, p2pService), "")
	if err != nil {
		return fmt.Errorf("failed to register container service: %w", err)
	}
	http.Handle("/jrpc", jrpcHandler)

	logrus.Infof("JRPC server listening on port %s", config.Port)
	address := fmt.Sprintf("%s:%s", config.ListenAddress, config.Port)
	if err := http.ListenAndServe(address, nil); err != nil {
		return fmt.Errorf("failed to start jrpc server: %w", err)
	}

	return nil
}
