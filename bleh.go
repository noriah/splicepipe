package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"syscall"
)

func main() {
	p := os.Args[1]

	if err := syscall.Mkfifo(p, 0666); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(p)

	count := 0

	for {
		f, err := os.Open(p)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(count, line)
			count++
		}

		f.Close()
	}
}
