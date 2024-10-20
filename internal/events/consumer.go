package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

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
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS logs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            message TEXT,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
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

	ticker := time.NewTicker(1 * time.Minute)
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
	_, err := c.db.Exec("INSERT INTO logs (message) VALUES (?)", message)
	return err
}

func (c *Consumer) Close() error {
	return c.db.Close()
}
