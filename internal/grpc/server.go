package grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"strings"
	"time"
	internal "vigilant/internal/events"
	pb "vigilant/internal/logger"
)

type GServer struct {
	pb.UnimplementedLogServiceServer
	producer *internal.KProducer
	db       *sql.DB
}

func NewGRPCServer(producer *internal.KProducer) (*GServer, error) {
	dbPath := "logs.db"
	db, err := sql.Open("sqlite3", dbPath)
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS logs (
            id INTEGER PRIMARY KEY,
            message TEXT,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            level TEXT,
            severity INTEGER,
            source TEXT,
            "group" TEXT,
            tags TEXT,
            type TEXT,
            origin TEXT,
            data TEXT
        )
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	return &GServer{producer: producer, db: db}, nil
}

func StartServer(address string, producer *internal.KProducer) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	grpcServer, _ := NewGRPCServer(producer)
	pb.RegisterLogServiceServer(server, grpcServer)
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

func (s *GServer) GetLogs(req *pb.GetLogsRequest, stream pb.LogService_GetLogsServer) error {
	// Apply default values
	logQuery := &pb.GetLogsRequest{
		Limit:     defaultUint32(req.Limit, 100),
		Offset:    req.Offset,
		Source:    req.Source,
		Type:      req.Type,
		Group:     req.Group,
		Tags:      req.Tags,
		Origin:    req.Origin,
		Message:   req.Message,
		Level:     defaultLogLevel(req.Level, pb.LogLevel_INFO),
		Severity:  defaultUint32(req.Severity, 1),
		Timestamp: defaultInt64(req.Timestamp, time.Now().Unix()),
	}

	log.Printf("Processed query: %+v", logQuery)

	// Build the query and arguments dynamically
	query := "SELECT id, message, timestamp, level, severity, source, \"group\", tags, type, origin, data FROM logs WHERE 1=1"
	args := []interface{}{}

	//if logQuery.Level != pb.LogLevel_ALL {
	//	query += " AND level = ?"
	//	args = append(args, logLevelToString(logQuery.Level))
	//}

	//if logQuery.Source != "" {
	//	query += " AND source = ?"
	//	args = append(args, logQuery.Source)
	//}
	//if logQuery.Type != "" {
	//	query += " AND type = ?"
	//	args = append(args, logQuery.Type)
	//}
	//if logQuery.Group != "" {
	//	query += " AND \"group\" = ?"
	//	args = append(args, logQuery.Group)
	//}
	//if logQuery.Tags != "" {
	//	query += " AND tags = ?"
	//	args = append(args, logQuery.Tags)
	//}
	//if logQuery.Origin != "" {
	//	query += " AND origin = ?"
	//	args = append(args, logQuery.Origin)
	//}
	//if logQuery.Message != "" {
	//	query += " AND message LIKE ?"
	//	args = append(args, "%"+logQuery.Message+"%")
	//}
	//if logQuery.Level != pb.LogLevel_NONE {
	//	query += " AND level = ?"
	//	args = append(args, logQuery.Level)
	//}
	//if logQuery.Severity != 0 {
	//	query += " AND severity = ?"
	//	args = append(args, logQuery.Severity)
	//}
	//if logQuery.Timestamp != 0 {
	//	query += " AND timestamp = ?"
	//	args = append(args, logQuery.Timestamp)
	//}

	//query += " ORDER BY timestamp DESC LIMIT ? OFFSET ?"
	args = append(args, logQuery.Limit, logQuery.Offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		return status.Errorf(codes.Internal, "failed to query database: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var log pb.Log
		var data []byte
		var timestamp time.Time
		var levelStr string
		err := rows.Scan(&log.Id, &log.Message, &timestamp, &levelStr, &log.Severity, &log.Source, &log.Group, &log.Tags, &log.Type, &log.Origin, &data)
		if err != nil {
			fmt.Printf("Error scanning row: %v", err)
			return status.Errorf(codes.Internal, "failed to scan row: %v", err)
		}

		// Convert time.Time to Unix timestamp (int64)
		log.Timestamp = timestamp.Unix()

		// Convert string level to LogLevel enum
		log.Level = stringToLogLevel(levelStr)

		// Parse the data JSON if needed
		if len(data) > 0 {
			log.Data = &pb.Log_Data{Fields: make(map[string]string)}
			if err := json.Unmarshal(data, &log.Data.Fields); err != nil {
				fmt.Printf("Error parsing JSON data: %v", err)
				return status.Errorf(codes.Internal, "failed to parse data JSON: %v", err)
			}
		}

		if err := stream.Send(&log); err != nil {
			fmt.Printf("Error sending log: %v", err)
			return status.Errorf(codes.Internal, "failed to send log: %v", err)
		}
		count++
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
		return status.Errorf(codes.Internal, "error during row iteration: %v", err)
	}

	log.Printf("Successfully sent %d logs", count)

	return nil
}

// Helper functions for default values
func defaultUint32(value, defaultValue uint32) uint32 {
	if value == 0 {
		return defaultValue
	}
	return value
}

func defaultInt64(value, defaultValue int64) int64 {
	if value == 0 {
		return defaultValue
	}
	return value
}

func defaultLogLevel(value, defaultValue pb.LogLevel) pb.LogLevel {
	if value == pb.LogLevel_NONE {
		return defaultValue
	}
	return value
}

// Helper function to convert LogLevel enum to string
func logLevelToString(level pb.LogLevel) string {
	switch level {
	case pb.LogLevel_ALL:
		return "ALL"
	case pb.LogLevel_NONE:
		return "NONE"
	case pb.LogLevel_INFO:
		return "INFO"
	case pb.LogLevel_WARN:
		return "WARN"
	case pb.LogLevel_ERROR:
		return "ERROR"
	case pb.LogLevel_DEBUG:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// Helper function to convert string to LogLevel enum
func stringToLogLevel(level string) pb.LogLevel {
	switch strings.ToUpper(level) {
	case "ALL":
		return pb.LogLevel_ALL
	case "NONE":
		return pb.LogLevel_NONE
	case "INFO":
		return pb.LogLevel_INFO
	case "WARN":
		return pb.LogLevel_WARN
	case "ERROR":
		return pb.LogLevel_ERROR
	case "DEBUG":
		return pb.LogLevel_DEBUG
	default:
		return pb.LogLevel_ALL // Default to ALL if unknown
	}
}
