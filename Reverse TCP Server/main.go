package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// Constants
const maxBufferSize = 1024
const address = "127.0.0.1:3333"

// TCPServer function that handles incoming client connections
func TCPServer(ready chan bool) {
	// Start listening on the specified address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	// Signal that the server is ready
	ready <- true
	fmt.Println("Server is ready and listening on", address)

	// Accept incoming client connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}
}

// handleConnection processes a single client connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the incoming message
	message := make([]byte, maxBufferSize)
	n, err := conn.Read(message)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading from connection:", err)
		return
	}

	// Reverse the message
	reversedMessage := reverseString(string(message[:n]))

	// Send the reversed message back to the client
	_, err = conn.Write([]byte(reversedMessage))
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		return
	}
}

// reverseString reverses the input string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Client-side function to send messages and receive reversed messages
func tcpClient(messages []string) ([]string, error) {
	// Resolve TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return []string{}, err
	}

	var reversed []string

	// For each message, create a connection to the server, send the message, and read the response
	for _, msg := range messages {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return []string{}, err
		}
		defer conn.Close()

		// Send the message
		_, err = conn.Write([]byte(msg))
		if err != nil {
			return []string{}, err
		}

		// Read the server's response
		reply := make([]byte, maxBufferSize)
		n, err := conn.Read(reply)
		if err != nil {
			return []string{}, err
		}

		// Append the reversed message to the result
		reversed = append(reversed, string(reply[:n]))
	}

	return reversed, nil
}

// Read a line from stdin
func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}
	return strings.TrimRight(string(str), "\r\n")
}

// Check if there's any error
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Set up the output file
	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)
	defer stdout.Close()

	// Set up input and output readers/writers
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)
	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	// Read the number of messages
	messagesCount, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)

	// Read the messages
	messages := make([]string, messagesCount)
	for i := int64(0); i < messagesCount; i++ {
		messages[i] = readLine(reader)
	}

	// Send the messages to the server and get the reversed messages
	reversedMessages, err := tcpClient(messages)
	checkError(err)

	// Write the reversed messages to stdout
	for _, message := range reversedMessages {
		_, err := writer.WriteString(message + "\n")
		checkError(err)
	}

	// Flush the writer
	writer.Flush()
}
