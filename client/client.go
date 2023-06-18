package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", ":50000")
	checkConnection(conn, err)
	if err != nil {

		// No connection could be established - con refused
		fmt.Println("Error dialing", err.Error())
		return
	}
	inr := bufio.NewReader(os.Stdin)
	fmt.Println("Enter client Identification")
	cid, _ := inr.ReadString('\n')
	fmt.Printf("Client ID: %s", cid)
	trimmedCid := strings.Trim(cid, "\r\n")

	// send to server until termination
	for {
		fmt.Println("Enter input to server or Q to quit.")
		in, _ := inr.ReadString('\n')
		trimmedIn := strings.Trim(in, "\r\n")
		if trimmedIn == "Q" {
			conn.Close()
			return
		}
		_, err = conn.Write([]byte(trimmedCid + " sent: " + trimmedIn))
		if err != nil {
			fmt.Printf("Error writing to server: %v", err)
		}
	}

}

func checkConnection(conn net.Conn, err error) {
	if err != nil {
		fmt.Printf("error %v connecting!", err)
		os.Exit(1)
	}
	fmt.Printf("Connection is made with %s\n", conn.RemoteAddr().String())
}
