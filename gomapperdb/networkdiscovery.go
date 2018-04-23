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

package gomapperdb

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/tinyzimmer/gomapper/config"
	"github.com/tinyzimmer/gomapper/logging"
	"github.com/tinyzimmer/gomapper/netutils"
	"github.com/tinyzimmer/gomapper/scanner"
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

func GetAddr() (net.IP, error) {
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
			logging.LogError("Failed to get address membership")
			return false
		}
		if ipnet.Contains(addr) {
			return true
		}
	}
	return false
}

func probeNetwork(db MemoryDatabase, network string, config config.Configuration) {
	logging.LogInfo(fmt.Sprintf("Probing network: %s", network))
	scanner, err := scanner.InitScanner(network, config.Discovery.Debug)
	if err != nil {
		return
	}
	logging.LogInfo(fmt.Sprintf("Using %s mode on network %s", config.Discovery.Mode, network))
	if config.Discovery.Mode == "ping" {
		scanner.SetPingDiscovery()
	} else if config.Discovery.Mode == "stealth" {
		scanner.SetStealthDiscovery()
	} else if config.Discovery.Mode == "connect" {
		scanner.SetConnectDiscovery()
	}
	scanner.RunScan()
	if !scanner.Failed {
		numHosts := len(scanner.Results.Hosts)
		logging.LogInfo(fmt.Sprintf("Probe of %s complete. Found %v hosts", network, numHosts))
		db.AddScanResultsByNetwork(network, scanner.Results)
	} else {
		logging.LogError(fmt.Sprintf("Scan of %s failed", network))
	}
}

func LocalNetworkDiscovery(addr net.IP, db MemoryDatabase, config config.Configuration) {
	for _, netw := range config.Discovery.Networks {
		logging.LogInfo(fmt.Sprintf("Adding %s to memory database", netw))
		db.AddNetwork(netw)
		go probeNetwork(db, netw, config)
	}
	networks, err := netutils.DetectLocalNetworks(addr)
	if err != nil {
		logging.LogError("Could not detect local networks. Discovery is disabled.")
	} else {
		for _, network := range networks {
			networkString := fmt.Sprintf("%s/%s", network.IP.String(), DEFAULT_ASSUMED_NETMASK)
			logging.LogInfo(fmt.Sprintf("Adding %s to memory database", networkString))
			db.AddNetwork(networkString)
			go probeNetwork(db, networkString, config)
		}
	}
}

func notifyDiscoveryDisabled() {
	logging.LogWarn("Network discovery is disabled")
}

func SetupNetworkDiscovery() (addr net.IP, db MemoryDatabase, err error) {
	addr, err = netutils.GetAddr()
	if err != nil {
		logging.LogError(err.Error())
		notifyDiscoveryDisabled()
		return
	}
	db, err = GetMemoryDatabase()
	if err != nil {
		logging.LogError(err.Error())
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
