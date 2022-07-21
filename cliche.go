package cliche

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

type Control struct {
	JARPath    string
	BinaryPath string
	DataDir    string

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

	var cmd *exec.Cmd
	var usingPath string
	if c.BinaryPath != "" {
		usingPath = c.BinaryPath
		cmd = exec.Command(
			c.BinaryPath,
			"-Dcliche.datadir="+c.DataDir,
			"-Dcliche.json.compact=true",
		)
	} else if c.JARPath != "" {
		usingPath = c.JARPath
		cmd = exec.Command(
			"java",
			"-Dcliche.datadir="+c.DataDir,
			"-Dcliche.json.compact=true",
			"-jar", c.JARPath,
		)
	} else {
		return fmt.Errorf("must specify BinaryPath or JARPath, but both are empty")
	}

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
					log.Printf("[go-cliche] failed to read from stderr: %s", err.Error())
					if err == io.EOF || err == io.ErrUnexpectedEOF {
						return
					}
					time.Sleep(30 * time.Second)
				} else {
					log.Print("[go-cliche] stderr: ", strings.TrimSpace(string(line)))
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
				log.Printf("[go-cliche] failed to read from stdout: %s", err.Error())
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					return
				}
				time.Sleep(30 * time.Second)
				continue
			}

			// is this an event?
			var event JSONRPCNotification
			if err = json.Unmarshal(line, &event); err == nil && event.Method != "" {
				switch event.Method {
				case "ready":
					ready <- struct{}{}
				case "payment_succeeded":
					var ps PaymentSucceededEvent
					json.Unmarshal(event.Params, &ps)
					c.PaymentSuccesses <- ps
				case "payment_failed":
					var ps PaymentFailedEvent
					json.Unmarshal(event.Params, &ps)
					c.PaymentFailures <- ps
				case "payment_received":
					var ps PaymentReceivedEvent
					json.Unmarshal(event.Params, &ps)
					c.IncomingPayments <- ps
				}
				continue
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
				log.Print("[go-cliche] stdout: ", strings.TrimSpace(string(line)))
			}
		}
	}()

	if err = cmd.Start(); err != nil {
		return fmt.Errorf("failed to start cliche (%s): %w", usingPath, err)
	}

	// wait until cliche is ready to receive commands
	<-ready

	return nil
}

func (c *Control) Call(method string, params interface{}) (json.RawMessage, error) {
	id := fmt.Sprintf("id:%d", rand.Int63())
	ch := make(chan JSONRPCResponse, 1)
	c.waiting[id] = ch
	err := c.stdin.Encode(JSONRPCRequest{id, method, params})
	if err != nil {
		return nil,
			fmt.Errorf("error writing json to cliche stdin ('%s'): %w", method, err)
	}
	response := <-ch
	if response.Error != nil {
		return nil, fmt.Errorf("'%s' error: '%s'", method, response.Error.Message)
	}

	return response.Result, nil
}
