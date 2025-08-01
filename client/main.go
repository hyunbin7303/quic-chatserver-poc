package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"

	quic "github.com/quic-go/quic-go"
)

func main() {

	session, err := quic.DialAddr(context.Background(), "localhost:4242", &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{}, // Load your certificate and key here
		NextProtos:         []string{"h3"},      // Add this line
		ServerName:         "kevin",
	}, nil)

	if err != nil {
		log.Fatal(err)
	}
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	_, err = stream.Write([]byte("Hello from QUIC client!"))
	if err != nil {
		log.Fatal(err)
	}
	reply := make([]byte, 1024)
	n, err := stream.Read(reply)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	fmt.Printf("Client received: %s\n", string(reply[:n]))
}
