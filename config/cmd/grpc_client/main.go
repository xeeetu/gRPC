package main

import (
	"context"
	"log"
	"time"

	"github.com/xeeetu/gRPC/config/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
	noteId  = 12
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", err)
		}
	}()

	c := note_v1.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &note_v1.GetRequest{Id: noteId})
	if err != nil {
		log.Fatalf("could not get note: %v", err)
	}

	log.Printf("note: %v", r.GetNote())
}
