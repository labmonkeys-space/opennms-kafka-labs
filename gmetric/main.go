package main

import (
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	metricPb "gmetric/proto/MetricsProtos"
	"google.golang.org/protobuf/proto"
	"os"
	"strings"
)

func main() {
	var bootstrap = flag.String("bootstrap-servers", "127.0.0.1:9092", "Kafka bootstrap servers in the format of host:port")
	var groupId = flag.String("group-id", "gmetric", "Kafka consumer group id")
	var SessionTimeoutMs = flag.Int("session-timeout-ms", 6000, "Kafka consumer session timeout in milliseconds")
	var autoOffsetReset = flag.String("auto-offset-reset", "earliest", "Offset topic reset")
	var topic = flag.String("topic", "flows", "Kafka topic with flows to consume")
	var help = flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	fmt.Println("Try connecting to " + *bootstrap)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  *bootstrap,
		"group.id":           *groupId,
		"session.timeout.ms": *SessionTimeoutMs,
		"auto.offset.reset":  *autoOffsetReset})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n\n", err)
		os.Exit(1)
	}

	err = consumer.SubscribeTopics([]string{*topic}, nil)
	if err != nil {
		fmt.Printf("Error subscribing topic; %s\n\n", err)
		os.Exit(1)
	}
	fmt.Println("Consumer is listening ...")
	defer consumer.Close()
	collectionSet := &metricPb.CollectionSet{}
	for {
		msg, msgErr := consumer.ReadMessage(-1)
		if msgErr == nil {
			proto.Unmarshal(msg.Value, collectionSet)
			if strings.Contains(collectionSet.String(), "rtr-calix-01") {
				fmt.Printf("CollectionSet: %v\n", collectionSet)
			}
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
