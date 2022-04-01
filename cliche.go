package cliche

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
)

type Control struct {
	JARPath string
	DataDir string

	DontLogStderr bool
	DontLogStdout bool

	stdin   *json.Encoder
	waiting map[string]chan JSONRPCResponse

	PaymentSuccesses chan PaymentSucceededEvent
	PaymentFailures  chan PaymentFailedEvent
	IncomingPayments chan PaymentReceivedEvent
}

func (c *Control) Start() error {
	c.waiting = make(map[string]chan JSONRPCResponse)
	c.PaymentSuccesses = make(chan PaymentSucceededEvent)
	c.PaymentFailures = make(chan PaymentFailedEvent)
	c.IncomingPayments = make(chan PaymentReceivedEvent)

	cmd := exec.Command(
		"java",
		"-Dcliche.datadir="+c.DataDir,
		"-Dcliche.json.compact=true",
		"-jar", c.JARPath,
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to open cliche stdin: %w", err)
	}
	c.stdin = json.NewEncoder(stdin)

	if !c.DontLogStderr {
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("failed to open cliche stderr: %w", err)
		}

		go func() {
			reader := bufio.NewReader(stderr)

			for {
				if line, err := reader.ReadBytes('\n'); err != nil {
					log.Printf("failed to read from stderr: %s", err.Error())
				} else {
					log.Print("stderr: ", strings.TrimSpace(string(line)))
				}
			}
		}()
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to open cliche stdout: %w", err)
	}

	ready := make(chan struct{})

	go func() {
		reader := bufio.NewReader(stdout)

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				log.Printf("failed to read from stdout: %s", err.Error())
				continue
			}

			// is this an event?
			var event Event
			if err = json.Unmarshal(line, &event); err == nil && event.Event != "" {
				switch event.Event {
				case "ready":
					ready <- struct{}{}
				case "payment_succeeded":
					var ps PaymentSucceededEvent
					json.Unmarshal(line, &ps)
					c.PaymentSuccesses <- ps
				case "payment_failed":
					var ps PaymentFailedEvent
					json.Unmarshal(line, &ps)
					c.PaymentFailures <- ps
				case "payment_received":
					var ps PaymentReceivedEvent
					json.Unmarshal(line, &ps)
					c.IncomingPayments <- ps
				}
			}

			// is this a response from a command?
			var response JSONRPCResponse
			if err = json.Unmarshal(line, &response); err == nil {
				if awaiter, ok := c.waiting[response.Id]; ok {
					awaiter <- response
				}
				continue
			}

			// it's not json
			if !c.DontLogStdout {
				log.Print("stdout: ", strings.TrimSpace(string(line)))
			}
		}
	}()

	if err = cmd.Start(); err != nil {
		return fmt.Errorf("failed to start cliche (%s): %w", c.JARPath, err)
	}

	// wait until cliche is ready to receive commands
	<-ready

	return nil
}

func (c *Control) Call(method string, params interface{}) (json.RawMessage, error) {
	id := fmt.Sprintf("id:%d", rand.Int63())
	ch := make(chan JSONRPCResponse, 1)
	c.waiting[id] = ch
	c.stdin.Encode(JSONRPCMessage{"2.0", id, method, params})
	response := <-ch
	if response.Error != nil {
		return nil, fmt.Errorf("'%s' error: '%s'", method, response.Error.Message)
	}

	return response.Result, nil
}
