package logconnection

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/modell-aachen/gologger/internal/interfaces"
)

type Closable interface {
	Close() error
}

type LogConnection struct {
	channel    Closable
	deliveries interfaces.DeliverySupplier
}

func CreateLogConnection(channel Closable, deliveries interfaces.DeliverySupplier) *LogConnection {
	return &LogConnection{
		channel,
		deliveries,
	}
}

func (c *LogConnection) Close() {
	c.channel.Close()
}

type rabbitMessage struct {
	Metadata interfaces.LogMetadata
	Log_data interfaces.LogRow
}

func (c *LogConnection) GetDelivery() (logMetadata interfaces.LogMetadata, logRow interfaces.LogRow, err error) {
	delivery := c.deliveries.GetDelivery()
	err = delivery.Ack(true)
	if err != nil {
		return logMetadata, logRow, errors.Wrap(err, "Could not ACK delivery")
	}

	message := rabbitMessage{}
	err = json.Unmarshal(delivery.GetBody(), &message)
	if err != nil {
		return logMetadata, logRow, errors.Wrapf(err, "Could not unmarshal json: %s", delivery.GetBody())
	}
	logMetadata = message.Metadata
	logRow = message.Log_data

	return logMetadata, logRow, nil
}
