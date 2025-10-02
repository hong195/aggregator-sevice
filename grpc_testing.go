package main

import (
	"context"
	"fmt"
	v1 "github.com/hong195/aggregator-sevice/docs/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	client := v1.NewAggregationServiceClient(conn)

	resp, err := client.FindPacketByID(context.Background(), &v1.FindPacketByIDRequest{
		Id: "6bc2ce92-3929-4cf2-a424-cc1e1c143105",
	})
	if err != nil {
		log.Fatalf("error calling FindPacketByID: %v", err)
	}

	fmt.Println("Response:", resp)

	// пример 2025-10-02 23:40:14 +05:00 до сейчас
	loc := time.FixedZone("UTC+5", 5*60*60) // +5
	startTime := time.Date(2025, 10, 2, 23, 40, 14, 0, loc)
	endTime := time.Now().UTC()

	period := &v1.ListPacketsByPeriodRequest{
		Start: timestamppb.New(startTime),
		End:   timestamppb.New(endTime),
	}

	items, err := client.ListPacketsByPeriod(context.Background(), period)

	fmt.Println("Items:", items)
}
