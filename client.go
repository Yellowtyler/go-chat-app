package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("connected")
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- struct{}{}
	}()

	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Fatalln(err)
	}
	conn.Close()
	<-done
}
