package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var clients = make(map[string]net.Conn)
var messages = make(chan Message)
var leaving = make(chan Message)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("connection error")
	}

	go broadcast()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handle(conn)
	}
}

type Message struct {
	address string
	message string
}

func handle(conn net.Conn) {

	address := conn.RemoteAddr().String()
	clients[address] = conn
	messages <- newMessage("joined", conn)
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- newMessage(address+": "+input.Text(), conn)
	}
	delete(clients, address)
	leaving <- newMessage(address+" left chat", conn)

	conn.Close()
}

func newMessage(message string, conn net.Conn) Message {
	address := conn.RemoteAddr().String()
	return Message{
		address: address,
		message: message,
	}
}

func broadcast() {
	for {
		select {
		case msg := <-messages:
			for _, conn := range clients {
				if msg.address == conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprintln(conn, msg.message)
			}
		case leave := <-leaving:
			for _, conn := range clients {
				fmt.Fprintln(conn, leave.message)
			}
		}
	}
}
