# go-cliche

This is a small library that will start [cliché](https://github.com/fiatjaf/cliche) and communicate with it via STDIN and STDOUT, allowing you to send commands, receive replies and also receive events.

See the example below or read the full [API docs](https://pkg.go.dev/github.com/fiatjaf/go-cliche) for more (not much more).

## Examples

### Starting it:

```go
c := &cliche.Control{
	DataDir: ".",
    BinaryPath: "/usr/local/bin/cliche",
	// or JARPath: "/home/name/Downloads/cliche.jar", (requires java)
}

err := c.Start()
if err != nil {
	log.Fatal(err)
}
```

### Getting general information:

```go
info, err := c.GetInfo()
if err != nil {
	log.Fatal(err)
}
log.Print(info)
```

### Listening for events

```go
go func () {
	for invpaid := range c.IncomingPayments {
		log.Printf("got payment: %d msat", invpaid.Msatoshi)
    }
}
```

### Creating an invoice

```go
inv, err := c.CreateInvoice(cliche.CreateInvoiceParams{
	Msatoshi: 100000, Description: "test invoice"})
if err != nil {
	log.Fatal(err)
}
log.Print(inv)
```

### Making an arbitrary call

```go
resp, err := c.Call("request-hc", map[string]interface{}{
	"pubkey": "02cd1b7bc418fac2dc99f0ba350d60fa6c45fde5ab6017ee14df6425df485fb1dd",
	"host":   "134.209.228.207",
	"port":   80,
})
# resp will be json.RawMessage aka []byte
```

Other commands are available in the [API docs](https://pkg.go.dev/github.com/fiatjaf/go-cliche).
