package main

import (
	"flag"
	"fmt"
	"net"
	"syscall"
)

var maxRead = 25

func main() {
	fmt.Println("Starting the server ...")
	// create listener:
	// listener, err := net.Listen("tcp", "0.0.0.0:3001")
	// local listener
	flag.Parse()
	if flag.NArg() != 2 {
		panic("Required usage: host port")
	}
	addr := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1))
	listener := initServer(addr)

	// listener, err := net.Listen("tcp", "localhost:50000")
	for {
		conn, err := listener.Accept()
		if err != nil {
			checkErr(err, "Accept")
		}
		go connectionHandler(conn)
	}
}

func initServer(addr string) *net.TCPListener {
	serverAddr, err := net.ResolveTCPAddr("tcp", addr)
	checkErr(err, "Resolving address:port failed: `"+addr+"`")
	listener, err := net.ListenTCP("tcp", serverAddr)
	checkErr(err, "ListenTCP: ")
	fmt.Println("Listening on: ", listener.Addr().String())
	return listener
}

func connectionHandler(conn net.Conn) {
	connFrom := conn.RemoteAddr().String()
	fmt.Println("Connection from: ", connFrom)
	prompt(conn)

	for {
		var inbuf []byte = make([]byte, maxRead+1)
		length, err := conn.Read(inbuf[0:maxRead])
		// prevent possibile overflow
		inbuf[maxRead] = 0
		switch err {
		case nil:
			handleMsg(length, err, inbuf)
		case syscall.EAGAIN:
			continue
		default:
			_ = conn.Close()
			fmt.Println("Connection terminated: ", connFrom)

			if err != nil {
				fmt.Println("Error closing conncetion: ", err)
				return
			}
			// checkErr(err, "Close: ")
		}
	}
}

func prompt(to net.Conn) {
	outbuf := []byte{}
	outbuf = append(outbuf, "Enter Input"...)
	wrote, err := to.Write(outbuf)
	checkErr(err, "Write: "+string(wrote)+" bytes")
}

func handleMsg(length int, err error, msg []byte) {
	if length > 0 {
		fmt.Print("<", length, ":")
		for i := 0; ; i++ {
			if msg[i] == 0 {
				break
			}
			fmt.Printf("%c", msg[i])
		}
		fmt.Print(">")
	}
}

func checkErr(error error, info string) {
	if error != nil {
		// panic("ERROR: " + info + error.Error())
		fmt.Printf("ERROR: "+info+"%s", error.Error())
		return // not panicing to not kill server on each error
	}
}
