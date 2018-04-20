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
	"log"
	"net"
)

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

func isPrivateAddr(addr net.IP) bool {
	var private_nets = [3]string{"192.168.0.0/8", "10.0.0.0/8", "172.0.0.0/8"}
	for _, network := range private_nets {
		_, ipnet, err := net.ParseCIDR(network)
		if err != nil {
			log.Println("Failed to get address membership")
			return false
		}
		if ipnet.Contains(addr) {
			return true
		}
	}
	return false
}
