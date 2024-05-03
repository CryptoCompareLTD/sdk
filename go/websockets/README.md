# websockets

Example app written in Go which connects to the CCData data streamer over websocket (wss://data-streamer.cryptocompare.com) & consumes CADLI tick updates.

## Dependencies

- [nhooyr/websocket](https://github.com/nhooyr/websocket)

## Run

```
export CCDATA_API_KEY=
go run .
```

Alternatively, you can run via the Makefile:

```
make run
```

### Lint

[golangci-lint](https://golangci-lint.run/) can be run from the Makefile:

```
make lint
```