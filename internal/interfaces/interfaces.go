package interfaces

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type LogReceiver interface {
	GetDelivery() (metadata LogMetadata, logRow LogRow, err error)
	Close()
}

type LogGeneric map[string]string
type LogMetadata map[string]string
type LogRow struct {
	Time   time.Time
	Source SourceString
	Level  LevelString
	Misc   LogGeneric
	Fields LogFields
}

type LogRpc interface {
	GetRequest() (request LogRpcRequest, startTime time.Time, endTime time.Time, source SourceString, levels []LevelString, err error)
	Close() error
}

type LogRpcRequest interface {
	Reply(fields []LogRow) (err error)
}

type LogReporter interface {
	Send(metadata LogMetadata, row LogRow) error
}

type DeliverySupplier interface {
	GetDelivery() (delivery Delivery)
}

type LogStore interface {
	Store(row LogRow) error
	Read(startTime time.Time, endTime time.Time, source SourceString, levels []LevelString) (rows []LogRow, err error)
	Close()
	CleanUp() error
}

type QueueInstance interface {
	GetRpcReceiver() (logRpc LogRpc, err error)
	GetReceiver(name string) (LogReceiver, error)
	Close() (err error)
}

type QueueChannel interface {
	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name string, key string, exchange string, noWait bool, args amqp.Table) error
	Qos(prefetchCount int, prefetchSize int, global bool) error
	Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table) error
	Close() error
}

type Delivery interface {
	GetReplyTo() string
	GetCorrelationId() string
	Ack(bool) error
	GetBody() []byte
}

type RpcChannel interface {
	Close() error
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

type TimeString string

type LevelString string

type SourceString string

type LogFields []string

func (logFields *LogFields) UnmarshalJSON(bytes []byte) (err error) {
	var data []interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}
	strings := (LogFields)(make([]string, len(data)))
	*logFields = strings
	for idx, value := range data {
		stringValue, ok := value.(string)
		if !ok {
			stringValue = fmt.Sprint(value)
		}
		strings[idx] = stringValue
	}
	return err
}

func (logGeneric *LogGeneric) UnmarshalJSON(bytes []byte) (err error) {
	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}
	strings := (LogGeneric)(make(map[string]string))
	*logGeneric = strings
	for key, value := range data {
		stringValue, ok := value.(string)
		if !ok {
			stringValue = fmt.Sprint(value)
		}
		strings[key] = stringValue
	}
	return err
}
