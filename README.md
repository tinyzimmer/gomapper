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
The container weighs in at about 33.5 MB and will have absolutely nothing but the statically compiled gomapper and nmap binaries, as well as nmap's support documents.

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

$> curl localhost:8080/scan -d '{"target": "127.0.0.1", "detection": "full", "ports": "8080"}'
{
    "Scanner": "nmap",
    "Args": "nmap -A -p 8080 -oX /tmp/425838479 127.0.0.1",
    "Start": 1524178362,
    "StartStr": "Thu Apr 19 22:52:42 2018",
    "Version": "7.70",
    "ScanInfo": {
        "Type": "syn",
        "Protocol": "tcp",
        "NumServices": 1,
        "Services": "8080"
    },
    "Verbose": {
        "Level": 0
    },
    "Debugging": {
        "Level": 0
    },
    "Hosts": [
        {
            "StartTime": 1524178362,
            "EndTime": 1524178369,
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
                        "PortId": 8080,
                        "State": {
                            "State": "open",
                            "Reason": "syn-ack",
                            "ReasonTTL": 64
                        },
                        "Service": {
                            "Name": "http",
                            "Product": "Golang net/http server",
                            "Version": "",
                            "ExtraInfo": "Go-IPFS json-rpc or InfluxDB API",
                            "Method": "probed",
                            "Conf": "10",
                            "Cpe": "cpe:/a:protocol_labs:go-ipfs"
                        },
                        "Script": {
                            "Id": "http-title",
                            "Ouput": "Site doesn't have a title (text/plain; charset=utf-8).",
                            "Tables": null
                        }
                    }
                ]
            },
            "Times": {
                "Srtt": 28,
                "Rttvar": 91,
                "To": 100000
            }
        }
    ],
    "RunStats": {
        "Finished": {
            "Time": 1524178369,
            "TimeStr": "Thu Apr 19 22:52:49 2018",
            "Elapsed": 8.08,
            "Summary": "Nmap done at Thu Apr 19 22:52:49 2018; 1 IP address (1 host up) scanned in 8.08 seconds",
            "Exit": "success"
        },
        "FinishedHosts": {
            "Up": 1,
            "Down": 0,
            "Total": 1
        }
    }
}

# Below examples tested from outside docker container using NSE scripts and service detection

$> curl localhost:8080/scan -d '{"target": "127.0.0.1", "script": "ssh-hostkey", "ports": "22"}'
{
    "Scanner": "nmap",
    "Args": "nmap --script=ssh-hostkey -p 22 -oX /tmp/377426807 127.0.0.1",
    ...
    "Hosts": [
        ...
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
                            "State": "open",
                            "Reason": "syn-ack",
                            "ReasonTTL": 0
                        },
                        "Service": {
                            "Name": "ssh",
                            "Product": "",
                            "Version": "",
                            "ExtraInfo": "",
                            "Method": "table",
                            "Conf": "3",
                            "Cpe": ""
                        },
                        "Script": {
                            "Id": "ssh-hostkey",
                            "Ouput": "\n  2048 7b:24:9e:6d:3a:1e:2d:5d:80:cc:fc:7a:84:a6:76:30 (RSA)\n  256 6d:bd:61:77:f7:0e:c6:18:05:ae:76:b7:48:23:cf:21 (ECDSA)\n  256 90:dc:9a:d0:c9:3a:42:4e:73:0c:51:b6:ce:86:60:08 (ED25519)",
                            "Tables": [
                                {
                                    "Elems": [
                                        {
                                            "Key": "key",
                                            "Value": "AAAAB3NzaC1yc2EAAAADAQABAAABAQDJmIxQQbk9y5Z+YfljWN98MlZnC52jR0cgMtxFItxmDVywkREhKOayjiHdA71+oMsJYKH3iDfO5hudtiDfbA83MfzdBJ6oFtVhUhldxDXb+R68fiHftUUZRuezvOrHUmyzsUE7TozJLIXp9xiew9aNPMZUb4urS1LHlu7irgRlAjgm9oXVB+vxwpcCsURahV6Nnr2cjhUK4/1R5QVrN4hE7sZbXue9FQCLst5jZUrN/KyulHXPijC/L5gknYci53diXv50kYpQ2+k498kS3t7VWGOgKCSQnCWrXrh5pRw3xDqxyMDUf6hYm5QWzKrRfLvEF5uyfFuQWaAvQmF4PRkd"
                                        },
                                        {
                                            "Key": "bits",
                                            "Value": "2048.0"
                                        },
                                        {
                                            "Key": "type",
                                            "Value": "ssh-rsa"
                                        },
                                        {
                                            "Key": "fingerprint",
                                            "Value": "7b249e6d3a1e2d5d80ccfc7a84a67630"
                                        }
                                    ]
                                },
                                ...
                            ]
                        }
                    }
                ]
            },
        ...
    }
}

$> curl localhost:8080/scan -d '{"target": "127.0.0.1", "method": "udp", "script": "snmp-sysdescr", "scriptArgs": "creds.snmp=public", "ports": "161"}'
{
    "Scanner": "nmap",
    ...
    "Hosts": [
        ...
            "Ports": {
                ...
                "Ports": [
                    {
                        "Protocol": "udp",
                        "PortId": 161,
                        "State": {
                            "State": "open",
                            "Reason": "udp-response",
                            "ReasonTTL": 64
                        },
                        "Service": {
                            "Name": "snmp",
                            "Product": "",
                            "Version": "",
                            "ExtraInfo": "",
                            "Method": "table",
                            "Conf": "3",
                            "Cpe": ""
                        },
                        "Script": {
                            "Id": "snmp-sysdescr",
                            "Ouput": "Linux base 4.16.2-1-MANJARO #1 SMP PREEMPT Thu Apr 12 17:46:07 UTC 2018 x86_64\n  System uptime: 5m23.92s (32392 timeticks)",
                            "Tables": null
                        }
                    }
                ]
            },
            ...
}
```

You can also use the docker container as a standalone nmap installation

```bash
$> docker run --rm tinyzimmer/gomapper /bin/nmap
$> # or
$> alias nmap='docker run --rm tinyzimmer/gomapper /bin/nmap'
$> nmap
```
