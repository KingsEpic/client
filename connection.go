package main

import (
	"github.com/KingsEpic/kinglib"
	// "bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	// "io"
	// "log"
	"net"
)

var connection *Connection

type Connection struct {
	NetConn net.Conn
	Decoder *gob.Decoder
	Encoder *gob.Encoder

	Outgoing chan interface{}
	Quit     chan bool

	recv_buffer_decoder *gob.Decoder
	recv_buffer_encoder *gob.Encoder
	recv_buffer         *bytes.Buffer

	send_buffer_decoder *gob.Decoder
	send_buffer_encoder *gob.Encoder
	send_buffer         *bytes.Buffer
}

func (c *Connection) Init() {
	c.send_buffer = new(bytes.Buffer)
	c.recv_buffer = new(bytes.Buffer)
	c.send_buffer_decoder = gob.NewDecoder(c.send_buffer)
	c.send_buffer_encoder = gob.NewEncoder(c.send_buffer)
	c.recv_buffer_decoder = gob.NewDecoder(c.recv_buffer)
	c.recv_buffer_encoder = gob.NewEncoder(c.recv_buffer)

	c.Quit = make(chan bool)
	c.Outgoing = make(chan interface{}, 1000)
}

func (c *Connection) DecodeData(data []byte, i interface{}) {
	c.recv_buffer.Write(data)
	c.recv_buffer_decoder.Decode(i)
	c.recv_buffer.Reset()
}

func (c *Connection) reader() {
	for {
		packet := &kinglib.Packet{}
		err := c.Decoder.Decode(packet)
		if err != nil {
			fmt.Printf("Error decoding: %s\n", err)
			break
		}

		packets <- packet
	}
	c.Quit <- true
}

func (c *Connection) writer() {
Loop:
	for {
		select {
		case obj := <-c.Outgoing:
			c.send_buffer.Reset()

			p := kinglib.Packet{}
			st := kinglib.GetSubType(obj)
			p.SubType = st

			if p.SubType != 0 {
				c.send_buffer_encoder.Encode(obj)
				p.Data = c.send_buffer.Bytes()

				err := c.Encoder.Encode(p)
				if err != nil {
					fmt.Printf("Error encoding: %s\n", err)
					break Loop
				}
				// fmt.Printf("Sending %d\n", atcount)
				// atcount++
			} else {
				fmt.Printf("Unknown object type, could not send message: %v", reflect.TypeOf(obj).String())
			}
		case <-c.Quit:
			break Loop
		}
	}
}

func (c *Connection) SendGob(obj interface{}) {
	c.Outgoing <- obj
}

func (c *Connection) HandleConnection() {
	go c.reader()
	c.writer()

	// TODO: cleaner handling of killed connections.
	// connection = nil
}

func connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", m.Address)
	if err != nil {
		// log.Fatal("Could not resovle address: ", err)
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
		// log.Fatal("Connection err", err)
	}

	connection = &Connection{NetConn: conn, Decoder: gob.NewDecoder(conn), Encoder: gob.NewEncoder(conn)}
	connection.Init()

	go connection.HandleConnection()

	return nil
}
