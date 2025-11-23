package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	bufCh := make(chan string, 2)

	go channel(bufCh, "buffered")

	for {
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled, exiting for loop")
			return

		case msg, ok := <-bufCh:
			if ok {
				log.Printf("Received data from buffered channel :%s", msg)
				continue
			}
			cancel()
			return
		default:
			log.Printf("No data is receiving from channel")
			time.Sleep(time.Second)
		}
	}
}

func channel(ch chan string, chName string) {
	for i := range 10 {
		log.Printf("Generating value #%d for %s channel", i+1, chName)
		ch <- fmt.Sprintf("value #%d", i+1)
		log.Printf("Length of %s channel: %d", chName, len(ch))
	}
	close(ch)
}
