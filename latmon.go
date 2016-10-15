package main

import (
	"fmt"
	"time"
	"flag"
	"strings"
	"github.com/jerrinot/latencymon/server"
	"github.com/jerrinot/latencymon/client"
	"strconv"
)

const DEFAULT_PORT = 3540
const DEFAULT_INTERVAL = 1000

func main() {
	hostArgument := flag.String("s", "localhost", "comma separated list of servers to connect to")
	port := flag.Int("p", DEFAULT_PORT, "tcp port to use")
	interval := flag.Int("i", DEFAULT_INTERVAL, "interval in ms")
	flag.Parse()

	server.StartServer(*port)
	hosts := strings.Split(*hostArgument, ",")

	chans  := startClients(hosts, *port)
	measure(chans, *interval)
}

func startClients(hosts []string, port int) []chan int {
	chans := make([]chan int, len(hosts))
	csvFormat := "# timestamp"
	for i, host := range hosts {
		csvFormat += ", "
		chans[i] = make(chan int)
		csvFormat += host
		address := fmt.Sprintf("%s:%d", host, port)
		go client.Measure(address, chans[i])
	}
	fmt.Println(csvFormat)
	return chans
}

func measure(chans []chan int, interval int) {
	for {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		var csv string
		sep := ""
		for _, latChan := range chans {
			csv += sep
			latency := <- latChan
			csv += strconv.Itoa(latency)
			sep = ", "
		}
		fmt.Printf("%v, %s\n", time.Now().Unix(), csv)
	}
}