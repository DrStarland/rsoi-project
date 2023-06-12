package initialization

import (
	"gateway/pkg/utils"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type KafkaSettings struct {
	Topic    string
	Producer sarama.SyncProducer
}

func InitKafka(logger *zap.SugaredLogger) *KafkaSettings {
	kafkaBrokers := utils.Config.Kafka.Endpoints
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	config := sarama.NewConfig()
	config.Net.TLS.Enable = false
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(kafkaBrokers, config)
	if err != nil {
		logger.Errorln("Error creating Kafka producer: %v", err)
	}

	for err != nil {
		time.Sleep(5 * time.Second)
		producer, err = sarama.NewSyncProducer(kafkaBrokers, config)
		logger.Errorln("Error creating Kafka producer: %v", err)
	}

	return &KafkaSettings{
		Topic:    utils.Config.Kafka.Topics[0],
		Producer: producer,
	}
}
