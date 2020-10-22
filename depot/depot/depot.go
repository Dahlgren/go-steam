// A simple example that uses the modules from the depot package and go-steam to log on
// to the Steam network.
//
// The command expects log on data, optionally with an auth code:
//
//     depot [username] [password]
//     depot [username] [password] [authcode]
package main

import (
	"fmt"
	"os"

	"github.com/Philipp15b/go-steam"
	"github.com/Philipp15b/go-steam/depot"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("depot example\nusage: \n\tdepot [username] [password] [authcode]")
		return
	}
	authcode := ""
	if len(os.Args) > 3 {
		authcode = os.Args[3]
	}

	bot := depot.Default()
	client := bot.Client
	auth := depot.NewAuth(bot, &depot.LogOnDetails{
		Username: os.Args[1],
		Password: os.Args[2],
		AuthCode: authcode,
	}, "sentry.bin")
	debug, err := depot.NewDebug(bot, "debug")
	if err != nil {
		panic(err)
	}
	client.RegisterPacketHandler(debug)
	serverList := depot.NewServerList(bot, "serverlist.json")
	serverList.Connect()

	for event := range client.Events() {
		auth.HandleEvent(event)
		debug.HandleEvent(event)
		serverList.HandleEvent(event)

		switch e := event.(type) {
		case error:
			fmt.Printf("Error: %v", e)
		case *steam.LoggedOnEvent:
			fmt.Print("Logged on")
		}
	}
}
