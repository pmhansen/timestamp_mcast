package main

import (
	"fmt"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
)

func main() {
	// Configure multicast address and ports
	multicastAddr := "225.0.0.1"    // IPv4 multicast address
	multicastListenerPort := "3001" // Multicast server port
	multicastClientPort := "3000"   // Multicast client port, where the server is casting
	networkInterface := "wlp4s0"    // Network interface name to use

	if len(os.Args[1:]) != 4 {
		fmt.Println("\nRunning without arguments, using defaults:")
		fmt.Printf("Run with arguments %s <mcast group> <mcast listener port> <mcast client port> <nwtwork interface>\n", os.Args[0])
		fmt.Printf("defaults: %s %s %s %s %s\n\n", os.Args[0], multicastAddr, multicastListenerPort, multicastClientPort, networkInterface)
	} else {
		multicastAddr = os.Args[1]
		multicastListenerPort = os.Args[2]
		multicastClientPort = os.Args[3]
		networkInterface = os.Args[4]
	}

	// Resolve the multicast address
	addr, err := net.ResolveUDPAddr("udp", multicastAddr+":"+multicastListenerPort)
	if err != nil {
		log.Fatal(err)
	}

	// Create a UDP connection
	var conn net.PacketConn
	conn, err = net.ListenPacket("udp4", ":"+multicastClientPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Find the network interface by name
	var nic *net.Interface
	nic, err = net.InterfaceByName(networkInterface)
	if err != nil {
		log.Fatal(err)
	}

	// Join the multicast group
	p := ipv4.NewPacketConn(conn)
	if err := p.JoinGroup(nic, addr); err != nil {
		log.Fatal(err)
	}

	// Print client information
	log.Printf("multicast client started\n")
	log.Printf("listening for multicast data at %s:%s\n", multicastAddr, multicastClientPort)

	// Read multicast packets and process them
	buffer := make([]byte, 128)
	for {
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err)
		}

		// Process the received data
		data := string(buffer[:n])
		fmt.Printf("received data: %s\n", data)
	}
}
