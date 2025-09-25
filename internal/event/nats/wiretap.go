package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	accountv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/account/v1"
	transactionv1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/transaction/v1"
	accountevent "github.com/hampgoodwin/GoLuca/internal/account/event"
	transactionevent "github.com/hampgoodwin/GoLuca/internal/transaction/event"
)

func WireTap(url string) (*nats.Conn, error) {
	nc, _ := nats.Connect(url)

	_, err := nc.Subscribe(">", messageHandler)
	if err != nil {
		return nil, fmt.Errorf("subscribing to all subjects: %w", err)
	}
	return nc, nil
}

func messageHandler(msg *nats.Msg) {
	switch msg.Subject {
	case accountevent.SubjectAccountCreated:
		account := &accountv1.Account{}
		err := proto.Unmarshal(msg.Data, account)
		if err != nil {
			fmt.Printf("error unmarshaling message on subject %q", accountevent.SubjectAccountCreated)
		}
		fmt.Printf("received %q\n%v\n", accountevent.SubjectAccountCreated, account)
	case transactionevent.SubjectTransactionCreated:
		transaction := &transactionv1.Transaction{}
		err := proto.Unmarshal(msg.Data, transaction)
		if err != nil {
			fmt.Printf("error unmarshaling message on subject %q", transactionevent.SubjectTransactionCreated)
		}
		fmt.Printf("received %q\n%v\n", transactionevent.SubjectTransactionCreated, transaction)
	default:
		fmt.Printf("unhandled event received on subject %q\n", msg.Sub.Subject)
	}
}
