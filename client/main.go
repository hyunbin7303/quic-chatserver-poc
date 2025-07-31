package client

// import (
// 	"crypto/tls"
// 	"io"
// 	"log"
// 	"net/http"

// 	"github.com/quic-go/quic-go/http3" // For HTTP/3 client
// )

// func main() {
// 	// For a generic QUIC client, you would use quic.DialAddr
// 	// For HTTP/3, use http3.RoundTripper

// 	// Create a custom HTTP client with the http3.RoundTripper
// 	client := &http.Client{
// 		Transport: &http3.RoundTripper{
// 			TLSClientConfig: &tls.Config{
// 				InsecureSkipVerify: true, // For testing, skip certificate verification
// 			},
// 		},
// 	}

// 	// Make an HTTP/3 request
// 	resp, err := client.Get("https://localhost:4242/") // Replace with your server address
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("Response: %s", string(body))
// }
