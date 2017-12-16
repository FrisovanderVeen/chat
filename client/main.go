package main

import (
	"os"

	"github.com/FrisovanderVeen/chat/client/cmd"
)

func main() {
	app := cmd.New()
	app.Run(os.Args)
}
