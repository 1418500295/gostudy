package utils

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func init() {
	fmt.Println("-- 自动结算正式开始 ---")
}

// JetStreamContext 返回JetStreamContext
func JetStreamContext(nc *nats.Conn) nats.JetStreamContext {
	var err error
	nc, err = nats.Connect("nats://123.51.206.118:4222")
	if err != nil {
		log.Fatalf("连接NATS失败: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("创建JetStream上下文失败: %v", err)
	}
	return js
}

// Create creates the named stream.
func Create(js nats.JetStreamContext, name string) *nats.StreamInfo {
	fmt.Printf("Creating stream: %q\n", name)
	strInfo, err := js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: []string{"Settlement.*", "Settlement"},
		MaxAge:   0, // 0 means keep forever
		Storage:  nats.FileStorage,
	})
	if err != nil {
		log.Panicf("could not create stream: %v", err)
	}

	prettyPrint(strInfo)
	return strInfo
}

func prettyPrint(x interface{}) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		log.Fatalf("could not prettyPrint: %v", err)
	}
	fmt.Println(string(b))
}

// Delete deletes the named stream.
func Delete(js nats.JetStreamContext, name string) {
	fmt.Printf("Deleting stream: %q\n", name)
	if err := js.DeleteStream(name); err != nil {
		log.Printf("error deleting stream: %v", err)
	}
}

func AddConsumer(js nats.JetStreamContext, strName, consName, consFilter string) {
	info, err := js.AddConsumer(strName, &nats.ConsumerConfig{
		Durable:       consName,
		AckPolicy:     nats.AckExplicitPolicy,
		MaxAckPending: 1, // default value is 20,000
		FilterSubject: consFilter,
	})
	if err != nil {
		log.Panicf("could not add consumer: %v", err)
	}
	prettyPrint(info)
}
