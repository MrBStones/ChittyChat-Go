package main

import (
	"bufio"
	pb "chittychat/stc"
	"context"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	name             string
	lamporttimestamp int
	mu               sync.Mutex
)

func Broadcast(cl pb.ChittyChatClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mu.Lock()
	lamporttimestamp++
	mu.Unlock()

	stream, err := cl.Broadcast(ctx, &pb.User{User: name, Logicaltimestamp: int32(lamporttimestamp)})
	if err != nil {
		return err
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive a message: %v", err)
		}

		mu.Lock()
		lamporttimestamp = max(int(lamporttimestamp), int(msg.Logicaltimestamp)) + 1
		mu.Unlock()

		log.Println("User:", msg.User, "; Message:", msg.Message, "; Lamport:", msg.Logicaltimestamp)
	}
}

func Publish(cl pb.ChittyChatClient) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		log.Printf("Enter message: ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		if len(text) > 128 {
			log.Fatalf("Message exceeds 128 characters")
			continue
		}
		mu.Lock()
		lamporttimestamp++
		mu.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := cl.Publish(ctx, &pb.ChatMessage{Message: text, Logicaltimestamp: int32(lamporttimestamp), User: name})
		if err != nil {
			log.Fatalf("Failed to publish message: %v", err)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s name", os.Args[0])
	}
	name = os.Args[1]
	conn, err := grpc.NewClient("localhost:3333", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewChittyChatClient(conn)

	go Broadcast(client)
	go Publish(client)

	// Wait for interrupt signal to gracefully shutdown the server.
	c := make(chan os.Signal, 1)
	<-c
}
