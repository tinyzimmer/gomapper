# gomapper
REST-like interface in go for running Nmap scans

Actually, I am also turning it into a queryable, passive network mapper that stores it's data in a go-memdb

# Building/Downloading

To build the docker container which will include an nmap installation:

*Using the Docker container has the benefit of providing access to root only nmap functions*

## Build Dependencies
* internet connection (required)
* docker (required)
* golang >=1.10 (required)
* upx (optional) - golang binaries apparently compress a crapload

## Build/Pull

```bash
$> git clone https://www.github.com/tinyzimmer/gomapper
$> cd gomapper
$> ./build.sh
$> docker run -p 8080:8080 gomapper
```
The container weighs in at about 30 MB and will have absolutely nothing but the statically compiled gomapper and nmap binaries, as well as nmap's support documents.

If I keep improving this I'll also try to keep an updated image on dockerhub

```bash
$> docker pull tinyzimmer/gomapper
$> docker run -p 8080:8080 tinyzimmer/gomapper
```

You can also use the docker container as a standalone nmap installation

```bash
$> docker run --rm tinyzimmer/gomapper /bin/nmap
$> # or
$> alias nmap='docker run --rm tinyzimmer/gomapper /bin/nmap'
$> nmap
```

## Building without docker

To compile and build locally:

```bash
# Requires local installation of nmap to use
$> go get github.com/tinyzimmer/gomapper
$> go install github.com/tinyzimmer/gomapper
```

If you run this way without root, network discovery will be disabled

# Running

![Server side](doc/server.apng)

And from the client

![Client side](doc/client.apng)

Eventually more information will get transparently stored in memory from various probes

## Configuration

Configuration is done either via a config.toml or environment variables.

See the included example configuration file. These are for use outside the container and can be specified with:

```bash
$> gomapper --config /path/to/config.toml
```

The configurations that can be passed via the environment are below:

|Environment Variable|Options|Default|
|----------|:-------------:|------:|
|GOMAPPER_LISTEN_ADDRESS|ip address|First non-local interface found|
|GOMAPPER_LISTEN_PORT|port|8080|
|GOMAPPER_DISCOVERY_ENABLED|0,1,false,true|true|
|GOMAPPER_DISCOVERY_MODE|ping,stealth,connect|ping|
|GOMAPPER_DISCOVERY_NETWORKS|comma separated list of networks or ip addresses|none|
|GOMAPPER_DISCOVERY_DEBUG|0,1,false,true|false|

Network discovery by default will only act on private networks. Override this behavior by specifying additional networks

## Commands

```bash
$> # Start the service via the command line or docker
$> docker run -p 8080:8080 tinyzimmer/gomapper
$> # Leave off ports for default options. Todo is to create more argument generation functions for different scan types and switches
$> # use "rawArgs" (list) to create a custom scan instead of method
$> # Current scan methods: ["tcp-connect", "tcp-ack", "tcp-syn", "udp", "ping"]
$> curl localhost:8080/scan -d '{"target": "127.0.0.1"}' # defaults
$> curl localhost:8080/scan -d '{"target": "127.0.0.1", "method": "tcp-connect", "ports": "22,8080"}'
$> curl localhost:8080/scan -d '{"target": "127.0.0.1", "detection": "full", "ports": "8080"}'
$> curl localhost:8080/scan -d '{"target": "127.0.0.1", "rawArgs": ["-f", "--data-length", "200", "-T3"]}'

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
$> curl localhost:8080/query # will eventually contain service/port discovery information
{
    "10.0.1.0/24": {
        "Hosts": [
            "10.0.1.24",
            "10.0.1.31",
            "10.0.1.40",
            "10.0.1.86",
            "10.0.1.93",
            "10.0.1.124",
            "10.0.1.126",
            "10.0.1.127",
            "10.0.1.141",
            "10.0.1.164",
            "10.0.1.1"
        ]
    },
    "10.7.7.0/24": {
        "Hosts": [
            "10.7.7.1",
            "10.7.7.86"
        ]
    },
    "192.168.1.0/24": {
        "Hosts": [
            "192.168.1.1",
            "192.168.1.2",
            "192.168.1.100",
            "192.168.1.7"
        ]
    }
}
```
