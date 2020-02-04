package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	host := flag.String("host", "", "Server transport address.")
	port := flag.Int("port", 8080, "Server transport address port.")
	flag.Parse()

	serverAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		log.Fatalf("connection failed: %s", err.Error())
	}

	fmt.Println("Successfully Connected")

	outputChan := make(chan string)

	go func() {

		for {
			time.Sleep(5 * time.Second)

			bytesWritten, err := conn.Write([]byte("Hello From Peer"))
			if err != nil {
				log.Fatal("Could not write to connection.")
			}

			if bytesWritten == 0 {
				log.Fatal("No bytes written.")
			} else {
				outputChan <- fmt.Sprint("Message Sent Bytes:", bytesWritten)
			}
		}
	}()

	go func() {
		for {
			buf := make([]byte, 1024)
			n, _ := conn.Read(buf)

			if len(buf[:n]) > 0 {
				outputChan <- fmt.Sprintf("Message from Server: %s", string(buf[:n]))
			}
		}
	}()

	for message := range outputChan {
		fmt.Printf("%s\n", message)
	}
}
