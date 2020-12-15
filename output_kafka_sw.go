package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/buger/goreplay/byteutils"
	"github.com/buger/goreplay/proto"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
)

// KafkaOutput is used for sending payloads to kafka in JSON format.
type KafkaSwOutput struct {
	config   *OutputKafkaConfig
	producer sarama.AsyncProducer
}

// KafkaOutputFrequency in milliseconds
const KafkaSwOutputFrequency = 500

// NewKafkaOutput creates instance of kafka producer client  with TLS config
func NewKafkaSwOutput(address string, config *OutputKafkaConfig, tlsConfig *KafkaTLSConfig) PluginWriter {
	c := NewKafkaConfig(tlsConfig)

	//SW
	println("SW-HTTP", config.Host, config.SwSource)

	var producer sarama.AsyncProducer

	if mock, ok := config.producer.(*mocks.AsyncProducer); ok && mock != nil {
		producer = config.producer
	} else {
		c.Producer.RequiredAcks = sarama.WaitForLocal
		c.Producer.Compression = sarama.CompressionSnappy
		c.Producer.Flush.Frequency = KafkaSwOutputFrequency * time.Millisecond

		brokerList := strings.Split(config.Host, ",")

		var err error
		producer, err = sarama.NewAsyncProducer(brokerList, c)
		if err != nil {
			log.Fatalln("Failed to start Sarama(Kafka) producer:", err)
		}
	}

	o := &KafkaSwOutput{
		config:   config,
		producer: producer,
	}

	// Start infinite loop for tracking errors for kafka producer.
	go o.ErrorHandler()

	return o
}

// ErrorHandler should receive errors
func (o *KafkaSwOutput) ErrorHandler() {
	for err := range o.producer.Errors() {
		Debug(1, "Failed to write access log entry:", err)
	}
}

// PluginWrite writes a message to this plugin
func (o *KafkaSwOutput) PluginWrite(msg *Message) (n int, err error) {
	var message sarama.StringEncoder

	mimeHeader := proto.ParseHeaders(msg.Data)
	meta := payloadMeta(msg.Meta)
	req := msg.Data

	//SW: Data
	kafkaMessage := buildSwMessage(o.config.SwSource, mimeHeader, meta, req)

	jsonMessage, _ := json.Marshal(&kafkaMessage)
	message = sarama.StringEncoder(byteutils.SliceToString(jsonMessage))

	o.producer.Input() <- &sarama.ProducerMessage{
		//SW TOPIC
		Topic: "collector_api_access",
		Value: message,
	}

	return len(message), nil
}
