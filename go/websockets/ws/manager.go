package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/CryptoCompareLTD/sdk/products"
	"nhooyr.io/websocket"
)

// Manager interface defines a WS manager
type Manager interface {
	HandleSubscription(context.Context, products.ProductSubscriptions)
	Stop() error
}

type manager struct {
	url          string
	wg           *sync.WaitGroup
	dataCh       chan Message
	errorCh      chan error
	clients      []Client
	nextClientID int
	stopping     bool
}

// Message describes data received via WS
type Message struct {
	Data   []byte
	T      websocket.MessageType
	Client int
}

// WSError is an error along with the originating WS client ID
type WSError struct {
	Err      error
	ClientID int
}

// Error returns a string containing client ID & error message
func (ws *WSError) Error() string {
	return fmt.Sprintf("Client: %d, Error: %s", ws.ClientID, ws.Err.Error())
}

// NewManager creates a new WS manager
func NewManager(url string, dataCh chan Message, errorCh chan error) (Manager, error) {
	return &manager{
		url:     url,
		clients: make([]Client, 0),
		dataCh:  dataCh,
		errorCh: errorCh,
		wg:      &sync.WaitGroup{},
	}, nil
}

// HandleSubscription is used to create a WS subscription on an available client
func (m *manager) HandleSubscription(ctx context.Context, prodSubs products.ProductSubscriptions) {
	var (
		client Client
		err    error
	)

	if len(m.clients) == 0 {
		client, err = New(ctx, m.url, m.dataCh, m.nextClientID)
		if err != nil {
			m.errorCh <- &WSError{
				Err:      err,
				ClientID: m.nextClientID,
			}
			return
		}

		m.nextClientID++

		m.clients = append(m.clients, client)

		m.wg.Add(1)
		go m.readFromClient(client)
	} else {
		client = m.clients[len(m.clients)-1]
	}

	for productName, instruments := range prodSubs {
		var groups []Group
		switch productName {
		case products.CADLITick:
			groups = DefaultCADLIGroups
		default:
			m.errorCh <- &WSError{
				Err:      errors.New("no groups found for product name"),
				ClientID: client.ID(),
			}
			continue
		}

		var subs []Subscription
		for _, v := range instruments {
			subs = append(subs, Subscription{
				Market:     v.Market,
				Instrument: v.Instrument,
			})
		}

		subscription := SubscriptionMsg{
			Action:       AddSubscription,
			Type:         CADLITick,
			Groups:       groups,
			Subcriptions: subs,
		}

		msg, err := json.Marshal(subscription)
		if err != nil {
			m.errorCh <- &WSError{
				Err:      err,
				ClientID: client.ID(),
			}
			continue
		}

		err = client.Write(ctx, websocket.MessageText, msg)
		if err != nil {
			m.errorCh <- &WSError{
				Err:      err,
				ClientID: client.ID(),
			}
		}
	}
}

func (m *manager) readFromClient(client Client) {
	var err error

	defer m.wg.Done()

readData:
	for !m.stopping {
		var (
			t    websocket.MessageType
			data []byte
		)

		t, data, err = client.Read(context.Background())
		if err != nil {
			break
		}

		msg := Message{
			Data:   data,
			T:      t,
			Client: client.ID(),
		}

		// We try to push the message to the dataCh. If the consumer can't keep
		// up, the channel will fill up and we stop reading.
		select {
		case m.dataCh <- msg:
		default:
			err = errors.New("can't keep up with upstream")
			break readData
		}
	}

	if err != nil {
		m.errorCh <- &WSError{
			Err:      err,
			ClientID: client.ID(),
		}
	}
}

// Stop gracefully ends the WS manager & underlying clients
func (m *manager) Stop() error {
	log.Printf("WS manager stopping.  No of WS clients: %d\n", len(m.clients))

	m.stopping = true

	// Wait for readFromClient goroutines to exit
	m.wg.Wait()

	for _, c := range m.clients {
		// Close WS clients
		err := c.Close()
		if err != nil {
			m.errorCh <- &WSError{
				Err:      err,
				ClientID: c.ID(),
			}
		}
	}

	// Close channels
	close(m.dataCh)
	close(m.errorCh)
	return nil
}
