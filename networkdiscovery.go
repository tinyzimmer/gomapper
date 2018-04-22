/**
    This file is part of gomapper.

    Gomapper is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Gomapper is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with gomapper.  If not, see <http://www.gnu.org/licenses/>.
**/

package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"
)

const DEFAULT_TIMEOUT_MS = 500
const DEFAULT_PING_HOST = "1.1.1.1"
const DEFAULT_PORT = 33434
const DEFAULT_PACKET_SIZE = 52
const DEFAULT_MAX_HOPS = 64
const DEFAULT_MAX_RETRIES = 3
const DEFAULT_ASSUMED_NETMASK = "24"

type Hop struct {
	Success     bool
	Address     [4]byte
	Host        string
	N           int
	ElapsedTime time.Duration
	TTL         int
}

func formatDefaultNetmask(ip string) string {
	split := strings.Split(ip, ".")
	return fmt.Sprintf("%s.%s.%s.0/%s", split[0], split[1], split[2], DEFAULT_ASSUMED_NETMASK)
}

func getAddr() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return net.IP{}, err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if len(ipnet.IP.To4()) == net.IPv4len {
				return ipnet.IP, nil
			}
		}
	}
	outErr := errors.New("Failed to retrieve socket address")
	return net.IP{}, outErr
}

func getIpObj(ip string) net.IP {
	return net.ParseIP(ip)
}

func destAddr(dest string) (destAddr [4]byte, err error) {
	addrs, err := net.LookupHost(dest)
	if err != nil {
		return
	}
	addr := addrs[0]

	ipAddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return
	}
	copy(destAddr[:], ipAddr.IP.To4())
	return
}

func isPrivateAddr(addr net.IP) bool {
	var private_nets = [4]string{"192.168.0.0/16", "10.0.0.0/8", "172.0.0.0/8", "127.0.0.0/24"}
	for _, network := range private_nets {
		_, ipnet, err := net.ParseCIDR(network)
		if err != nil {
			logError("Failed to get address membership")
			return false
		}
		if ipnet.Contains(addr) {
			return true
		}
	}
	return false
}

func probeNetwork(graph Graph, network string, config Configuration) {
	logInfo(fmt.Sprintf("Probing network: %s", network))
	scanner, err := InitScanner(network, config.Discovery.Debug)
	if err != nil {
		return
	}
	if config.Discovery.Mode == "ping" {
		logInfo(fmt.Sprintf("Using ping mode on network %s", network))
		scanner.SetPingDiscovery()
	} else if config.Discovery.Mode == "stealth" {
		logInfo(fmt.Sprintf("Using stealth mode on network %s", network))
		scanner.SetStealthDiscovery()
	} else if config.Discovery.Mode == "connect" {
		scanner.SetConnectDiscovery()
	}
	scanner.RunScan()
	if !scanner.Failed {
		numHosts := len(scanner.Results.Hosts)
		logInfo(fmt.Sprintf("Probe of %s complete. Found %v hosts", network, numHosts))
		graph.AddScanResultsByNetwork(network, scanner.Results)
	} else {
		logError(fmt.Sprintf("Scan of %s failed", network))
	}
}

func localNetworkDiscovery(addr net.IP, graph Graph, config Configuration) {
	for _, netw := range config.Discovery.Networks {
		logInfo(fmt.Sprintf("Adding %s to memory graph", netw))
		graph.AddNetwork(netw)
		go probeNetwork(graph, netw, config)
	}
	networks, err := detectLocalNetworks(addr)
	if err != nil {
		logError("Could not detect local networks. Discovery is disabled.")
	} else {
		for _, network := range networks {
			networkString := fmt.Sprintf("%s/%s", network.IP.String(), DEFAULT_ASSUMED_NETMASK)
			logInfo(fmt.Sprintf("Adding %s to memory graph", networkString))
			graph.AddNetwork(networkString)
			go probeNetwork(graph, networkString, config)
		}
	}
}

func notifyDiscoveryDisabled() {
	logWarn("Network discovery is disabled")
}

func setupNetworkDiscovery() (addr net.IP, graph Graph, err error) {
	addr, err = getAddr()
	if err != nil {
		logError(err.Error())
		notifyDiscoveryDisabled()
		return
	}
	graph, err = getMemoryGraph()
	if err != nil {
		logError(err.Error())
		notifyDiscoveryDisabled()
		return
	}
	return
}

func netContains(slice []net.IPNet, item net.IPNet) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		netStr := s.String()
		set[netStr] = struct{}{}
	}
	itemStr := item.String()
	_, ok := set[itemStr]
	return ok
}

func detectLocalNetworks(addr net.IP) ([]net.IPNet, error) {
	var networks []net.IPNet
	destAddr, err := destAddr(DEFAULT_PING_HOST)
	if err != nil {
		logError(fmt.Sprintf("Trace Detection Error: %s", err.Error()))
		return networks, err
	}
	var sourceAddr [4]byte
	copy(sourceAddr[:], addr.To4())
	timeoutMs := (int64)(DEFAULT_TIMEOUT_MS)
	tv := syscall.NsecToTimeval(1000 * 1000 * timeoutMs)
	ttl := 1
	retry := 0
	for {
		recvSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
		if err != nil {
			logWarn("Could not create raw socket for ping probe, are you running in the docker container? If not, do that, or try root.")
			logWarn(fmt.Sprintf("Traceroute Detection Error: %s", err.Error()))
			return networks, err
		}
		sendSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
		if err != nil {
			logWarn(fmt.Sprintf("Traceroute Detection Error: %s", err.Error()))
			return networks, err
		}
		syscall.SetsockoptInt(sendSocket, 0x0, syscall.IP_TTL, ttl)
		syscall.SetsockoptTimeval(recvSocket, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv)

		defer syscall.Close(recvSocket)
		defer syscall.Close(sendSocket)

		syscall.Bind(recvSocket, &syscall.SockaddrInet4{Port: DEFAULT_PORT, Addr: sourceAddr})

		syscall.Sendto(sendSocket, []byte{0x0}, 0, &syscall.SockaddrInet4{Port: DEFAULT_PORT, Addr: destAddr})

		var p = make([]byte, DEFAULT_PACKET_SIZE)
		_, from, err := syscall.Recvfrom(recvSocket, p, 0)
		if err == nil {
			currAddr := from.(*syscall.SockaddrInet4).Addr
			netObj := net.IPv4(currAddr[0], currAddr[1], currAddr[2], byte(0))
			if isPrivateAddr(netObj) {
				network := net.IPNet{IP: netObj, Mask: net.IPv4Mask(255, 255, 255, 0)}
				if !netContains(networks, network) {
					networks = append(networks, network)
					logInfo(fmt.Sprintf("Local Network Detected: %s/%s", network.IP, DEFAULT_ASSUMED_NETMASK))
				}
			}
			ttl += 1
			retry = 0
			if ttl > DEFAULT_MAX_HOPS || currAddr == destAddr {
				return networks, nil
			}
		} else {
			retry += 1
			if retry > DEFAULT_MAX_RETRIES {
				ttl += 1
				retry = 0
			}

			if ttl > DEFAULT_MAX_HOPS {
				return networks, nil
			}
		}
	}
	return networks, nil
}