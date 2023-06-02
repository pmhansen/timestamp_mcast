package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	// Configure multicast address, port and nodeid
	multicastAddr := "225.0.0.1"         // IPv4 multicast address
	multicastPort := "3000"              // Multicast port
	networkInterfaceIP := "172.19.10.12" // IP address of the network interface
	nodeId := 1                          // ID of this node

	if len(os.Args[1:]) != 4 {
		fmt.Println("\nRunning without arguments, using defaults:")
		fmt.Printf("Run with arguments %s <mcast group> <mcast port> <nic ip> <nodeid>\n", os.Args[0])
		fmt.Printf("defaults: %s %s %s %s %d\n\n", os.Args[0], multicastAddr, multicastPort, networkInterfaceIP, nodeId)
	} else {
		multicastAddr = os.Args[1]
		multicastPort = os.Args[2]
		networkInterfaceIP = os.Args[3]
		nodeId, _ = strconv.Atoi(os.Args[4])
	}

	// Resolve the multicast address
	addr, err := net.ResolveUDPAddr("udp", multicastAddr+":"+multicastPort)
	if err != nil {
		log.Fatal(err)
	}

	// Create a UDP address for the network interface
	var localAddr *net.UDPAddr
	localAddr, err = net.ResolveUDPAddr("udp", networkInterfaceIP+":0")
	if err != nil {
		log.Fatal(err)
	}

	// Create a UDP connection
	var conn *net.UDPConn
	conn, err = net.DialUDP("udp", localAddr, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Print server information
	log.Printf("multicast server started\n")
	log.Printf("multicasting to %s:%s\n", multicastAddr, multicastPort)

	var tno uint = 0
	// Multicast nodeid, tno and timestamp every second
	for {
		// Get the current timestamp
		timestamp := time.Now().Format(time.RFC3339)

		// Prepare message
		tno++
		data := fmt.Sprintf("%d,%d,%s", nodeId, tno, timestamp)

		// Multicast the data
		_, err := conn.Write([]byte(data))
		if err != nil {
			log.Printf("error while multicasting: %v", err)
		}

		// Wait for 1 second before multicasting the next timestamp
		time.Sleep(1 * time.Second)
	}
}
