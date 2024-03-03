package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/segmentio/kafka-go"
)

type Registration struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
}

const (
	broker  = "localhost:29092"
	topic   = "user.registration"
	groupId = "kafka-golang-example"
)

func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  groupId,
		Topic:    topic,
		MaxBytes: 10e6,
	})

	go func() {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Printf("message at topic %v/partition %v/offset %v = %s\r\n", m.Topic, m.Partition, m.Offset, string(m.Value))
		}
	}()

	handleRequests()

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader: ", err)
	}
}

func handleRequests() {
	http.HandleFunc("/registration", registrationRequest)
	log.Fatal(http.ListenAndServe(":1000", nil))
}

func registrationRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		registration = Registration{}
		payload, _   = io.ReadAll(r.Body)
	)

	if err := json.Unmarshal(payload, &registration); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := produce(registration); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func produce(registration Registration) error {
	w := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    "user.registration",
		Balancer: &kafka.LeastBytes{},
	}

	payload, _ := json.Marshal(registration)

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Value: payload,
		},
	)

	return err
}
