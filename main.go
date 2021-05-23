//shows how to watch for new devices and list them
package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	driver "gitlab.com/gomidi/rtmididrv"
)

func main() {

	myDriver, err := driver.New()

	if err != nil {
		panic(fmt.Errorf("creating driver failed: %s", err))
	}

	defer myDriver.Close()

	ins, err := myDriver.Ins()
	if err != nil {
		panic(fmt.Errorf("get input streams failed: %s", err))
	}

	for _, input := range ins {
		log.Printf("input device info: %s", input.String())
	}

	in := ins[1]

	err = in.Open()

	defer in.Close()

	if err != nil {
		panic(fmt.Errorf("opening input failed: %s", err))
	}

	myReader := reader.New(reader.Each(func(pos *reader.Position, msg midi.Message) {
		// TODO, This function will trigger
		log.Printf("got message %s\n", msg)
	}),
	)

	err = myReader.ListenTo(in)

	if err != nil {
		panic(fmt.Errorf("reading from input failed: %s", err))
	}

	log.Println("Midi listener added without errors!")

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func(killChannel chan os.Signal, Exit func(int)) {
		for {
			select {
			case <-killChannel:
				Exit(0)
			}
		}
	}(sig, os.Kill)

	waitGroup.Wait()
}
