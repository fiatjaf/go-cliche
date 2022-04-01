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

	LogStderr           bool
	LogIrrelevantStdout bool

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

	if c.LogStderr {
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

	stdout, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to open cliche stdout: %w", err)
	}

	go func() {
		reader := bufio.NewReader(stdout)

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				log.Printf("failed to read from stdout: %s", err.Error())
				continue
			}

			// is this a response from a command?
			var response JSONRPCResponse
			if err = json.Unmarshal(line, &response); err == nil {
				if awaiter, ok := c.waiting[response.Id]; ok {
					awaiter <- response
				}
			}

			// is this an event?
			var event Event
			if err = json.Unmarshal(line, &event); err == nil {
				switch event.Event {
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

			// it's not json
			log.Print("stdout: ", strings.TrimSpace(string(line)))
		}
	}()

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

func (c *Control) GetInfo() (result GetInfoResult, err error) {
	resultJson, err := c.Call("get-info", map[string]interface{}{})
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(resultJson, &result)
	return result, err
}
