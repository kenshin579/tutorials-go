package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

var (
	wg      sync.WaitGroup
	brokers = flag.String("brokers", "", "broker 주소 목록")
	topic   = flag.String("topic", "test", "topic 주제")
)

//do not specify group for consumption
func main() {
	//Create consumer
	consumer, err := sarama.NewConsumer(strings.Split(*brokers, ","), nil)
	if err != nil {
		fmt.Println("Failed to start consumer: %s", err)
		return
	}
	//Set partition
	partitionList, err := consumer.Partitions("nginx_log")
	if err != nil {
		fmt.Println("Failed to get the list of partitions: ", err)
		return
	}
	fmt.Println(partitionList)
	//Cyclic partition
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("nginx_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("Failed to start consumer for partition %d: %s\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				fmt.Println()
			}

		}(pc)
	}
	//time.Sleep(time.Hour)
	wg.Wait()
	consumer.Close()
}
