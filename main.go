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
	listener, err := quic.ListenAddr("localhost:4242", generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server listening on localhost:4242 ...")
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		go func(sess *quic.Conn) {
			stream, err := sess.AcceptStream(context.Background())
			if err != nil {
				log.Println(err)
				return
			}
			buf := make([]byte, 1024)
			n, err := stream.Read(buf)
			if err != nil && err != io.EOF {
				log.Println(err)
			}
			fmt.Printf("Server received: %s\n", string(buf[:n]))

			stream.Write([]byte("Hello from server over QUIC!"))
		}(sess)
	}
}

func generateTLSConfig() *tls.Config {
	cert, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert}, // Load your certificate and key here
		NextProtos:         []string{"h3"},          // Add this line
	}
}

var serverCert = []byte(`-----BEGIN CERTIFICATE-----
MIIB...
-----END CERTIFICATE-----`)
var serverKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIE...
-----END RSA PRIVATE KEY-----`)
