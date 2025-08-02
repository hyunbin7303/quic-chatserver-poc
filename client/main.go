package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	quic "github.com/quic-go/quic-go"
)

type ChatClient struct {
	session  *quic.Conn
	stream   *quic.Stream
	username string
}

func NewChatClient() *ChatClient {
	return &ChatClient{}
}

func (cc *ChatClient) connect(serverAddr string) error {
	session, err := quic.DialAddr(context.Background(), serverAddr, generateTLSConfig(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		return fmt.Errorf("failed to open stream: %v", err)
	}

	cc.session = session
	cc.stream = stream
	return nil
}

func (cc *ChatClient) startReceiving() {
	scanner := bufio.NewScanner(cc.stream)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("\n%s\n", message)
		fmt.Print("> ") // Re-print prompt
	}

	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			log.Printf("Error reading from server: %v", err)
		}
	}
}

func (cc *ChatClient) startSending() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your username: ")

	if scanner.Scan() {
		cc.username = strings.TrimSpace(scanner.Text())
		if cc.username == "" {
			cc.username = "Anonymous"
		}

		// Send username to server
		_, err := cc.stream.Write([]byte(cc.username + "\n"))
		if err != nil {
			log.Printf("Error sending username: %v", err)
			return
		}

		fmt.Print("> ")
	}

	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())

		if message == "" {
			fmt.Print("> ")
			continue
		}

		if message == "/quit" || message == "/exit" {
			fmt.Println("Disconnecting from chat...")
			break
		}

		// Send message to server
		_, err := cc.stream.Write([]byte(message + "\n"))
		if err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}

		fmt.Print("> ")
	}
}

func (cc *ChatClient) run(serverAddr string) error {
	// Connect to server
	if err := cc.connect(serverAddr); err != nil {
		return err
	}
	defer cc.session.CloseWithError(0, "client disconnecting")

	fmt.Printf("Connected to QUIC chat server at %s\n", serverAddr)
	fmt.Println("Type /quit or /exit to disconnect")
	fmt.Println("----------------------------------------")

	// Start receiving messages in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cc.startReceiving()
	}()

	// Start sending messages in main goroutine
	cc.startSending()

	// Wait for receiving goroutine to finish
	wg.Wait()
	return nil
}

func generateTLSConfig() *tls.Config {
	cert, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		NextProtos:         []string{"h3"},
	}
}

func main() {
	serverAddr := "localhost:4242"

	client := NewChatClient()
	if err := client.run(serverAddr); err != nil {
		log.Fatal(err)
	}
}

var serverCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDCTCCAfGgAwIBAgIUV+x7pvcNOQlpUZIsXRjMYd7dzuUwDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI1MDgwMjE4NDcwMFoXDTI2MDgw
MjE4NDcwMFowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAt8oTcOXZXBumFlc4MhSe7Kvoa6I0bRNYU1DdbhTbuvUB
URamm+52WJHQCpra8vqGX+5zVXVwsnQpXjJsuytne7vUx9melowX564FDZjm3FGV
uqdoK/6zevxleS2s5dEUr27sn+6rmEOOEYXXoVWhIAh0CBW/CKqcuL15zRAAzx04
u7e4M4cqlU58NlV1na5GG/92B0VsmSD1AJzCFxAisounGhcWvEhBdfrHChb7p4kO
EWGu+EQOYx29taw5c63/AYIEz9Sadw7iyxAXXXEIUrIRJ2IiBjpYH9IQuLyHSRFK
wn06niWb8QCJR2oDj2G9Dp9SJyOP29h5pDfrxf1YOQIDAQABo1MwUTAdBgNVHQ4E
FgQUx8mUJw2R+u/CttRMoF0d7ntndrkwHwYDVR0jBBgwFoAUx8mUJw2R+u/CttRM
oF0d7ntndrkwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAHnkS
ojfbVOzhcq17KYK6z93BzdyAvVkp6q0RSIeiIYre2rTEce3IaVHc9FWN0sTLgt1J
fKlGo2ZqF7Vm7QY/AMTu9PaHHnGooQA34e106BKM4PYzMhQHnFzL5jc8n84//uL+
zVfUcJDNuvzrJtNAuJcNTppZoaKMmsxa1fUTb3z6eSIBm0g5NVTM1SbHFwuv5g5n
cxnN0HJM2En7R5KekDlvp3QnuEATnMXALItClYaGhAG9Iiqi6vjRmxZGHGRXwcz5
M7NS2omf6mzkUJv6V3gq5VCBrRcyLoOV6lnjVi4jHl2cBnE3B1QU4TNsUbVFW3Fc
gH6JrNH5JWZGFqqy8Q==
-----END CERTIFICATE-----`)
var serverKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC3yhNw5dlcG6YW
VzgyFJ7sq+hrojRtE1hTUN1uFNu69QFRFqab7nZYkdAKmtry+oZf7nNVdXCydCle
Mmy7K2d7u9TH2Z6WjBfnrgUNmObcUZW6p2gr/rN6/GV5Lazl0RSvbuyf7quYQ44R
hdehVaEgCHQIFb8Iqpy4vXnNEADPHTi7t7gzhyqVTnw2VXWdrkYb/3YHRWyZIPUA
nMIXECKyi6caFxa8SEF1+scKFvuniQ4RYa74RA5jHb21rDlzrf8BggTP1Jp3DuLL
EBddcQhSshEnYiIGOlgf0hC4vIdJEUrCfTqeJZvxAIlHagOPYb0On1InI4/b2Hmk
N+vF/Vg5AgMBAAECggEAFngQgg6/pj3vw1nHZ4FTip4GV1OoctCuItrHL27tDBa+
90xnHbd8ushHdbxyhHPHW3wPQpqN0WtX0V8lzDN6vnWmBlwfF3WqHdYNt8NX2gXp
Yb1CBDrqc2CJq/5pswZUL+cbllPl9kKwCt62OtBOrK1hOMnRJedB1hNoDCoRh0Bs
rCw31BI+MkAzg/pFZHMwbTfP80B2UCONE6RycdA7HZhKT7Zpl1Aude5TJT5JKzAa
7h7KvV2ZFhVLwcE/38DK+CuoWPEoUqbM/rpuDrvvZjlBonofm5NMPCFvtpqDBWKV
ycoGsBnhLDUjFCufb+IlOBMDntbhtHu+wJU47d4FbwKBgQDm2QP6+41i1Ear7QRm
RxCblph5DjAM9Ocw/Hfo1tW+xbWWDeCnABxQjtzUkCA7PudCx/bsMFPzH7nn+GCo
0b+sPkCjTjMSIc6FenUm/JfcsgHJXO73ObAD7yTdE/qNKT0Sdaj1DRNw/umbvhKU
NF+7/vC3DbG6pT8fNvaQMicTVwKBgQDL0HjU06rUmUPTaXNabt1JCPW8gkvZQuU1
xeLH74CouNqxH8lQDaeuGXaZZ4T7E+H/Jbm9Epj5/vxcLb4cqjIZ4TYK+pEgsFHT
m9OZ2p7lvrD6LGN/ZQawly/rgFjj5Q/G8iQbv7ZpGVkS+UWDG0nsaeOv2Ipm7pZl
QM5u9rvG7wKBgG4V7Z0B2wHXQ0B3zhJML3JTFbEc//Md0yZ8L16dHN9V/2togMqC
9f3AszS26nf2XmhtXaZywYX+ijRCMS4woFwub7qw7w/liUwEAtwttunrBYkWRWsm
Wnb10zmObnxkvxgPfhwmOTA4kATSVp/Qfhrzz60r3aapaPmkx14qXJIPAoGBALNg
bCfrjosDxOT5BvQNZKYVw6jACB9Tt8VGvxv2FwbnglmnPxc8nVolwPKsYCZVzm4v
drQH/SjxGIvMGmjCBcwvINAyzK23YJzbpTTgaz6KQNo9XOhPMr8SoLMkx5bzD5qp
m8vsQ49mJrYDOwFzb/EpFKG787s5upWsnsKcVpFzAoGAPEkVAIZrRE+lCM03pfSl
cMpw5H2zmPwj82akVTzZEoPK/DLbHUg02dvrZB/Lz1rqSLDFM26FGqFCEsOUA7s6
zsTRgEyD8+EPt8J6KTuioOYofoD5SE7HNkpoqd6vzkp0n/978yemHKHm+wu7duzz
ePggLf6+ubBk/eDfJA4e/Vk=
-----END PRIVATE KEY-----`)
