package main

import (
	"flag"
	"fmt"
	eventPb "gevent/proto/OpennmsModelProtos"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
	"os"
	"time"
)

func main() {
	var bootstrap = flag.String("bootstrap-servers", "127.0.0.1:9092", "Kafka bootstrap servers in the format of host:port")
	var groupId = flag.String("group-id", "gevent", "Kafka consumer group id")
	var SessionTimeoutMs = flag.Int("session-timeout-ms", 6000, "Kafka consumer session timeout in milliseconds")
	var autoOffsetReset = flag.String("auto-offset-reset", "earliest", "Offset topic reset")
	var topic = flag.String("topic", "events", "Kafka topic with flows to consume")
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
		fmt.Println("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	err = consumer.SubscribeTopics([]string{*topic}, nil)
	if err != nil {
		fmt.Println("Error subscribing topic; %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Consumer is listening ...")
	defer consumer.Close()
	event := &eventPb.Event{}
	for {
		msg, msgErr := consumer.ReadMessage(-1)
		if msgErr == nil {
			proto.Unmarshal(msg.Value, event)
			latency := time.Now().UnixMilli() - int64(event.GetCreateTime())
			fmt.Printf("Latency: %vms, Event ID: %v, Node Label: %v (%v), UEI: %v\n", latency, event.GetId(), event.GetNodeCriteria().GetNodeLabel(), event.GetNodeCriteria().GetId(), event.GetUei())
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
