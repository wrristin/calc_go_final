package orchestrator

import (
	"calc_service/pkg/grpc"
	"context"
	"log"
	"net"

	msgrpc "google.golang.org/grpc"
)

type TaskServer struct {
	grpc.UnimplementedTaskServiceServer
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := msgrpc.NewServer()
	grpc.RegisterTaskServiceServer(s, &TaskServer{})
	log.Printf("gRPC server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *TaskServer) GetTask(ctx context.Context, req *grpc.TaskRequest) (*grpc.TaskResponse, error) {
	task, err := DB.GetPendingTask()
	if err != nil || task == nil {
		return nil, err
	}
	return &grpc.TaskResponse{
		Id:            task.ID,
		Arg1:          task.Arg1,
		Arg2:          task.Arg2,
		Operation:     task.Operation,
		OperationTime: int64(task.OperationTime.Milliseconds()),
	}, nil
}

func (s *TaskServer) SubmitResult(ctx context.Context, req *grpc.ResultRequest) (*grpc.ResultResponse, error) {
	err := DB.UpdateTaskResult(req.Id, req.Result)
	return &grpc.ResultResponse{Success: err == nil}, err
}
