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
	fmt.Println("I am here infra")
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
	fmt.Println("I am here before snd msg")
	go func() {
		kep.producer.Input() <- msg
		fmt.Println("Message was sent")
	}()
	fmt.Println("I am here after snd msg")
	return nil
}
