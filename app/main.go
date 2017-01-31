package main

import (
  // "os"
  "log"

  // "github.com/urfave/cli"
  "github.com/rjeczalik/notify"
)

func main() {
  // cli.NewApp().Run(os.Args)

  // Make the channel buffered to ensure no event is dropped. Notify will drop
  // an event if the receiver is not able to keep up the sending pace.
  c := make(chan notify.EventInfo, 1)

  // Set up a watchpoint listening for events within a directory tree rooted
  // at current working directory. Dispatch remove events to c.
  if err := notify.Watch("/tmp/...", c, notify.Create); err != nil {
      log.Fatal(err)
  }

  defer notify.Stop(c)

  // Block until an event is received.
  for {
    ei := <-c
    log.Println("Got event:", ei)
  }
}
