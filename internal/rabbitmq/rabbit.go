package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"image"

	"github.com/Dyleme/image-coverter/internal/logging"
	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/streadway/amqp"
)

// RabbitSender is a struct, which is used to send data
// to the  image converter using RabbitMQ.
type RabbitSender struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// Config to connect to the message broker.
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
}

// Name of the queue, which is used to communicate with the RabbitMQ.
var queueName = "convert"

// ConversionData is struct, which contains images and all needed
// information to convert images.
type ConversionData struct {
	Ctx       context.Context
	ImageInfo model.ConversionInfo `json:"imageInfo"`
	UserID    int                  `json:"userID"`
	ReqID     int                  `json:"reqID"`
	OldType   string               `json:"oldType"`
	Pic       []byte               `json:"pic"`
	FileName  string               `json:"fileName"`
}

// NewRabbitSender returns *RabbitSender, which is ready to send messages.
// NewRabbitSender at first initialize connection with RabbitMQ server,
// than it initialize channel with broker.
func NewRabbitSender(c *Config) (*RabbitSender, error) {
	connStr := fmt.Sprintf("amqps://%s:%s@%s:%s/", c.User, c.Password, c.Host, c.Port)
	conn, err := amqp.Dial(connStr)

	if err != nil {
		// logrus.Fatalf("unable to make connection to rabbitMQ: %v", err)
		return nil, fmt.Errorf("unable to make connection to rabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		// logrus.Fatalf("falied in open a channel: %v", err)
		return nil, fmt.Errorf("falied in open a channel: %w", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prfectSize
		false, // global
	)

	if err != nil {
		// logrus.Fatalf("falied in open a channel: %v", err)
		return nil, fmt.Errorf("falied in open a channel: %w", err)
	}

	return &RabbitSender{conn: conn, ch: ch}, nil
}

// This function is used to send images and data to convert it, to the message broker.
func (r *RabbitSender) ProcessImage(ctx context.Context, data *ConversionData) {
	r.SendJSON(ctx, data)
}

// This function send data to the message broker.
// At first this function initialize queue to communicate with message broker,
// Then it marshals data in json and send this json to the queue.
// If any error occurs, this function log it to the logger, getted from context.
func (r *RabbitSender) SendJSON(ctx context.Context, data interface{}) {
	logger := logging.FromContext(ctx)
	q, err := r.ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delte when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		logger.Fatalf("unable to make a queue: %v", err)
	}

	jsn, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("rabbitmq: %v", err)
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
		logger.Fatal("uanble to publish message")
	}
}

// Converter is an interface which provide functions to convert images.
type Converter interface {
	Convert(ctx context.Context, data *ConversionData) image.Image
	ProcessResizedImage(ctx context.Context, im image.Image, data *ConversionData)
}

// Receive is method which is used to get messages from RabbitMQ and then convert images.
// At first this function initialize connection, channel and queue to with RabbitMQ.
// Then it in infinite loop get messages from queue, convert image and process it.
func Receive(ctx context.Context, conv Converter, conf *Config) {
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
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
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

	logger.Info("start conversion server")

	go func() {
		for d := range msgs {
			logger.Println("get conversion reqeust")

			var data ConversionData
			err := json.Unmarshal(d.Body, &data)

			if err != nil {
				logger.Print("Umarshaling error")
			}

			im := conv.Convert(ctx, &data)
			conv.ProcessResizedImage(ctx, im, &data)
			logger.Println("conversion request is handled")
		}
	}()

	<-forever
}
