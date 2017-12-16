package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not upgrade: %v", err)
		return
	}
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Could not read: %v", err)
			break
		}
		log.Printf("Message: %s", message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Printf("Could not write: %v", err)
			break
		}
	}
}

//Run runs the server
func Run(addr string) error {
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return fmt.Errorf("Could not listen and serve: %v", err)
	}
	return nil
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
	<meta charset="utf-8">
	<script>  
	window.addEventListener("load", function(evt) {
		var output = document.getElementById("output");
		var input = document.getElementById("input");
		var ws;
		var print = function(message) {
			var d = document.createElement("div");
			d.innerHTML = message;
			output.appendChild(d);
		};
		document.getElementById("open").onclick = function(evt) {
			if (ws) {
				return false;
			}
			ws = new WebSocket("{{.}}");
			ws.onopen = function(evt) {
				print("OPEN");
			}
			ws.onclose = function(evt) {
				print("CLOSE");
				ws = null;
			}
			ws.onmessage = function(evt) {
				print("RESPONSE: " + evt.data);
			}
			ws.onerror = function(evt) {
				print("ERROR: " + evt.data);
			}
			return false;
		};
		document.getElementById("send").onclick = function(evt) {
			if (!ws) {
				return false;
			}
			print("SEND: " + input.value);
			ws.send(input.value);
			return false;
		};
		document.getElementById("close").onclick = function(evt) {
			if (!ws) {
				return false;
			}
			ws.close();
			return false;
		};
	});
	</script>
	</head>
	<body>
	<table>
	<tr><td valign="top" width="50%">
	<p>Click "Open" to create a connection to the server, 
	"Send" to send a message to the server and "Close" to close the connection. 
	You can change the message and send multiple times.
	<p>
	<form>
	<button id="open">Open</button>
	<button id="close">Close</button>
	<p><input id="input" type="text" value="Hello world!">
	<button id="send">Send</button>
	</form>
	</td><td valign="top" width="50%">
	<div id="output"></div>
	</td></tr></table>
	</body>
	</html>
	`))
