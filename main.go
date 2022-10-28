package main

import (
	"bufio"
	pb "chat/chat_msg"
	"context"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedChatServer
}

func (s *server) SendMessage(ctx context.Context, in *pb.MessageRequest) (*pb.Message, error) {
	log.Printf("Received: %v", in.GetMsg())
	return &pb.Message{Msg: "GET REKT 42069: " + in.GetMsg()}, nil
}

func runServer(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	lis, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatalf("WE FUCKED UP SERVER: %#v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("WE FUCKED UP: %#v", err)
	}
}

func runClient(wg *sync.WaitGroup) {
	conn, err := grpc.Dial("127.0.0.1:12345", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("WE FUCKED UP CLIENT: %#v", err)
	}
	reader := bufio.NewReader(os.Stdin)
	for true {
		str, _, err := reader.ReadLine()
		if err != nil {
			log.Fatalf("WE FUCKED UP CLIENT: %#v", err)
		}
		defer conn.Close()
		c := pb.NewChatClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.SendMessage(ctx, &pb.MessageRequest{Msg: string(str)})
		if err != nil {
			log.Fatalf("could not send message: %v", err)
		}
		log.Printf("RESPONSE IS: %s", r.GetMsg())
	}
}

func main() {
	var wg sync.WaitGroup
	go runServer(&wg)
	wg.Wait()
	runClient(&wg)
}
