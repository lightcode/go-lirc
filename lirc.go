package lirc

import (
	"errors"
	"fmt"
	"net/textproto"
	"strconv"
	"strings"
)

const DefaultLircSocket = "/var/run/lirc/lircd"

type Client struct {
	Socket string
}

type Event struct {
	Code              uint64
	RepeatCount       uint64
	ButtonName        string
	RemoteControlName string
}

func NewClient() *Client {
	return &Client{
		Socket: DefaultLircSocket,
	}
}

func (c *Client) Listen(fn func(*Event, error)) error {
	conn, err := textproto.Dial("unix", c.Socket)
	if err != nil {
		return err
	}

	defer conn.Close()

	for {
		line, err := conn.ReadLine()
		if err != nil {
			fn(nil, err)
			continue
		}

		event, err := parseLine(line)
		if err != nil {
			fn(nil, err)
			continue
		}

		fn(event, nil)
	}
}

func parseLine(line string) (*Event, error) {
	s := strings.Split(line, " ")
	if len(s) != 4 {
		return nil, errors.New("Invalid input from lirc")
	}

	code, err := strconv.ParseUint(s[0], 16, 64)
	if err != nil {
		return nil, fmt.Errorf("Fail to parse key code: %s", err)
	}

	rc, err := strconv.ParseUint(s[1], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("Fail to parse repeat count: %s", err)
	}

	e := &Event{
		Code:              code,
		RepeatCount:       rc,
		ButtonName:        s[2],
		RemoteControlName: s[3],
	}

	return e, nil
}
