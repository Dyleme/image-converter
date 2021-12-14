package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"image"

	"github.com/Dyleme/image-coverter/internal/logging"
	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitSender struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
}

func NewRabbitSender(c Config) *RabbitSender {
	connStr := fmt.Sprintf("amqps://%s:%s@%s:%s/", c.User, c.Password, c.Host, c.Port)
	conn, err := amqp.Dial(connStr)

	if err != nil {
		logrus.Fatalf("unable to make connection to rabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("falied in open a channel: %v", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prfectSize
		false, // global
	)
	if err != nil {
		logrus.Fatalf("falied in open a channel: %v", err)
	}

	return &RabbitSender{conn: conn, ch: ch}
}

func (r *RabbitSender) ProcessImage(data *model.ConversionData) {
	r.SendJSON(data)
}

func (r *RabbitSender) SendJSON(data interface{}) {
	q, err := r.ch.QueueDeclare(
		"hello",
		true,  // durable
		false, // delte when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		logrus.Fatalf("unable to make a queue: %v", err)
	}

	jsn, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("rabbitmq: %v", err)
	}

	err = r.ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsn,
		})

	if err != nil {
		logrus.Fatal("uanble to publish message")
	}
}

type Converter interface {
	Convert(ctx context.Context, data *model.ConversionData) image.Image
}

type ResizingProcesser interface {
	ProcessResizedImage(ctx context.Context, im image.Image, data *model.ConversionData)
}

func Receive(ctx context.Context, conv Converter, proc ResizingProcesser, conf Config) {
	logger := logging.FromContext(ctx)
	connStr := fmt.Sprintf("amqps://%s:%s@%s:%s/", conf.User, conf.Password, conf.Host, conf.Port)
	conn, err := amqp.Dial(connStr)

	if err != nil {
		logger.Fatalf("unable to make connection to rabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("falied in open a channel: %v", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prfectSize
		false, // global
	)
	if err != nil {
		logger.Fatalf("falied in open a channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		logger.Fatalf("failed to declare a queue: %v", err)
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
		logger.Fatalf("failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var data model.ConversionData
			err := json.Unmarshal(d.Body, &data)

			if err != nil {
				logger.Print("Umarshaling error")
			}

			im := conv.Convert(ctx, &data)
			proc.ProcessResizedImage(ctx, im, &data)
			logger.Println("End receiving")
		}
	}()

	<-forever
}
