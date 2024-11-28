package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"github.com/xclamation/go-bank-transaction-system/internal/database"
	"github.com/xclamation/go-bank-transaction-system/internal/server"
)

func processTransaction(ctx context.Context, db *database.Queries, req server.TransactionRequest) error {
	// Проверка баланса отправителя
	fromClient, err := db.GetClient(ctx, req.FromClientID)
	if err != nil {
		return fmt.Errorf("invalid from_client_id: %v", err)
	}

	if fromClient.Balance.LessThan(req.Amount) {
		return fmt.Errorf("insufficient balance")
	}

	// Обновление баланса отправителя
	err = db.UpdateClientBalance(ctx, database.UpdateClientBalanceParams{
		ID:      fromClient.ID,
		Balance: fromClient.Balance.Sub(req.Amount),
	})
	if err != nil {
		return fmt.Errorf("invalid fromClient balance update: %v", err)
	}

	// Обновление баланса получателя
	toClient, err := db.GetClient(ctx, req.ToClientID)
	if err != nil {
		return fmt.Errorf("invalid toClient balance update: %v", err)
	}

	err = db.UpdateClientBalance(ctx, database.UpdateClientBalanceParams{
		ID:      toClient.ID,
		Balance: toClient.Balance,
	})
	if err != nil {
		return fmt.Errorf("invalid fromClient balance update: %v", err)
	}

	log.Printf("transaction processed: %v", req)

	return nil
}

func StartWorker(db *database.Queries) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"transactions", //name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var req server.TransactionRequest
			err := json.Unmarshal(d.Body, &req)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			log.Printf("Recieved a transaction request: %v", req)

			err = processTransaction(context.Background(), db, req)
			if err != nil {
				log.Printf("Failed to process transaction: %v", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit precc CTRL+C")

	<-forever
}
