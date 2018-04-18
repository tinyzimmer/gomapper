# gomapper
REST-like interface in go for running Nmap scans

# Building

To build the docker container which will include an nmap installation:

*Using the Docker container has the benefit of providing access to root only scans*

```bash
$> git clone https://www.github.com/tinyzimmer/gomapper
$> cd gomapper
$> ./build.sh
$> docker run -p 8080:8080 gomapper
```

If I keep improving this I'll also try to keep an updated image on dockerhub

```bash
$> docker pull tinyzimmer/gomapper
$> docker run -p 8080:8080 tinyzimmer/gomapper
```

To compile and build locally:

```bash
# Requires local installation of nmap to use
$> go get github.com/tinyzimmer/gomapper
```

# Running 

```bash
$> # Start the service via the command line or docker
$> docker run -p 8080:8080 tinyzimmer/gomapper
$> # Leave off ports for default options. Todo is to create more argument generation functions for different scan types and switches
$> # use "rawArgs" (list) to create a custom scan instead of method
$> # Current scan methods: ["tcp-connect", "tcp-ack", "tcp-syn", "udp"]
$> curl localhost:8080/scan -d '{"target": "127.0.0.1", "method": "tcp-connect", "ports": "22,8080"}'
{
    "Scanner": "nmap",
    "Args": "nmap -sT -p 22,8080 -oX /tmp/731340889 127.0.0.1",
    "Start": 1524018327,
    "StartStr": "Wed Apr 18 02:25:27 2018",
    "Version": "7.60",
    "ScanInfo": {
        "Type": "connect",
        "Protocol": "tcp",
        "NumServices": 2,
        "Services": "22,8080"
    },
    "Verbose": {
        "Level": 0
    },
    "Debugging": {
        "Level": 0
    },
    "Hosts": [
        {
            "StartTime": 1524018327,
            "EndTime": 1524018327,
            "Status": {
                "State": "up",
                "Reason": "localhost-response",
                "ReasonTTL": 0
            },
            "Address": {
                "Addr": "127.0.0.1",
                "AddrType": "ipv4"
            },
            "Hostnames": {
                "Hostnames": [
                    {
                        "Name": "localhost",
                        "Type": "PTR"
                    }
                ]
            },
            "Ports": {
                "ExtraPorts": {
                    "State": "",
                    "Count": 0,
                    "ExtraReasons": null
                },
                "Ports": [
                    {
                        "Protocol": "tcp",
                        "PortId": 22,
                        "State": {
                            "State": "closed",
                            "Reason": "conn-refused",
                            "ReasonTTL": 0
                        },
                        "Service": {
                            "Name": "ssh",
                            "Method": "table",
                            "Conf": "3"
                        }
                    },
                    {
                        "Protocol": "tcp",
                        "PortId": 8080,
                        "State": {
                            "State": "open",
                            "Reason": "syn-ack",
                            "ReasonTTL": 0
                        },
                        "Service": {
                            "Name": "http-proxy",
                            "Method": "table",
                            "Conf": "3"
                        }
                    }
                ]
            },
            "Times": {
                "Srtt": 61,
                "Rttvar": 3757,
                "To": 100000
            }
        }
    ],
    "RunStats": {
        "Finished": {
            "Time": 1524018327,
            "TimeStr": "Wed Apr 18 02:25:27 2018",
            "Elapsed": 0.06,
            "Summary": "Nmap done at Wed Apr 18 02:25:27 2018; 1 IP address (1 host up) scanned in 0.06 seconds",
            "Exit": "success"
        },
        "FinishedHosts": {
            "Up": 1,
            "Down": 0,
            "Total": 1
        }
    }
}
```
