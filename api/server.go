package api

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/firecracker-microvm/firecracker-go-sdk"

	"google.golang.org/grpc/keepalive"

	log "github.com/sirupsen/logrus"

	node "github.com/PUMATeam/catapult-node/pb"
	"github.com/PUMATeam/catapult-node/service"

	"google.golang.org/grpc"
)

const (
	logFile = "catapult-node.log"
)

func init() {
	// TODO make configurable
	var f *os.File
	if _, err := os.Stat(logFile); err != nil {
		f, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

	}

	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)
}

// Start starts catapult node server
func Start(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer(grpc.KeepaliveParams(
		keepalive.ServerParameters{
			Timeout: 1 * time.Minute,
		}),
	)

	node.RegisterNodeServer(server, &service.NodeService{
		Machines: make(map[string]*firecracker.Machine),
	})
	if err := server.Serve(lis); err != nil {
		log.Error(err)
	}
}
