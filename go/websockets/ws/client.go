// Package ws defines a way to interact with external data over websocket
package ws

import (
	"context"
	"log"

	"nhooyr.io/websocket"
)

// Client interface defines a websocket client
type Client interface {
	Read(context.Context) (websocket.MessageType, []byte, error)
	Write(context.Context, websocket.MessageType, []byte) error
	Close() error
	ID() int
}

type wsClient struct {
	conn   *websocket.Conn
	dataCh chan Message
	id     int
}

// New creates a new wsClient and returns it as Client
func New(ctx context.Context, url string, dataCh chan Message, id int) (Client, error) {
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	return &wsClient{
		conn:   conn,
		dataCh: dataCh,
		id:     id,
	}, nil
}

func (c *wsClient) Read(
	ctx context.Context,
) (websocket.MessageType, []byte, error) {
	return c.conn.Read(ctx)
}

func (c *wsClient) Write(
	ctx context.Context,
	t websocket.MessageType,
	data []byte,
) error {
	return c.conn.Write(ctx, t, data)
}

func (c *wsClient) Close() error {
	log.Printf("Closing client: %d\n", c.id)
	return c.conn.Close(websocket.StatusNormalClosure, "close requested")
}

func (c *wsClient) ID() int {
	return c.id
}
