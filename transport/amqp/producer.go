package amqp

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

var (
	ProducerClient *Producer
)

func NewProducer(amqpURI, exchange, exchangeType string) (*Producer, error) {

	producer := &Producer{
		conn:    nil,
		channel: nil,
	}

	var err error

	log.Printf("dialing %q", amqpURI)
	producer.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		log.Printf("Dial: %s", err)
		return nil, err
	}

	log.Printf("got Connection, getting Channel")
	producer.channel, err = producer.conn.Channel()
	if err != nil {
		log.Printf("Channel: %s", err)
		return nil, err
	}

	log.Printf("got Channel, declaring %q Exchange (%q)", exchangeType, exchange)
	if err := producer.channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		log.Printf("Exchange Declare: %s", err)
		return nil, err
	}

	return producer, nil
}

func (producer *Producer) Publish(exchange, exchangeType, routingKey, body string, headers amqp.Table, reliable bool) error {
	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if reliable {
		// log.Printf("enabling publishing confirms.")
		// if err := channel.Confirm(false); err != nil {
		// 	return log.Printf("Channel could not be put into confirm mode: %s", err)
		// }
		//
		// confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		//
		// defer confirmOne(confirms)
	}

	var err error
	log.Printf("Publishing %dB body with routingKey %q (%q)", len(body), routingKey, body)
	if err = producer.channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         headers,
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		log.Printf("Exchange Publish: %s", err)
	}

	return err
}

func (producer *Producer) Shutdown() error {

	if producer == nil{
		return nil
	}

	if err := producer.channel.Cancel("", true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := producer.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP Producer shutdown OK")

	return nil
}

// // One would typically keep a channel of publishings, a sequence number, and a
// // set of unacknowledged sequence numbers and loop until the publishing channel
// // is closed.
// func confirmOne(confirms <-chan amqp.Confirmation) {
// 	log.Printf("waiting for confirmation of one publishing")
//
// 	if confirmed := <-confirms; confirmed.Ack {
// 		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
// 	} else {
// 		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
// 	}
// }
