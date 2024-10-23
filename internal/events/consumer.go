package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
	pb "vigilant/internal/logger"

	"github.com/IBM/sarama"
	_ "github.com/mattn/go-sqlite3"
)

type Consumer struct {
	brokers []string
	topic   string
	db      *sql.DB
}

func NewConsumer(brokers []string, topic string, dbPath string) (*Consumer, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &Consumer{
		brokers: brokers,
		topic:   topic,
		db:      db,
	}, nil
}

func createTable(db *sql.DB) error {
	_, _ = db.Exec(`
		DROP TABLE IF EXISTS logs
`)
	_, err := db.Exec(`
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
	return err
}

func (c *Consumer) Start() error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(c.brokers, config)
	if err != nil {
		return fmt.Errorf("error creating consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Reading messages:")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			c.readMessages(ctx, partitionConsumer)
			cancel()
		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v", err)
		}
	}
}

func (c *Consumer) readMessages(ctx context.Context, partitionConsumer sarama.PartitionConsumer) {
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Message: %s\n", string(msg.Value))
			if err := c.saveToDatabase(string(msg.Value)); err != nil {
				log.Printf("Error saving to database: %v", err)
			}
		case <-ctx.Done():
			return
		default:
			return
		}
	}
}

func (c *Consumer) saveToDatabase(message string) error {
	var logMessage pb.Log
	if err := json.Unmarshal([]byte(message), &logMessage); err != nil {
		return fmt.Errorf("failed to unmarshal message: %v", err)
	}

	// Serialize the Data map to a JSON string
	dataJSON, err := json.Marshal(logMessage.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	_, err = c.db.Exec(`INSERT INTO logs (id, message, timestamp, level, severity, source, "group", tags, type, origin, data) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		logMessage.Id,
		logMessage.Message,
		time.Unix(logMessage.Timestamp, 0),
		logMessage.Level.String(),
		logMessage.Severity,
		logMessage.Source,
		logMessage.Group,
		logMessage.Tags,
		logMessage.Type,
		logMessage.Origin,
		string(dataJSON),
	)
	return err
}

func (c *Consumer) Close() error {
	return c.db.Close()
}
