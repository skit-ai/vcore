package amqp

import (
	"log"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewProducer(amqpURI, exchange, exchangeType string) (*Producer, error) {

	producer := &Producer{
			conn:    nil,
			channel: nil,
		}

	var err error

	log.Printf("dialing %q", amqpURI)
	producer.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		log.Fatalf("Dial: %s", err)
	}
	//defer connection.Close()

	log.Printf("got Connection, getting Channel")
	producer.channel, err = producer.conn.Channel()
	if err != nil {
		log.Fatalf("Channel: %s", err)
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
		log.Fatalf("Exchange Declare: %s", err)
	}

	return producer,err
}

func (producer *Producer) Publish(exchange, exchangeType,routingKey, body string, headers amqp.Table,reliable bool) error {
	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if reliable {
		// log.Printf("enabling publishing confirms.")
		// if err := channel.Confirm(false); err != nil {
		// 	return log.Fatalf("Channel could not be put into confirm mode: %s", err)
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
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		log.Fatalf("Exchange Publish: %s", err)
	}

	return err
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
