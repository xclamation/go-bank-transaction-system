package server

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/streadway/amqp"
)

type TransactionRequest struct {
	FromClientID uuid.UUID       `json:"from_client_id"`
	ToClientID   uuid.UUID       `json:"to_client_id"`
	Amount       decimal.Decimal `json:"amount"`
}

func SendTransaction(req TransactionRequest) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"transactions", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	// ИСПОЛЬЗОВАТЬ EASYJSON ВМЕСТО ОБЫЧНОГО JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil
	}

	err = ch.Publish(
		"",     // exhange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	log.Printf("Sent transaction request: %v", req)

	return nil

}
