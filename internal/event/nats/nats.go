package nats

import (
	"fmt"

	"github.com/hampgoodwin/GoLuca/internal/event"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
)

func NewNATSEncodedConn(url string) (*nats.EncodedConn, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connecting to nats: %w", err)
	}
	nec, err := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		return nil, fmt.Errorf("creating new encoded connection for protobuf: %w", err)
	}
	return nec, nil
}

func NewNATSJetStream(nec *nats.EncodedConn) (nats.JetStreamContext, error) {
	jsc, err := nec.Conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("getting jetstream context: %w", err)
	}

	// Configure stream
	accountStreamConfiguration := &nats.StreamConfig{
		Name:        "account",
		Description: "all account subjects",
		Subjects:    []string{event.SubjectAccountCreated},
	}
	_, err = jsc.AddStream(accountStreamConfiguration)
	if err != nil {
		return nil, fmt.Errorf("adding account stream: %w", err)
	}
	transactionStreamConfiguration := &nats.StreamConfig{
		Name:        "transaction",
		Description: "all transaction subjects",
		Subjects:    []string{event.SubjectTransactionCreated},
	}
	_, err = jsc.AddStream(transactionStreamConfiguration)
	if err != nil {
		return nil, fmt.Errorf("adding transaction stream: %w", err)
	}

	return jsc, nil
}
