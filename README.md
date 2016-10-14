# LatencyMon
Dead-Simple TCP Full Round-Trip Measurement.

It combines both server and clients in one tool and produces a simple CSV format suitable for further processing. 

## Output Format Description
### Syntax
`unix_timestamp, round_trip_in_nanos_to_first_server, round_trip_in_nanos_to_second_server, ....`
#### Example: 
`1476454782, 72709, 73136`

## Example Usage
Start the latencymon on all servers you want to measure and use the `-s` flag to specify IP addresses. 
I have 2 boxes and I want to measure latency among them. 
```
[jaromir@Server1 ~]$ ./latencymon -s 10.212.40.101,10.212.40.102
# timestamp, 10.212.40.101, 10.212.40.102
# server: connection no. 1: accepted. 10.212.40.101:3540 <-> 10.212.40.101:43824
# server: connection no. 2: accepted. 10.212.40.101:3540 <-> 10.212.40.102:58815
1476454782, 72709, 73136
1476454783, 39054, 52702
1476454784, 31709, 49425
1476454785, 28480, 56950
1476454786, 28715, 44806
```
```
[jaromir@SERVER2 ~]$ ./latencymon -s 10.212.40.101,10.212.40.102
# timestamp, 10.212.40.101, 10.212.40.102
# server: connection no. 1: accepted. 10.212.40.102:3540 <-> 10.212.40.102:54749
1476454780, -1, 98348
# server: connection no. 2: accepted. 10.212.40.102:3540 <-> 10.212.40.101:51083
1476454781, -1, 38944
1476454782, 128347, 32958
1476454783, 46465, 38848
1476454784, 48947, 33906
1476454785, 46224, 26431
```


```
$ ./latencymon -h
Usage of ./latencymon:
  -i int
    	interval in ms (default 1000)
  -p int
    	tcp port to use (default 3540)
  -s string
    	comma separated list of servers to connect to (default "localhost")
```
