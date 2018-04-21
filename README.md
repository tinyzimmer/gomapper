# gomapper
REST-like interface in go for running Nmap scans

Actually, I am also turning it into a queryable, passive network mapper that stores it's data in memory graphs

# Building/Downloading

To build the docker container which will include an nmap installation:

*Using the Docker container has the benefit of providing access to root only nmap functions*

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

To compile and build locally:

```bash
# Requires local installation of nmap to use
$> go get github.com/tinyzimmer/gomapper
```

# Running

![Server side](doc/server.apng)

And from the client

![Client side](doc/client.apng)

Eventually more information will get transparently stored in memory from various probes

```bash
$> # Start the service via the command line or docker
$> docker run -p 8080:8080 tinyzimmer/gomapper
$> # Leave off ports for default options. Todo is to create more argument generation functions for different scan types and switches
$> # use "rawArgs" (list) to create a custom scan instead of method
$> # Current scan methods: ["tcp-connect", "tcp-ack", "tcp-syn", "udp", "ping"]
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
```
