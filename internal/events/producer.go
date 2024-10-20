package internal

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"vigilant/internal/logger"
)

type KProducer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewKProducer(topic string, producer sarama.SyncProducer) *KProducer {
	return &KProducer{topic: topic, producer: producer}
}

func (p *KProducer) Close() {
	p.producer.Close()
}

func (p *KProducer) AddLog(log logger.Log) error {
	logJson, err := json.Marshal(log)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(logJson),
	}
	_, _, err = p.producer.SendMessage(msg)
	return err
}
