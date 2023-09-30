package eventstorer

type Message interface {
	Type() string
}
