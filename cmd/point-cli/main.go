package main

import (
	"context"
	"encoding/json"
	point_service "homework/pkg/point-service/v1"
	"log"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	methodFlag   = pflag.String("method", "{}", "method in API")
	dataFlag     = pflag.String("data", "{}", "data in JSON format")
	metadataFlag = pflag.String("metadata", "{}", "metadata in JSON format")
)

const (
	grpcServerHost = "127.0.0.1:7001"
)

func main() {
	pflag.Parse()
	conf := newConfig(cliFlags)

	prod, err := NewKafkaProducer(conf.kafka.Brokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer prod.Close()

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
		sendOrderEvent(prod, conf.producer.topic, "AddOrder", req.OrderId)

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
		sendOrderEvent(prod, conf.producer.topic, "GiveOrder", req.OrderIds)

	case "AcceptReturn":
		req := &point_service.AcceptReturnRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, respErr = pointServiceClient.AcceptReturn(ctx, req)
		sendOrderEvent(prod, conf.producer.topic, "AcceptReturn", req.OrderId)

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

func sendOrderEvent(prod KafkaProducerInterface, topic, eventType string, id any) {
	event := map[string]interface{}{
		"eventType": eventType,
		"id":        id,
		"timestamp": time.Now().UnixNano(), // временная метка в миллисекундах
	}
	eventData, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}
	err = prod.SendMessage(topic, eventType, string(eventData))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
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
