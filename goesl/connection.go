// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"bufio"
	"bytes"
	toolkit "github.com/charles/toolkit"
	"container/list"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Main connection against ESL - Gotta add more description here
type SocketConnection struct {
	net.Conn
	err chan error
	m   chan *Message
	mtx sync.Mutex
	R   *bufio.Reader
	L   *list.List // list of Future
}

// Dial - Will establish timedout dial against specified address. In this case, it will be freeswitch server
func (c *SocketConnection) Dial(network string, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

// Send - Will send raw message to open net connection
func (c *SocketConnection) Send(cmd string) (*toolkit.Future, error) {

	if strings.Contains(cmd, "\r\n") {
		return nil, fmt.Errorf(EInvalidCommandProvided, cmd)
	}

	// lock mutex
	c.mtx.Lock()
	defer c.mtx.Unlock()

	_, err := io.WriteString(c, cmd)
	if err != nil {
		return nil, err
	}

	_, err = io.WriteString(c, "\r\n\r\n")
	if err != nil {
		return nil, err
	}

	// create a future and enlist and wait
	f := toolkit.NewFuture()
	c.L.PushBack(f)

	return f, nil
}

// SendMany - Will loop against passed commands and return 1st error if error happens
func (c *SocketConnection) SendMany(cmds []string) error {

	for _, cmd := range cmds {
		if _, err := c.Send(cmd); err != nil {
			return err
		}
	}

	return nil
}

// SendEvent - Will loop against passed event headers
func (c *SocketConnection) SendEvent(eventHeaders []string) error {
	if len(eventHeaders) <= 0 {
		return fmt.Errorf(ECouldNotSendEvent, len(eventHeaders))
	}

	// lock mutex to prevent event headers from conflicting
	c.mtx.Lock()
	defer c.mtx.Unlock()

	_, err := io.WriteString(c, "sendevent ")
	if err != nil {
		return err
	}

	for _, eventHeader := range eventHeaders {
		_, err := io.WriteString(c, eventHeader)
		if err != nil {
			return err
		}

		_, err = io.WriteString(c, "\r\n")
		if err != nil {
			return err
		}

	}

	_, err = io.WriteString(c, "\r\n")
	if err != nil {
		return err
	}

	return nil
}

// Execute - Helper fuck to execute commands with its args and sync/async mode
func (c *SocketConnection) Execute(command, args string, sync bool) (m *Message, err error) {
	return c.SendMsg(map[string]string{
		"call-command":     "execute",
		"execute-app-name": command,
		"execute-app-arg":  args,
		"event-lock":       strconv.FormatBool(sync),
	}, "", "")
}

// ExecuteUUID - Helper fuck to execute uuid specific commands with its args and sync/async mode
func (c *SocketConnection) ExecuteUUID(uuid string, command string, args string, sync bool) (m *Message, err error) {
	return c.SendMsg(map[string]string{
		"call-command":     "execute",
		"execute-app-name": command,
		"execute-app-arg":  args,
		"event-lock":       strconv.FormatBool(sync),
	}, uuid, "")
}

// ExecuteUUID - Helper fuck to execute uuid specific commands with its args and sync/async mode
func (c *SocketConnection) ExecuteUUIDAsync(uuid string, command string, args string, sync bool) (f *toolkit.Future, err error) {
	return c.SendMsgAsync(map[string]string{
		"call-command":     "execute",
		"execute-app-name": command,
		"execute-app-arg":  args,
		"event-lock":       strconv.FormatBool(sync),
	}, uuid, "")
}
// SendMsgAsync - Basically this func will send message to the opened connection
// donot wait for response
//
func (c *SocketConnection) SendMsgAsync(msg map[string]string, uuid, data string) (*toolkit.Future, error) {

	b := bytes.NewBufferString("sendmsg")

	if uuid != "" {
		if strings.Contains(uuid, "\r\n") {
			return nil, fmt.Errorf(EInvalidCommandProvided, msg)
		}

		b.WriteString(" " + uuid)
	}

	b.WriteString("\n")

	for k, v := range msg {
		if strings.Contains(k, "\r\n") {
			return nil, fmt.Errorf(EInvalidCommandProvided, msg)
		}

		if v != "" {
			if strings.Contains(v, "\r\n") {
				return nil, fmt.Errorf(EInvalidCommandProvided, msg)
			}

			b.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}

	b.WriteString("\n")

	if msg["content-length"] != "" && data != "" {
		b.WriteString(data)
	}

	// lock mutex
	c.mtx.Lock()
	_, err := b.WriteTo(c)
	if err != nil {
		c.mtx.Unlock()
		return nil, err
	}

	// wait for response
	f := toolkit.NewFuture()
	c.L.PushBack(f)

	c.mtx.Unlock()
	
	return f, nil
}



func (c *SocketConnection) SendMsg(msg map[string]string, uuid, data string) (m *Message, err error) {

	f, err := c.SendMsgAsync(msg, uuid, data)
	if err != nil {
		return nil, err
	}
	
	ret := f.GetResult()
	if ret == nil {
		return nil, fmt.Errorf("future timeout")
	}
	
	if m, ok := ret.(*Message); !ok {
		return nil, fmt.Errorf("Invalid response")
	} else {
		return m, nil
	}
}

// OriginatorAdd - Will return originator address known as net.RemoteAddr()
// This will actually be a freeswitch address
func (c *SocketConnection) OriginatorAddr() net.Addr {
	return c.RemoteAddr()
}

// ReadMessage - Will read message from channels and return them back accordingy.
// If error is received, error will be returned. If not, message will be returned back!
func (c *SocketConnection) ReadMessage() (*Message, error) {
	Debug("Waiting for connection message to be received ...")

	select {
	case err := <-c.err:
		return nil, err
	case msg := <-c.m:
		return msg, nil
	}
}

// Handle - Will handle new messages and close connection when there are no messages left to process
func (c *SocketConnection) Handle() {

	Debug("buffer size %d", ReadBufferSize)
	done := make(chan bool)

	go func() {
		for {
			msg, err := newMessage(c.R, true)

			if err != nil {
				Error(err.Error())

				Debug(err.Error())
				c.err <- err
				done <- true
				break
			}

			// check type of message
			msgType, ok := msg.Headers["Content-Type"]
			if !ok {
				Error("can not find Content-Type")
//				fmt.Println(msg.Headers)
//				fmt.Println(msg)
			} else {
				if msgType == "command/reply" || msgType == "api/response" {
					// if it is reply, dequeue
					c.dequeue(msgType, msg)
					continue
				}
			}

			c.m <- msg
		}
		Error("connection handler exits")
	}()

	<-done

	// Closing the connection now as there's nothing left to do ...
	c.Close()
}

func (c *SocketConnection) dequeue(msgType string, msg *Message) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	v := c.L.Front()
	if v == nil {
		// this is a big error
		Error("Received " + msgType + ", but List is empty, should not happen")
		return
	}

	c.L.Remove(v)

	if v.Value == nil {
		Error("List element is empty, should not happen")
		return
	}

	// we have return message

	if f, ok := v.Value.(*toolkit.Future); ok { // the future
		// it is the future
		f.SetResult(msg)
		// don't put it into channel
	} else {
		Error("List value is not a future, should not happen")
	}

}

// Close - Will close down net connection and return error if error happen
func (c *SocketConnection) Close() error {
	if err := c.Conn.Close(); err != nil {
		return err
	}

	return nil
}
