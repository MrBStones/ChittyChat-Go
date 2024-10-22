package main

import (
	pb "chittychat/stc"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
)

type ChittyChatServer struct {
	// type embedded to comply with Google lib
	pb.UnimplementedChittyChatServer

	lamporttime  int
	mu           sync.Mutex
	participants map[string]pb.ChittyChat_BroadcastServer
}

func (s *ChittyChatServer) Broadcast(usr *pb.User, stream pb.ChittyChat_BroadcastServer) error {
	s.mu.Lock()
	s.participants[usr.User] = stream
	s.lamporttime = max(int(usr.Logicaltimestamp), int(s.lamporttime)) + 1
	joinMessage := &pb.ChatMessage{
		Message:          fmt.Sprintf("Participant %s joined Chitty-Chat at Lamport time %d", usr.User, s.lamporttime),
		User:             usr.User,
		Logicaltimestamp: int32(s.lamporttime),
	}
	log.Printf("Participant %s joined Chitty-Chat at Lamport time %d", usr.User, s.lamporttime)
	for _, participantStream := range s.participants {
		if err := participantStream.Send(joinMessage); err != nil {
			log.Printf("error sending join message to %s", usr.User)
		}
	}
	s.mu.Unlock()

	<-stream.Context().Done()

	s.mu.Lock()
	delete(s.participants, usr.User)
	s.lamporttime++
	leaveMessage := &pb.ChatMessage{
		Message:          fmt.Sprintf("Participant %s left Chitty-Chat at Lamport time %d", usr.User, s.lamporttime),
		User:             usr.User,
		Logicaltimestamp: int32(s.lamporttime),
	}
	log.Printf("Participant %s left Chitty-Chat at Lamport time %d", usr.User, s.lamporttime)
	for _, participantStream := range s.participants {
		if err := participantStream.Send(leaveMessage); err != nil {
			log.Printf("error sending leave message to %s", usr.User)
		}
	}
	s.mu.Unlock()

	return nil
}

func (s *ChittyChatServer) Publish(ctx context.Context, message *pb.ChatMessage) (*pb.PublishResponse, error) {
	s.mu.Lock()
	s.lamporttime = max(int(message.Logicaltimestamp), int(s.lamporttime)) + 1
	if len(message.Message) > 128 {
		s.mu.Unlock()
		return &pb.PublishResponse{Success: false}, fmt.Errorf("message exceeds 128 characters")
	}
	for _, participantStream := range s.participants {
		if err := participantStream.Send(message); err != nil {
			log.Printf("error sending message to %s", message.User)
		}
	}
	s.mu.Unlock()
	s.lamporttime++
	log.Println("from:", message.User, "; message:", message.Message, "; timestamp:", message.Logicaltimestamp)
	return &pb.PublishResponse{Success: true, Logicaltimestamp: int32(s.lamporttime)}, nil
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
	chittyChatSyncServer := &ChittyChatServer{
		lamporttime:  0,
		participants: make(map[string]pb.ChittyChat_BroadcastServer),
	}
	pb.RegisterChittyChatServer(s, chittyChatSyncServer)

	log.Printf("servers listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
