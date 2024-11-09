package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	ch := make(chan string)

	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			log.Printf("Generating value #%d in goroutine", i+1)
			ch <- time.Now().Format(time.RFC1123)
		}
		close(ch)
	}()

	go interruptSignal(cancel)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled, exiting for loop")
			return
		case msg, ok := <-ch:
			if ok {
				log.Printf("Received value from channel: %s", msg)
				continue
			}
			log.Printf("Channel closed, stopping for loop")
			cancel()
			return
		}
	}
}

func interruptSignal(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Received interrupt signal")
		cancel()
	}()
}
