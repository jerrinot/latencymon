package server

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

func StartServer(port int) {
	go func() {
		server, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if server == nil {
			panic("server: couldn't start listening: " + err.Error())
		}
		conns := clientConns(server)
		for {
			go handleConn(<-conns)
		}
	}()

}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Printf("# server: couldn't accept: " + err.Error())
				continue
			}
			i++
			fmt.Printf("# server: connection no. %d: accepted. %v <-> %v\n", i,
				client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(conn net.Conn) {
	for {
		var time uint64
		err := binary.Read(conn, binary.LittleEndian, &time)
		if err != nil {
			onConnectionError(conn, err)
			break
		}
		err = binary.Write(conn, binary.LittleEndian, time)
		if err != nil {
			onConnectionError(conn, err)
			break
		}
	}
}

func onConnectionError(conn net.Conn, err error) {
	fmt.Println("# server: Connection to " + conn.RemoteAddr().String() + " lost. Reason: " + err.Error())
}
