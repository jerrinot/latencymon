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
const DEFAULT_MODE = "hybrid"
const DEFAULT_BACKLOG = 1000

func main() {
	hostArgument := flag.String("s", "localhost", "comma separated list of servers to connect to")
	port := flag.Int("p", DEFAULT_PORT, "tcp port to use")
	interval := flag.Int("i", DEFAULT_INTERVAL, "interval in ms")
	mode := flag.String("m", DEFAULT_MODE, "Operation mode: {hybrid|server|client}")
	flag.Parse()

	switch *mode {
	case "hybrid", "client", "server":
	default:
		fmt.Printf("Unknown mode '%s'. Valid modes: hybrid|server|client\n", *mode)
		return
	}

	if (*mode != "client") {
		server.StartServer(*port)
	}

	if (*mode != "server") {
		hosts := strings.Split(*hostArgument, ",")
		clientChans := startClients(hosts, *port)
		timerChan := startTimer(*interval, DEFAULT_BACKLOG)
		go startMeasuring(clientChans, timerChan)
	}

	sleepIndefinitely()
}


func sleepIndefinitely() {
	select{}
}

func startTimer(rateMs int, backlog int) chan int {
	chn := make(chan int, backlog)
	go func() {
		for i := 0; ; i++ {
			chn <- i
			time.Sleep(time.Duration(rateMs) * time.Millisecond)
		}
	}()
	return chn
}

func startClients(hosts []string, port int) []chan int {
	chans := make([]chan int, len(hosts))
	csvFormat := "# timestamp"
	for i, host := range hosts {
		csvFormat += ", "
		chans[i] = make(chan int)
		csvFormat += host
		if strings.Index(host, ":") == -1 {
			host = fmt.Sprintf("%s:%d", host, port)
		}
		go client.Measure(host, chans[i])
	}
	fmt.Println(csvFormat)
	return chans
}

func startMeasuring(clientChans []chan int, timerChan chan int) {
	for {
		<- timerChan
		var csv string
		sep := ""
		for _, latChan := range clientChans {
			csv += sep
			latency := <- latChan
			csv += strconv.Itoa(latency)
			sep = ", "
		}
		fmt.Printf("%v, %s\n", time.Now().Unix(), csv)
	}
}