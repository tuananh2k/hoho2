package kafka

/*
create by: Hoangnd
create at: 2023-01-01
des: Message bus kafka client
*/

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hoho-framework-v2/library"
	"os"
	"sync"

	"github.com/segmentio/kafka-go"
)

const LUFFY_URL = "https://luffy.symper.vn/consumers"

type CallBack func(int64, string, []byte, error)

var ConsumerRunning = map[string]*kafka.Reader{}
var someMapMutex = sync.RWMutex{}

func getAddress() string {
	return os.Getenv("BOOTSTRAP_BROKER")
}
func PublishBulk(topic, event string, resources []interface{}, tenantId int) error {
	// to produce messages
	topic = library.GetPrefixEnvironment() + topic
	if ok := CheckTopicExist(topic); !ok {
		conn, err := kafka.DialLeader(context.Background(), "tcp", getAddress(), topic, 0)
		if err != nil {
			panic(err.Error())
		}
		defer conn.Close()
	}
	// make a writer that produces to topic-A, using the least-bytes distribution
	w := &kafka.Writer{
		Addr:     kafka.TCP(getAddress()),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	for i := 0; i < len(resources); i++ {
		dataPayload := map[string]interface{}{
			"event":     event,
			"time":      library.GetTimeMiliseconds(library.GetCurrentTimeStamp()),
			"data":      resources[i],
			"tenant_id": tenantId,
		}
		dataPayloadJson, _ := json.Marshal(dataPayload)
		err := w.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(event),
			Value: []byte(dataPayloadJson),
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	if err := w.Close(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

func GetListTopic() map[string]struct{} {
	conn, err := kafka.Dial("tcp", getAddress())
	if err != nil {
		// panic(err.Error())
		fmt.Println(err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		// panic(err.Error())
		fmt.Println(err)
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	return m
}

func CheckTopicExist(topic string) bool {
	l := GetListTopic()
	if _, found := l[topic]; found {
		return true
	} else {
		return false
	}
}
func Publish(topic, event string, resource interface{}, tenantId int) error {
	resources := []interface{}{resource}
	return PublishBulk(topic, event, resources, tenantId)
}

func SubscribeMultiTopic(topics []string, consumerId string, callback CallBack, triggerUrl, stopUrl string) error {
	if len(topics) == 0 {
		return errors.New("require topic")
	} else {
		for i := 0; i < len(topics); i++ {
			go GetData(topics[i], consumerId, callback, topics)
		}
	}
	return nil
}
func GetData(topic string, consumerId string, callback CallBack, topics []string) {
	t := library.GetPrefixEnvironment() + topic
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{getAddress()},
		GroupID:   consumerId,
		Topic:     t,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	someMapMutex.Lock()
	ConsumerRunning[consumerId] = r
	someMapMutex.Unlock()
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		callback(m.Offset, topic, m.Value, err)
	}
	if err := r.Close(); err != nil {
		fmt.Println("failed to close reader:", err)
	}
}

func PingLuffy(topics []string, status string, consumerId string) {
	// t, err := json.Marshal(topics)
	// lt := ""
	// if err == nil {
	// 	lt = string(t)
	// }
	// dataPost := map[string]string{
	// 	"serviceId":  consumerId,
	// 	"triggerUrl": "",
	// 	"topics":     lt,
	// 	"processId":  "000",
	// 	"stopUrl":    "",
	// 	"status":     status,
	// }
	// req := new(Request)
	// req.Url = LUFFY_URL
	// req.Body = dataPost
	// req.Method = "POST"
	// res, e := req.Send()
	// fmt.Println(res, e)
}

func StopSubscribe(processId string) {
	someMapMutex.Lock()
	ConsumerRunning[processId].Close()
	delete(ConsumerRunning, processId)
	someMapMutex.Unlock()

}
