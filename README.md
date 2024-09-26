# Reproduction

## Start the server

```sh
$ cd server
$ go run main.go
```

By default the server accepts up to 64 concurrent h2 steams on a single connection.

## Start the proxy

```sh
$ cd proxy
$ cargo run
```

The proxy configures the peer with 100 initial concurrent streams and doesn't limit the number of
streams from a downstream client.

## Start the client

```sh
$ cd client
$ go run main.go
```

The client attempts to open 100 concurrent streams on a single connection, but fails once 64 streams are opened.

