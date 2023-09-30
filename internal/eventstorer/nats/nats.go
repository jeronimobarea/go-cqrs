package nats

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/nats-io/nats.go"

	"github.com/jeronimobarea/go-cqrs/internal/eventstorer"
)

var _ eventstorer.EventStorer = &eventStorer{}

const (
	bufferSize = 64
)

type eventStorer struct {
	conn    *nats.Conn
	sub     *nats.Subscription
	channel chan eventstorer.Message
}

func NewEventStorer(conn *nats.Conn) *eventStorer {
	return &eventStorer{
		conn: conn,
	}
}

func (e *eventStorer) Publish(ctx context.Context, msg eventstorer.Message) error {
	data, err := e.encodeMessage(msg)
	if err != nil {
		return err
	}
	return e.conn.Publish(msg.Type(), data)
}

func (e eventStorer) encodeMessage(msg eventstorer.Message) ([]byte, error) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(msg)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (e *eventStorer) Subscribe(ctx context.Context) (<-chan eventstorer.Message, error) {
	e.channel = make(chan eventstorer.Message, bufferSize)
	ch := make(chan *nats.Msg, bufferSize)

	var (
		m   eventstorer.Message
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
	return (<-chan eventstorer.Message)(e.channel), nil
}

func (e eventStorer) OnCreate(f func(eventstorer.Message)) error {
	var (
		msg eventstorer.Message
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
