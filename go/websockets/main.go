// Package main creates the application
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CryptoCompareLTD/sdk/data"
	"github.com/CryptoCompareLTD/sdk/ws"
)

const (
	msgBufferSize = 100
)

func main() {
	dataCh := make(chan ws.Message, msgBufferSize)
	errorCh := make(chan error)

	apiKey := os.Getenv("CCDATA_API_KEY")

	wsManager, err := ws.NewManager(fmt.Sprintf("wss://data-streamer.cryptocompare.com/?api_key=%s", apiKey), dataCh, errorCh)
	if err != nil {
		panic(err)
	}

	dataManager, err := data.NewManager(wsManager, dataCh, errorCh)
	if err != nil {
		panic(err)
	}

	err = dataManager.Start()
	if err != nil {
		panic(err)
	}

	// Prepare for graceful shutdown
	endSignal := make(chan os.Signal, 1)
	signal.Notify(endSignal, syscall.SIGINT, syscall.SIGTERM)

	<-endSignal

	err = wsManager.Stop()
	if err != nil {
		log.Print(err)
	}

	err = dataManager.Stop()
	if err != nil {
		log.Print(err)
	}
}
