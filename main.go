package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const (
	ServerPort = "12345"
)

func main() {
	if os.Args[1] == "server" {
		server()
	} else if os.Args[1] == "client" {
		client()
	} else {
		fmt.Printf("Argument not recognized: %s\n", os.Args[1])
	}
}

func server() {
	listener, err := net.Listen("tcp", ":"+ServerPort)
	if err != nil {
		log.Fatalf("Error occurred: %s", err.Error())
	}
	defer listener.Close()

	fmt.Printf("Server is listening for connections on port %s...\n", ServerPort)

	for {
		fmt.Printf("waiting and listening!\n")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func client() {
	guru_ip := os.Getenv("MY_IP")
	addr, err := net.ResolveTCPAddr("tcp", guru_ip+":"+ServerPort)
	if err != nil {
		log.Fatalf("couldn't resolve server IP: %s", addr)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatalf("couldn't connect to server: %s", err)
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hello!"))
	if err != nil {
		log.Fatalf("couldn't write to connection")
	}
	fmt.Printf("%v bytes written\n", n)

	go handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Connected to: %s\n", conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Printf("Client at %s has disconnected.\n", conn.RemoteAddr().String())
			break
		}
		if err != nil {
			log.Printf("Error reading message: %s\n", err.Error())
			break
		}

		message = message[:len(message)-1]
		fmt.Printf("Client: %s\n", message)

		fmt.Print("You: ")
		scanner.Scan()
		reply := scanner.Text() + "\n"

		if _, err := writer.WriteString(reply); err != nil {
			log.Printf("Error sending reply: %s\n", err.Error())
			break
		}
		writer.Flush()
	}
}
