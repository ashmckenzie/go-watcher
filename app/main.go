package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	pushbullet "github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/rjeczalik/notify"
	"github.com/urfave/cli"
)

func sendNote(apiKey string, title string, message string) {
	pb := pushbullet.New(apiKey)

	n := requests.NewNote()
	n.Title = title
	n.Body = message

	if _, err := pb.PostPushesNote(n); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}
}

func main() {
	var title string

	app := cli.NewApp()
	app.Name = "watcher"
	app.Usage = "Watch a bunch of directories and notify of change"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "title, t",
			Usage:       "note title",
			EnvVar:      "TITLE",
			Destination: &title,
		},
	}

	app.Action = func(c *cli.Context) error {
		apiKey := os.Getenv("PUSHBULLET_API_KEY")

		if len(apiKey) == 0 {
			return cli.NewExitError("ERROR: API key is empty!", 1)
		}

		spew.Dump(c.Args())

		eventChannel := make(chan notify.EventInfo, 1)

		if err := notify.Watch("/tmp/...", eventChannel, notify.Create); err != nil {
			log.Fatal(err)
		}

		defer notify.Stop(eventChannel)

		for {
			ei := <-eventChannel
			log.Println("Got event:", ei)
		}

		// sendNote(apiKey, title, message)
	}

	app.Run(os.Args)
}
