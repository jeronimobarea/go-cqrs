package nats

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/nats-io/nats.go"

	"github.com/jeronimobarea/go-cqrs/internal/pubsub"
)

var _ pubsub.EventStorer = &eventStorer{}

const (
	bufferSize = 64
)

type eventStorer struct {
	conn    *nats.Conn
	sub     *nats.Subscription
	channel chan pubsub.Message
}

func NewEventStorer(conn *nats.Conn) *eventStorer {
	return &eventStorer{
		conn: conn,
	}
}

func (e *eventStorer) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
	if e.sub != nil {
		e.sub.Unsubscribe()
	}
	close(e.channel)
}

func (e *eventStorer) Publish(ctx context.Context, msg pubsub.Message) error {
	data, err := e.encodeMessage(msg)
	if err != nil {
		return err
	}
	return e.conn.Publish(msg.Type(), data)
}

func (e eventStorer) encodeMessage(msg pubsub.Message) ([]byte, error) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(msg)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (e *eventStorer) Subscribe(ctx context.Context) (<-chan pubsub.Message, error) {
	e.channel = make(chan pubsub.Message, bufferSize)
	ch := make(chan *nats.Msg, bufferSize)

	var (
		m   pubsub.Message
		err error
	)
	e.sub, err = e.conn.ChanSubscribe(m.Type(), ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				e.decodeMessage(msg.Data, &m)
				e.channel <- m
			}
		}
	}()
	return (<-chan pubsub.Message)(e.channel), nil
}

func (e eventStorer) OnCreate(f func(pubsub.Message)) error {
	var (
		msg pubsub.Message
		err error
	)
	e.sub, err = e.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		e.decodeMessage(m.Data, &msg)
		f(msg)
	})
	return err
}

func (e eventStorer) decodeMessage(data []byte, m any) error {
	b := bytes.NewBuffer(data)
	return gob.NewDecoder(b).Decode(m)
}
