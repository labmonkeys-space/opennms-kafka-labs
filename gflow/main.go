package main

import (
	"flag"
	"fmt"
	flowPb "gflow/proto/EnrichedFlowProtos"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
	"os"
)

func main() {
	var reset = "\033[0m"
	var green = "\033[32m"
	var red = "\033[31m"

	var bootstrap = flag.String("bootstrap-servers", "127.0.0.1:9092", "Kafka bootstrap servers in the format of host:port")
	var groupId = flag.String("group-id", "gflow", "Kafka consumer group id")
	var SessionTimeoutMs = flag.Int("session-timeout-ms", 6000, "Kafka consumer session timeout in milliseconds")
	var autoOffsetReset = flag.String("auto-offset-reset", "earliest", "Offset topic reset")
	var topic = flag.String("topic", "flows", "Kafka topic with flows to consume")
	var exporter = flag.String("exporter", "127.0.0.1", "Listen flows for a specific flow exporter")
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
	enFlowDoc := &flowPb.FlowDocument{}
	for {
		msg, msgErr := consumer.ReadMessage(-1)
		if msgErr == nil {
			proto.Unmarshal(msg.Value, enFlowDoc)
			if enFlowDoc.Host == *exporter {
				if enFlowDoc.Direction == 0 {
					// INGRESS
					fmt.Printf("Node ID: %v, Input IfIndex: %v, Flow Version: %s %s←%s Source Address: %s, Source rDNS: %s\n", enFlowDoc.ExporterNode.NodeId, enFlowDoc.InputSnmpIfindex.Value, enFlowDoc.NetflowVersion, red, reset, enFlowDoc.SrcAddress, enFlowDoc.SrcHostname)
					fmt.Printf(" └─ %s (%v) %v bytes\n\n", enFlowDoc.Application, enFlowDoc.DstPort.Value, enFlowDoc.NumBytes.Value)
				} else {
					// EGRESS
					fmt.Printf("Node ID: %v, Output IfIndex: %v, Flow Version: %s %s→%s Source Address: %s, Source rDNS: %s\n", enFlowDoc.ExporterNode.NodeId, enFlowDoc.OutputSnmpIfindex.Value, enFlowDoc.NetflowVersion, green, reset, enFlowDoc.DstAddress, enFlowDoc.DstHostname)
					fmt.Printf(" └─ %s (%v) %v bytes\n\n", enFlowDoc.Application, enFlowDoc.DstPort.Value, enFlowDoc.NumBytes.Value)
				}
			}
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
