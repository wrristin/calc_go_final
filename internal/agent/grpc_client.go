package agent

import (
	"calc_service/internal/orchestrator"
	mygrpc "calc_service/pkg/grpc"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcClient mygrpc.TaskServiceClient
	conn       *grpc.ClientConn // Исправлено: grpc.ClientConn из пакета "google.golang.org/grpc"
)

func InitGRPCClient() {
	var err error
	// Исправлено: grpc.Dial и grpc.WithTransportCredentials
	conn, err = grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	grpcClient = mygrpc.NewTaskServiceClient(conn)
}

func CloseGRPCClient() {
	conn.Close()
}

func getTaskGRPC() *orchestrator.Task {
	resp, err := grpcClient.GetTask(context.Background(), &mygrpc.TaskRequest{})
	if err != nil {
		log.Printf("gRPC GetTask error: %v", err)
		return nil
	}
	return &orchestrator.Task{
		ID:            resp.Id,
		Arg1:          resp.Arg1,
		Arg2:          resp.Arg2,
		Operation:     resp.Operation,
		OperationTime: time.Duration(resp.OperationTime) * time.Millisecond,
	}
}

func sendResultGRPC(taskID string, result float64) {
	_, err := grpcClient.SubmitResult(context.Background(), &mygrpc.ResultRequest{
		Id:     taskID,
		Result: result,
	})
	if err != nil {
		log.Printf("gRPC SubmitResult error: %v", err)
	}
}
