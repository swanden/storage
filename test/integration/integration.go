package main

import (
	"context"
	grpcStorage "github.com/swanden/storage/api/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	allStart := time.Now()

	creds := insecure.NewCredentials()

	conn, err := grpc.DialContext(ctx, "localhost:8001", grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		log.Panicln("error while grpc connect:", err)
	}
	defer conn.Close()

	log.Println("start")

	client := grpcStorage.NewStorageClient(conn)

	start := time.Now()

	key := "TestKey"
	value := "TestValue"
	ttl := int64(5)

	if _, err := client.Set(ctx, &grpcStorage.SetRequest{
		Key:   key,
		Value: value,
		Ttl:   ttl,
	}); err != nil {
		log.Panicln("unable to set key-value pair:", err)
	}

	log.Printf("key-value pair seccessfully set. spent: %v\n", time.Since(start))

	start = time.Now()

	resp, err := client.Get(ctx, &grpcStorage.GetRequest{
		Key: key,
	})
	if err != nil {
		log.Panicln("unable to get value:", err)
	}
	if resp.GetValue() != value {
		log.Panicf("values are not equal, want: %q, got: %q\n", value, resp.GetValue())
	}

	log.Printf("value seccessfully got. spent: %v\n", time.Since(start))

	start = time.Now()

	if _, err := client.Delete(ctx, &grpcStorage.DeleteRequest{
		Key: key,
	}); err != nil {
		log.Panicln("unable to delete value:", err)
	}

	resp, err = client.Get(ctx, &grpcStorage.GetRequest{
		Key: key,
	})

	st, ok := status.FromError(err)
	if !ok {
		log.Panicln("error was not a status error:", err)
	}
	if st.Code() != codes.NotFound {
		log.Panicln("value is still in storage:", err)
	}

	if err != nil && st.Code() != codes.NotFound {
		log.Panicln("unable to get value:", err)
	}

	log.Printf("value seccessfully deleted. spent: %v\n", time.Since(start))

	log.Println("time spent:", time.Since(allStart))
	log.Println("bye")
}
