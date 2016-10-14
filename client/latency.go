package client

import (
	"net"
	"github.com/aristanetworks/goarista/monotime"
	"encoding/binary"
	"time"
)

type Latency struct {
	address string
	conn net.Conn
}

func (l *Latency) reconnect() {
	var duration time.Duration = 100 * time.Millisecond
	l.conn, _ = net.DialTimeout("tcp", l.address, duration)
}

func (l *Latency) onFailure(latencyChannel chan int) {
	latencyChannel <- -1
	l.reconnect()
}

func Measure(address string, latencyChannel chan int) *Latency {
	l := Latency{address:address}
	l.reconnect()

	for {
		if (l.conn == nil) {
			l.onFailure(latencyChannel)
			continue
		}

		startTime := monotime.Now()
		err := binary.Write(l.conn, binary.LittleEndian, startTime)
		if (err != nil) {
			l.onFailure(latencyChannel)
			continue
		}

		var receivedTime uint64
		err = binary.Read(l.conn, binary.LittleEndian, &receivedTime)
		if (err != nil) {
			l.onFailure(latencyChannel)
			continue
		}

		delta := monotime.Now() - startTime
		latencyChannel <- int(delta)
	}
}