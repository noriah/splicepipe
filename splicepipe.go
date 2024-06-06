package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(0)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	data := make(chan string)

	go pipeReader(ctx, data)
	go pipeWriter(ctx, data)

	<-ctx.Done()

	os.Remove(os.Args[1])
	os.Remove(os.Args[2])
}

func pipeReader(ctx context.Context, data chan<- string) {
	p := os.Args[1]

	if err := syscall.Mkfifo(p, 0666); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(p)

	log.Println("created input fifo")

outer_read:
	for {
		f, err := os.OpenFile(p, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("opened input fifo")

		scanner := bufio.NewScanner(f)

	inner_read:
		for scanner.Scan() {
			line := scanner.Text()
			select {
			case <-ctx.Done():
				break inner_read
			case data <- line:
			}
		}

		f.Close()

		log.Println("input fifo died")

		select {
		case <-ctx.Done():
			break outer_read
		default:
			continue
		}
	}
}

func pipeWriter(ctx context.Context, data <-chan string) {
	p := os.Args[2]

	if err := syscall.Mkfifo(p, 0666); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(p)

	log.Println("created output fifo")

outer_write:
	for {
		f, err := os.OpenFile(p, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("opened output fifo")

		count := 0

	inner_write:
		for {
			select {
			case <-ctx.Done():
				break inner_write
			case line := <-data:

				if _, err := fmt.Fprintln(f, count, line); err != nil {
					break inner_write
				}

				count++
			}
		}

		f.Close()

		log.Println("output fifo died")

		select {
		case <-ctx.Done():
			break outer_write
		default:
			continue
		}
	}
}
