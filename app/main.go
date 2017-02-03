package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/go-nats-streaming"
	"github.com/rjeczalik/notify"
	"github.com/urfave/cli"
)

func publish(subject string, path string) {
	clusterID := "test-cluster"
	clientID := "go-watcher-new-movies"
	URL := stan.DefaultNatsURL

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	defer sc.Close()

	ch := make(chan bool)
	var glock sync.Mutex
	var guid string

	acb := func(lguid string, err error) {
		glock.Lock()
		// log.Printf("Received ACK for guid %s\n", lguid)

		defer glock.Unlock()
		if err != nil {
			log.Fatalf("Error in server ack for guid %s: %v\n", lguid, err)
		}
		if lguid != guid {
			log.Fatalf("Expected a matching guid in ack callback, got %s vs %s\n", lguid, guid)
		}
		ch <- true
	}

	glock.Lock()
	guid, err = sc.PublishAsync(subject, []byte(path), acb)
	if err != nil {
		log.Fatalf("Error during async publish: %v\n", err)
	}
	glock.Unlock()
	if guid == "" {
		log.Fatal("Expected non-empty guid to be returned.")
	}
	// log.Printf("Published [%s] : '%s' [guid: %s]\n", subject, path, guid)

	select {
	case <-ch:
		break
	case <-time.After(5 * time.Second):
		log.Fatal("timeout")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "watcher"
	app.Usage = "Watch a bunch of directories and notify of change"
	app.Version = "0.1.0"

	app.Action = func(c *cli.Context) error {
		if len(c.Args()) != 1 {
			log.Fatal("Please specify a directory to watch")
		}

		path := strings.TrimSpace(strings.TrimRight(c.Args()[0], "/"))
		fullPath := fmt.Sprintf("%s/...", path)

		eventChannel := make(chan notify.EventInfo, 1)

		if err := notify.Watch(fullPath, eventChannel, notify.Create); err != nil {
			log.Fatal(err)
		}

		defer notify.Stop(eventChannel)

		log.Printf("Watching %s", fullPath)

		for {
			ei := <-eventChannel
			path := ei.Path()
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				log.Println("New movie:", path)
				publish("new-movie", path)
			}
		}
	}

	app.Run(os.Args)
}
