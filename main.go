package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const droneAddress = "192.168.10.1:8889"

func main() {
	// Set up a UDP client to send/recieve commands from the drone
	udpAddr, err := net.ResolveUDPAddr("udp", droneAddress)

	if err != nil {
		fmt.Printf("Error resolving drone UDP address: %s", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		fmt.Printf("Error connecting to drone: %s", err)
		os.Exit(1)
	}

	defer conn.Close()
	
	fmt.Println("Entering SDK mode...")
	conn.Write([]byte("command"))

	fmt.Println("Ready for commands!")

	go printMessagesFromDrone(conn)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		newText := strings.Split(text, "\n")
		
		conn.Write([]byte(newText[0]))

		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("Stopping")
			return
		}
	}
}

func printMessagesFromDrone(conn *net.UDPConn) {
	for {
		message := make([]byte, 256)

		conn.ReadFromUDP(message)

		fmt.Printf("🚁 Message from drone: %s\n>> ", message)
	}
}
