package services

import (
	"encoding/json"
	"gateway/pkg/models/statistic"
	"log"

	"github.com/Shopify/sarama"
)

func SendRequestStatToKafka(stat *statistic.RequestStat, topic string, producer sarama.SyncProducer) {
	statBytes, err := json.Marshal(stat)
	if err != nil {
		log.Printf("Error encoding request stats")
		return
	}

	// Создаем сообщение Kafka
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(statBytes),
	}

	// Отправляем сообщение в Kafka
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Printf("Error sending request stat to Kafka: %v", err)
		return
	}

	log.Printf("Request stat sent to Kafka: %s", string(statBytes))
}
