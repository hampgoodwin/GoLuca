package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

func NewNATSConn(url string) (*nats.Conn, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connecting to nats: %w", err)
	}
	return nc, nil
}
