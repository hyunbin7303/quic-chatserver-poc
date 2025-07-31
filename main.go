package main

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/quic-go/quic-go" // Assuming you have imported quic-go
)

func main() {
	// Create a TLS configuration (you'll need to provide your own certificate and key)
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{}, // Load your certificate and key here
		NextProtos:   []string{"h3"},      // For HTTP/3
	}

	// Listen on a UDP address
	addr := "localhost:4242"
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Create a QUIC listener
	listener, err := quic.Listen(udpConn, tlsConf, nil) // The third argument is quic.Config
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on %s", addr)

	// Accept incoming QUIC connections
	for {
		conn, err := listener.Accept(nil) // The argument is context.Context
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		log.Printf("Accepted new connection from %s", conn.RemoteAddr())

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn quic.Connection) {
	// Implement your application logic here to handle streams and data
	// For example, accept a stream and read/write data
	stream, err := conn.AcceptStream(nil) // The argument is context.Context
	if err != nil {
		log.Println("Error accepting stream:", err)
		return
	}
	defer stream.Close()

	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil {
		log.Println("Error reading from stream:", err)
		return
	}
	log.Printf("Received: %s", string(buf[:n]))

	_, err = stream.Write([]byte("Hello from QUIC server!"))
	if err != nil {
		log.Println("Error writing to stream:", err)
	}
}
