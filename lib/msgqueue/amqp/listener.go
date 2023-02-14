package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/GO_NATIVE/contracts"
	"github.com/GO_NATIVE/lib/msgqueue"
	"github.com/streadway/amqp"
)

type ampqEventListener struct {
	connection *amqp.Connection
	queue      string
}

func (a *ampqEventListener) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		panic(err)
	}

	defer channel.Close()

	_, err = channel.QueueDeclare(a.queue, true, false, false, false, nil)
	return err
}

func NewAMQPEventListener(conn *amqp.Connection, queue string) (msgqueue.EventListener, error) {
	listener := &ampqEventListener{
		connection: conn,
		queue:      queue,
	}

	err := listener.setup()
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func (a *ampqEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil, nil, err
	}

	defer channel.Close()

	for _, eventName := range eventNames {
		err := channel.QueueBind(a.queue, eventName, "events", false, nil)
		if err != nil {
			return nil, nil, err
		}
	}

	msgs, err := channel.Consume(a.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	events := make(chan msgqueue.Event)
	errors := make(chan error)

	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errors <- fmt.Errorf("message did not contain x-event-name header")
				msg.Nack(false, false)
				continue
			}

			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("x-header-name is not a string but a %t", rawEventName)
				msg.Nack(false, false)
				continue
			}

			var event msgqueue.Event

			switch eventName {
			case "event.created":
				event = new(contracts.EventCreatedEvent)
			default:
				errors <- fmt.Errorf("event type %s is unknown", &eventName)
				msg.Nack(false, false)
				continue
			}

			err := json.Unmarshal(msg.Body, event)
			if err != nil {
				errors <- err
				continue
			}
			events <- event
		}
	}()
	return events, errors, nil
}
