package main

import (
	pb "chittychat/stc"
	"log"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
)

type ChittyChatServer struct {
	// type embedded to comply with Google lib
	pb.UnimplementedChittyChatServer

	lamporttime int
	mu          sync.Mutex // protects chatMessages

}

func (s *ChittyChatServer) Chat(stream pb.ChittyChat_ChatServer) error {

}

func broadcast() {
	log.Println("broadcasting")
}

func GetLamportTimestamp() int {
	return 0
}

func main() {
	var port string
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = ":3333"
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	chittyChatSyncServer := &ChittyChatServer{}
	pb.RegisterChittyChatServer(s, chittyChatSyncServer)

	log.Printf("servers listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
