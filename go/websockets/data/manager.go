// Package data handles the coordination with the WS manager, the data & any errors received back
package data

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/CryptoCompareLTD/sdk/products"
	"github.com/CryptoCompareLTD/sdk/ws"
)

// Manager interface defines a data manager
type Manager interface {
	Start() error
	Stop() error
}

type dataManager struct {
	wsManager          ws.Manager
	wg                 *sync.WaitGroup
	dataCh             chan ws.Message
	errorCh            chan error
	cadliTickCh        chan ws.Message
	bufferedCadliTicks []ws.CADLITickMsg
}

type metadata struct {
	CCType string `json:"TYPE"`
}

// NewManager creates a new data manager
func NewManager(wsManager ws.Manager, dataCh chan ws.Message, errorCh chan error) (Manager, error) {
	return &dataManager{
		wsManager:          wsManager,
		dataCh:             dataCh,
		errorCh:            errorCh,
		cadliTickCh:        make(chan ws.Message, 1),
		bufferedCadliTicks: make([]ws.CADLITickMsg, 0, 10),
		wg:                 &sync.WaitGroup{},
	}, nil
}

// Start begins the data management handlers & requests the initial WS subscriptions
func (d *dataManager) Start() error {
	log.Printf("Data manager starting...\n")

	d.wg.Add(3)

	go d.handleMessages()
	go d.handleErrors()
	go d.handleCADLITick()

	expectedProducts := products.Load()

	log.Printf("Expected products: %+v\n", expectedProducts)

	d.wsManager.HandleSubscription(context.Background(), expectedProducts)

	return nil
}

// Stop gracefully ends the manager
func (d *dataManager) Stop() error {
	log.Printf("Data manager stopping...\n")
	d.wg.Wait()
	return nil
}

func (d *dataManager) handleMessages() {
	defer func() {
		d.wg.Done()
		close(d.cadliTickCh)
	}()

	for msg := range d.dataCh {
		// Detect message type
		// Alternatively, could read the first few bytes
		var m metadata
		err := json.Unmarshal(msg.Data, &m)
		if err != nil {
			log.Print(err)
			continue
		}

		switch m.CCType {
		case string(ws.CADLITick):
			log.Printf("CADLITick msg received, client ID: %d, msg: %+v\n", msg.Client, m)
			d.cadliTickCh <- msg
		default:
			// Other message types should be handled appropriately in order to manage the state of each subscription
			log.Printf("Other msg received, client ID: %d, msg: %+v\n", msg.Client, string(msg.Data))
		}
	}
}

func (d *dataManager) handleErrors() {
	defer d.wg.Done()

	for err := range d.errorCh {
		// Detect errors in order to manage the state of each subscription
		log.Printf("Error received: %s\n", err.Error())
	}
}

func (d *dataManager) handleCADLITick() {
	//nolint:gomnd // 5 seconds used as an example ticker duration
	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
		d.wg.Done()
	}()

outer:
	for {
		select {
		case msg, ok := <-d.cadliTickCh:
			if !ok {
				break outer
			}

			var m ws.CADLITickMsg
			err := json.Unmarshal(msg.Data, &m)
			if err != nil {
				log.Print(err)
				continue
			}

			log.Printf("CADLITick msg handled, client ID: %d, tick: %+v\n", msg.Client, m)

			// Example: buffer updates
			d.bufferedCadliTicks = append(d.bufferedCadliTicks, m)
		case <-ticker.C:
			if len(d.bufferedCadliTicks) == 0 {
				log.Print("No buffered updates to process...")
				continue
			}

			log.Printf("Processing buffered updates: %d\n", len(d.bufferedCadliTicks))

			// Example: write buffered updates & clear the buffer
			d.bufferedCadliTicks = d.bufferedCadliTicks[len(d.bufferedCadliTicks):]
		}
	}

	log.Printf("Processing any remaining buffered updates: %d\n", len(d.bufferedCadliTicks))
}
