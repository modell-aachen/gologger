package rabbitmq

import (
	"fmt"
	"github.com/pkg/errors"
	"os"

	"github.com/modell-aachen/gologger/interfaces"

	"github.com/streadway/amqp"
)

const (
	routingKey = "qwiki.foswiki_perl"
	exchange   = "ma_logs"
	rpc_queue  = "rpc_queue"
)

type rabbitInstance struct {
	connection *amqp.Connection
}

func CreateInstance() (instance interfaces.QueueInstance, err error) {
	user := os.Getenv("RABBIT_USER")
	if user == "" {
		return nil, errors.New("RABBIT_USER not set")
	}
	password := os.Getenv("RABBIT_PASSWORD")
	host := os.Getenv("RABBIT_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("RABBIT_PORT")
	if port == "" {
		port = "5672"
	}
	dialString := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	connection, err := amqp.Dial(dialString)
	if err != nil {
		errors.WithStack(err)
		return nil, errors.Wrap(err, "Failed to connect to RabbitMQ")
	}

	instance = rabbitInstance{
		connection,
	}

	return instance, nil
}

func (instance rabbitInstance) Close() (err error) {
	if instance.connection != nil {
		err := instance.connection.Close()
		if err != nil {
			errors.WithStack(err)
			return errors.Wrap(err, "Could not close instance")
		}
		instance.connection = nil
	}

	return nil
}

func (instance rabbitInstance) GetRpcReceiver() (logRpc interfaces.LogRpc, err error) {
	if instance.connection == nil {
		return logRpc, errors.New("No open connection")
	}

	channel, err := instance.connection.Channel()
	if err != nil {
		return logRpc, errors.Wrap(err, "Failed to open a channel")
	}

	_, err = channel.QueueDeclare(
		rpc_queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return logRpc, errors.Wrap(err, "Could not declare queue")
	}
	err = channel.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		return logRpc, errors.Wrap(err, "Could not set QoS")
	}

	msgs, err := channel.Consume(
		rpc_queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		errors.WithStack(err)
		return logRpc, errors.Wrap(err, "Failed to register a consumer")
	}

	logRpc = rpcConnection{
		channel,
		msgs,
	}

	return logRpc, nil
}

func (instance rabbitInstance) GetReceiver(name string) (interfaces.LogReceiver, error) {
	if instance.connection == nil {
		return nil, errors.New("No open connection")
	}

	ch, err := instance.connection.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open a channel")
	}

	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to declare an exchange")
	}

	var queue amqp.Queue
	if name != "" {
		queue, err = ch.QueueDeclare(
			name,  // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
	} else {
		queue, err = ch.QueueDeclare(
			"",    // name
			false, // durable
			true,  // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to declare a queue with name \"%s\"", name)
	}

	err = ch.QueueBind(
		queue.Name,
		"#",
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to register a consumer")
	}

	c := logConnection{
		ch,
		msgs,
	}

	return c, nil
}
