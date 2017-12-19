package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

//Client is a mid-piece between the server and the connection
type Client struct {
	conn *websocket.Conn
	send chan Message
	rec  chan Message

	srv *Server
}

//ServeClient serves a connected client
func ServeClient(s *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not create websocket: %v", err)
		return
	}

	c := &Client{
		conn: conn,
		send: make(chan Message),
		rec:  make(chan Message),
		srv:  s,
	}
	c.srv.register <- c
	log.Println("Client connected")

	c.Run()
}

//Run runs the client connection
func (c *Client) Run() {
	go func() {
		for {
			mt, cont, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("Read: %v", err)
				break
			}

			c.rec <- Message{typ: mt, cont: cont}
		}
	}()
	if c.srv.echo {
		go func() {
			for {
				mess := <-c.rec
				c.send <- mess
			}
		}()
	} else {
		go func() {
			for {
				mess := <-c.rec
				if mess.typ == websocket.TextMessage {
					c.srv.messages <- mess
				} else if mess.typ == websocket.CloseMessage {
					c.send <- mess
					c.srv.unregister <- c
				} else {
					c.send <- mess
				}
			}
		}()
	}
	go func() {
		for {
			mess := <-c.send
			if err := c.conn.WriteMessage(mess.typ, mess.cont); err != nil {
				log.Printf("Could not write: %v", err)
				break
			}
		}
	}()
}
