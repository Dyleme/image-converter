package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Dyleme/image-coverter/pkg/conversion"
	"github.com/Dyleme/image-coverter/pkg/model"
	"github.com/Dyleme/image-coverter/pkg/repository"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type rabbitMQ struct {
	Service *RequestService
	conn    *amqp.Connection
	ch      *amqp.Channel
}

func initRabbit() *rabbitMQ {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

	return &rabbitMQ{conn: conn, ch: ch}
}

func (r *rabbitMQ) send(data *ConvesionData) {
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

func (r *rabbitMQ) receive() {
	q, err := r.ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		logrus.Fatalf("failed to declare a queue: %v", err)
	}

	msgs, err := r.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logrus.Fatalf("failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var data ConvesionData
			err := json.Unmarshal(d.Body, &data)

			if err != nil {
				logrus.Print("Umarshaling error")
			}

			r.convert(&data)

			time.Sleep(1 * time.Second)

			logrus.Println("End receiving")
		}
	}()

	logrus.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (r *rabbitMQ) convert(data *ConvesionData) {
	ctx := context.Background()
	err := r.Service.repo.UpdateRequestStatus(ctx, data.ReqID, repository.StatusProcessing)

	if err != nil {
		logrus.Warn(fmt.Errorf("repo update status in request: %w", err))
	}

	begin := time.Now()

	logrus.WithField("name", data.FileName).Info("start image conversion")

	im, err := decodeImage(bytes.NewBuffer(data.Pic), data.OldType)

	if err != nil {
		logrus.Error(err)
	}

	if data.ImageInfo.Ratio != 1 {
		im = conversion.Convert(im, data.ImageInfo.Ratio)
	}

	pointIndex := strings.LastIndex(data.FileName, ".")
	convFileName := data.FileName[:pointIndex] + "_conv." + data.ImageInfo.Type

	bts, err := encodeImage(im, data.ImageInfo.Type)
	if err != nil {
		logrus.Warn(fmt.Errorf("encode image: %w", err))
	}

	newURL, err := r.Service.uploadFile(ctx, bts, convFileName, data.UserID)
	if err != nil {
		logrus.Warn(fmt.Errorf("upload: %w", err))
	}

	newX, newY := getResolution(im)
	newImageInfo := model.Info{
		Width:  newX,
		Height: newY,
		URL:    newURL,
		Type:   data.ImageInfo.Type,
	}

	newImageID, err := r.Service.repo.AddImage(ctx, data.UserID, newImageInfo)
	if err != nil {
		logrus.Warn(fmt.Errorf("repo add image: %w", err))
	}

	err = r.Service.repo.AddProcessedImageIDToRequest(ctx, data.ReqID, newImageID)
	if err != nil {
		logrus.Warn(fmt.Errorf("repo update image in request: %w", err))
	}

	completionTime := time.Now()

	err = r.Service.repo.AddProcessedTimeToRequest(ctx, data.ReqID, completionTime)
	if err != nil {
		logrus.Warn(fmt.Errorf("repo update time in request: %w", err))
	}

	err = r.Service.repo.UpdateRequestStatus(ctx, data.ReqID, repository.StatusDone)
	if err != nil {
		logrus.Warn(fmt.Errorf("repo update status in request: %w", err))
	}

	logrus.WithFields(logrus.Fields{
		"time for conversion": time.Since(begin),
		"name":                data.FileName,
	}).Info("end image conversion")
}
