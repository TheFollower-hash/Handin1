package main


import (
"bufio";
"fmt";
"net";
"os"
"strconv"
"strings"
"sync"
)
// vi slemme piger vi til slemme drenge.

var lock sync.RWMutex
var connections []net.Conn
var ln net.Listener
var MessagesSent map[string]bool

func main() {
	connections = make([]net.Conn, 0, 10) // Holds references to established connections
	MessagesSent = make(map[string]bool, 100)
	//Ask for IP and Port on peer
	ipPort := promptForIpPort()

	//Try to Connect
	conn, err := net.Dial("tcp", ipPort)
	if err == nil {
		fmt.Println("SUCCESSFULLY connected to "+ipPort)
		defer conn.Close()

		ln, _ = net.Listen("tcp", conn.LocalAddr().String())
		defer ln.Close()
		go readFrom(conn)

		connections = append(connections, conn)
	} else {
		fmt.Println("The attempt to connect to "+ipPort+" was UNSUCCESSFUL")
		//Create new network, with only itself as member }
		ln, _ = net.Listen("tcp", ":")
		defer ln.Close()
		fmt.Println("A new network has been created")
	}

	// Print own IP address and listening port
	printMyIP()

	printListeningPort(ln)

	//wait for connections
	go waitForConnections(ln)

	// Iteratively prompt the user for text strings
	fmt.Println("Write a message to all peers. Press enter to send:")
	fmt.Print("> ")
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		fmt.Print("> ")
		checkIfMessageHasBeenSentIfNotPrintsAndSents(msg)
	}

	// make sure that When the user enters a text string at any connected client, then it will eventually be printed at all other clients. (Note there should only be the message - no info about who the sender is)
}

func promptForIpPort() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Supply IP-address and Port [xxx.xxx.xxx.xxx:xxxxx]")
	fmt.Print("> ")
	ipPort, _ := reader.ReadString('\n')
	return strings.TrimSpace(ipPort)
}

func printMyIP() {
	name, _ := os.Hostname()
	addrs, _ := net.LookupHost(name)
	fmt.Println("Name: " + name)
	for indx, addr := range addrs {
		fmt.Println("Address number " + strconv.Itoa(indx) + ": " + addr)
	}
}

func printListeningPort(ln net.Listener) {
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	fmt.Println("Port number " + port)
}

func waitForConnections(ln net.Listener) {
	for {
		fmt.Println("Listening for connection...")
		conn, _ := ln.Accept()
		fmt.Println("Got a connection.")
		connections = append(connections, conn)
		go readFrom(conn)
	}
}

func readFrom(conn net.Conn) {
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil { return }
		checkIfMessageHasBeenSentIfNotPrintsAndSents(msg)
	}
}

func checkIfMessageHasBeenSentIfNotPrintsAndSents(msg string)  {
	lock.Lock() //lock MessagesSent
	fmt.Println(MessagesSent)
	if !MessagesSent[msg] {
		MessagesSent[msg] = true
		lock.Unlock() //unlock MessagesSent
		fmt.Println(msg)
		for _, conn := range connections[:cap(connections)] {
			if conn != nil {
				fmt.Fprintf(conn, msg)
			}
		}
	} else {
		lock.Unlock() //unlock MessagesSent
	}
}
