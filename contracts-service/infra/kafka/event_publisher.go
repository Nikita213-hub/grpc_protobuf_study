package kafkaContracts

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/domain/events"
)

type KafkaEventPublisher struct {
	topic    string
	producer sarama.AsyncProducer
}

func NewKafkaEventPublisher(brokers []string, topic string) *KafkaEventPublisher {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(brokers, config)
	go func() {
		for e := range producer.Errors() {
			fmt.Println(e)
		}
	}()
	if err != nil {
		panic(err) //TODO: handle it
	}
	return &KafkaEventPublisher{
		topic:    topic,
		producer: producer,
	}
}

func (kep *KafkaEventPublisher) PublishContractEvent(event *events.ContractEvent) error {
	msgString, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: kep.topic,
		Value: sarama.ByteEncoder(msgString),
	}
	go func() {
		kep.producer.Input() <- msg
	}()
	return nil
}
