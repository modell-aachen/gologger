package rabbitmq

import "github.com/streadway/amqp"

type DeliveryWrapper interface {
	Ack(multiple bool) (err error)
	GetReplyTo() (replyTo string)
	GetCorrelationId() (correlationId string)
	GetBody() (body []byte)
}
type amqpDeliveryWrapper struct {
	delivery *amqp.Delivery
}

func (wrapper *amqpDeliveryWrapper) Ack(multiple bool) (err error) {
	return wrapper.delivery.Ack(multiple)
}

func (wrapper *amqpDeliveryWrapper) GetReplyTo() (replyTo string) {
	return wrapper.delivery.ReplyTo
}

func (wrapper *amqpDeliveryWrapper) GetCorrelationId() (correlationId string) {
	return wrapper.delivery.CorrelationId
}

func (wrapper *amqpDeliveryWrapper) GetBody() (body []byte) {
	return wrapper.delivery.Body
}

type ChannelWrapper interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}
type channelWrapper struct {
	channel *amqp.Channel
}

func (wrapper channelWrapper) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return wrapper.channel.Publish(exchange, key, mandatory, immediate, msg)
}
