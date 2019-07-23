# AMQP

## Producer

#### Usage

```go
// AMQP Producer - Initialize and connect
producer, err := amqp.NewProducer(*uri, *exchange, *exchangeType)
if err != nil {
  log.Printf("%s", err)
}
// AMQP Producer - Close connection
defer producer.Shutdown()

producer.Publish(*exchange, *exchangeType, routing_key, body, headers, false)
```

## Consumer

#### Usage

```go
// AMQP Consumer - Initialize and connect
consumer, err := amqp.NewConsumer(*uri, *exchange, *exchangeType, *queue, *consumerTag,<array of bindingKey strings>)
if err != nil {
  log.Printf("%s", err)
}
// AMQP Consumer - Close connection
defer consumer.Shutdown()

// Consume messages
msgs, _ := consumer.Consume(*queue)
```

## Guidelines

##### When RabbitMQ goes down

1. Starting a new process
Don't start the process

2. Existing process

Close the existing connection and shutdown the process.
Try reconnecting to RabbitMQ for configurable x times.

```error
closing: Exception (320) Reason: "CONNECTION_FORCED - broker forced connection closure with reason 'shutdown'"
```
