package pubsub

type Message interface {
	Type() string
}
