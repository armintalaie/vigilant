package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	internal "vigilant/internal/events"
	pb "vigilant/internal/logger"
)

type GServer struct {
	pb.UnimplementedLogServiceServer
	producer *internal.KProducer
}

func NewGRPCServer(producer *internal.KProducer) *GServer {
	return &GServer{producer: producer}
}

func StartServer(address string, producer *internal.KProducer) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterLogServiceServer(server, NewGRPCServer(producer))
	reflection.Register(server)

	return server.Serve(lis)
}

func (s *GServer) SendLog(ctx context.Context, req *pb.Log) (*pb.LogResponse, error) {
	log := pb.Log{
		Id:        req.Id,
		Message:   req.Message,
		Level:     pb.LogLevel(req.Level),
		Severity:  uint32(int(req.Severity)),
		Timestamp: req.Timestamp,
		Source:    req.Source,
		Data:      req.Data,
		Group:     req.Group,
		Tags:      req.Tags,
		Type:      req.Type,
		Origin:    req.Origin,
	}

	fmt.Printf("Received log: %v\n", log.Message)
	err := s.producer.AddLog(log)
	if err != nil {
		fmt.Printf("Failed to send log: %v\n", err)
		return nil, err
	}
	return &pb.LogResponse{
		Success: true,
	}, nil
}
