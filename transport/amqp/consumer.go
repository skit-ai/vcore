package amqp

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	//done    chan error
}

const (
  ExchangeTopic string = amqp.ExchangeTopic
)


func NewConsumer(amqpURI, exchange, exchangeType, queueName, ctag string, keys []string) (*Consumer, error) {

  c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
//		done:    make(chan error),
	}

	var err error

	log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}
	log.Printf("declared Queue (%q %d messages, %d consumers)", queue.Name, queue.Messages, queue.Consumers)

	for _,key := range keys {
		log.Printf("declared Queue (binding to Exchange (keys %q)", key)

		if err = c.channel.QueueBind(
			queue.Name, // name of the queue
			key,        // bindingKey
			exchange,   // sourceExchange
			false,      // noWait
			nil,        // arguments
		); err != nil {
			return nil, fmt.Errorf("Queue Bind: %s", err)
		}
	}


	return c, nil
}


func (consumer *Consumer) Consume(queue_name string) (<-chan amqp.Delivery, error) {

  log.Printf("starting Consume on queue %q (consumer tag %q)", queue_name,consumer.tag)
  deliveries, err := consumer.channel.Consume(
    queue_name, // name
    consumer.tag,      // consumerTag,
    true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
  )
  if err != nil {
    return deliveries, fmt.Errorf("Queue Consume: %s", err)
  }

  return deliveries, err
}
