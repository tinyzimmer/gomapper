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
	"syscall"
	"time"
)

const DEFAULT_TIMEOUT_MS = 500
const DEFAULT_PING_HOST = "1.1.1.1"
const DEFAULT_PORT = 33434
const DEFAULT_PACKET_SIZE = 52
const DEFAULT_MAX_HOPS = 64
const DEFAULT_MAX_RETRIES = 3

type Hop struct {
	Success     bool
	Address     [4]byte
	Host        string
	N           int
	ElapsedTime time.Duration
	TTL         int
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
	out_err := errors.New("Failed to retrieve socket address")
	return net.IP{}, out_err

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
	var private_nets = [3]string{"192.168.0.0/8", "10.0.0.0/8", "172.0.0.0/8"}
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

func detectLocalNetworks(addr net.IP) ([]net.IPNet, error) {
	var networks []net.IPNet
	destAddr, err := destAddr(DEFAULT_PING_HOST)
	if err != nil {
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
			return networks, err
		}
		sendSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
		if err != nil {
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
				networks = append(networks, network)
				logInfo(fmt.Sprintf("Local Network Detected: %s/24", network.IP))
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
