package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
)

var (
	brokers = flag.String("brokers", "", "broker 주소 목록")
	topic   = flag.String("topic", "test", "topic 주제")
	message = flag.String("message", "hello world", "보낼 메시지")
)

type producer struct {
	Producer sarama.SyncProducer
}

func NewProducer() *producer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer(strings.Split(*brokers, ","), config)
	if err != nil {
		fmt.Println("producer close, err:", err)
	}
	return &producer{Producer: syncProducer}
}

func (p *producer) sendMsg(message string) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = *topic
	msg.Value = sarama.StringEncoder(message)

	pid, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		fmt.Println("send message failed,", err)
		return
	}

	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}

func main() {
	producer := NewProducer()
	defer producer.Producer.Close()

	producer.sendMsg(*message)
}
