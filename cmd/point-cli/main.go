package main

import (
	"context"
	"encoding/json"
	"flag"
	point_service "homework/pkg/point-service/v1"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	methodFlag   = flag.String("method", "{}", "method in API")
	dataFlag     = flag.String("data", "{}", "data in JSON format")
	metadataFlag = flag.String("metadata", "{}", "metadata in JSON format")
)

const (
	grpcServerHost = "127.0.0.1:7001"
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(grpcServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}
	defer conn.Close()

	pointServiceClient := point_service.NewPointServiceClient(conn)

	ctx := context.Background()
	md := metadataParse()
	ctx = metadata.AppendToOutgoingContext(ctx, md...)

	var (
		resp    proto.Message
		respErr error
	)

	switch *methodFlag {
	case "AddOrder":
		req := &point_service.AddOrderRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, respErr = pointServiceClient.AddOrder(ctx, req)
	case "DeleteOrder":
		req := &point_service.DeleteOrderRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, respErr = pointServiceClient.DeleteOrder(ctx, req)
	case "GetOrders":
		req := &point_service.GetOrdersRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, respErr = pointServiceClient.GetOrders(ctx, req)
	case "GiveOrders":
		req := &point_service.GiveOrderRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, respErr = pointServiceClient.GiveOrder(ctx, req)
	case "AcceptReturn":
		req := &point_service.AcceptReturnRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, respErr = pointServiceClient.AcceptReturn(ctx, req)
	case "GetReturns":
		req := &point_service.GetReturnsRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, respErr = pointServiceClient.GetReturns(ctx, req)
	default:
		log.Fatalf("unknown command: %s", *methodFlag)
	}

	data, err := protojson.Marshal(resp)
	if err != nil {
		log.Fatalf("failed to unmarshal data: %v", err)
	}
	log.Printf("resp: %v; err %v\n", string(data), respErr)
}

func metadataParse() []string {
	md := make(map[string]string)
	if err := json.Unmarshal([]byte(*metadataFlag), &md); err != nil {
		log.Fatalf("failed to unmarshal metadata: %v", err)
	}

	kv := []string{
		"x-point-cli", "125234321241241",
	}
	for k, v := range md {
		kv = append(kv, k, v)
	}
	return kv
}
