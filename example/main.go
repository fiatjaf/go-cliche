package main

import (
	"log"
	"time"

	"github.com/fiatjaf/go-cliche"
)

func main() {
	c := &cliche.Control{
		DataDir: ".",
		JARPath: "/home/fiatjaf/comp/cliche/target/scala-2.13/cliche-assembly-0.1.0.jar",
	}

	err := c.Start()
	if err != nil {
		log.Fatal(err)
	}

	info, err := c.GetInfo()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(info)

	c.Call("request-hc", map[string]interface{}{
		"pubkey": "02cd1b7bc418fac2dc99f0ba350d60fa6c45fde5ab6017ee14df6425df485fb1dd",
		"host":   "134.209.228.207",
		"port":   80,
	})

	time.Sleep(2 * time.Second)

	inv, err := c.CreateInvoice(cliche.CreateInvoiceParams{
		Msatoshi: 100000, Description: "test invoice"})
	if err != nil {
		log.Fatal(err)
	}
	log.Print(inv)
}
