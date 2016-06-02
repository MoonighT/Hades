package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MoonighT/Hades/Service"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: mybot slack-bot-token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	ws, id := service.SlackConnect(os.Args[1])
	fmt.Println("mybot ready, ^C exits")

	for {
		// read each incoming message
		m, err := service.GetMessage(ws)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(m)
		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			parts := strings.Fields(m.Text)
			if len(parts) == 3 && parts[1] == "stock" {
				// looks good, get the quote and reply with the result
				go func(m service.Message) {
					m.Text = service.GetQuote(parts[2])
					service.PostMessage(ws, m)
				}(m)
				// NOTE: the Message object is copied, this is intentional
			} else {
				// huh?
				m.Text = fmt.Sprintf("sorry, that does not compute\n")
				service.PostMessage(ws, m)
			}
		}
	}
}
