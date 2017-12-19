package server

import (
	"log"
	"net/http"
)

//Server is a server for the chat
type Server struct {
	//all the clients
	clients map[*Client]bool

	//the message pool
	messages chan Message

	register   chan *Client
	unregister chan *Client

	//the user decided values
	addr string
	echo bool
}

//NewServer creates a new server
func NewServer(addr string, echo bool) *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		messages:   make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		addr:       addr,
		echo:       echo,
	}
}

//Run runs the server
func (s *Server) Run() {
	http.HandleFunc("/", serveDefault)
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		ServeClient(s, w, r)
	})

	go func() {
		for {
			select {
			case c := <-s.register:
				s.clients[c] = true
			case c := <-s.unregister:
				delete(s.clients, c)
				close(c.send)
				close(c.rec)
			case mess := <-s.messages:
				for c := range s.clients {
					c.send <- mess
				}
			}
		}
	}()

	log.Printf("Running on: %s", s.addr)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		log.Printf("Could not listen and serve: %v", err)
		return
	}
}

func serveDefault(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.Error(w, "Not found", 404)
	return
}
