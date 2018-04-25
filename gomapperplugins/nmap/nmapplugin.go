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

package nmap

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/tinyzimmer/gomapper/formats"
	"github.com/tinyzimmer/gomapper/gomapperplugins"
	"github.com/tinyzimmer/gomapper/logging"
	"github.com/tinyzimmer/gomapper/netutils"
)

const NMAP_DISCOVERY_ENV_VAR = "GOMAPPER_NMAP_DISCOVERY_MODE"
const NMAP_DEFAULT_DISCOVERY_MODE = "nmap_ping"

type NmapPlugin struct{}

func LoadPlugin() ([]string, gomapperplugins.PluginInterface, error) {
	methods := []string{"nmap_connect", "nmap_syn", "nmap_ack", "nmap_udp", "nmap_ping"}
	p := NmapPlugin{}
	return methods, p, nil
}

func (p NmapPlugin) ScanNetwork(conf map[string]interface{}, network string) (dbEntry formats.DbNetwork, err error) {
	scanner, err := InitScanner(network, false)
	if err != nil {
		return
	}
	def, ok := conf["discovery_mode"].(string)
	if !ok {
		def = checkEnvMode()
	}
	debug, _ := conf["debug"].(bool)
	scanner.SetDebug(debug)
	if debug {
		logging.LogDebug(fmt.Sprintf("nmap: %s", conf))
	}
	logging.LogInfo(fmt.Sprintf("nmap: Using %s mode on network %s", def, network))
	if def == "nmap_ping" {
		scanner.SetPingDiscovery()
	} else if def == "nmap_syn" {
		scanner.SetSynDiscovery()
	} else if def == "nmap_connect" {
		scanner.SetConnectDiscovery()
	} else if def == "nmap_ack" {
		scanner.SetAckDiscovery()
	} else if def == "nmap_udp" {
		scanner.SetUdpDiscovery()
	}
	scanner.RunScan()
	if !scanner.Failed {
		numHosts := len(scanner.Results.Hosts)
		logging.LogInfo(fmt.Sprintf("nmap: Probe of %s complete. Found %v hosts", network, numHosts))
		dbEntry = ParseScanToDbNetwork(network, scanner.Results)
		return
	} else {
		err = errors.New("Scan Failed")
		logging.LogError(fmt.Sprintf("nmap: Scan of %s failed", network))
		return
	}
	return
}

func (p NmapPlugin) DiscoverNetworks(conf map[string]interface{}) (networks []net.IPNet, err error) {
	debug, _ := conf["debug"].(bool)
	if debug {
		logging.LogDebug(fmt.Sprintf("nmap: %s", conf))
		logging.LogDebug("nmap: Finding local address")
	}
	addr, err := netutils.GetAddr()
	if debug {
		logging.LogDebug("nmap: Got local address")
	}
	if err != nil {
		return
	}
	networks, err = netutils.DetectLocalNetworks(addr)
	return
}

func (p NmapPlugin) HandleScanRequest(conf map[string]interface{}, input *formats.ReqInput) (response interface{}, dbEntry formats.DbNetwork) {
	logging.LogInfo("nmap: Parsing nmap request")
	debug, _ := conf["debug"].(bool)
	if debug {
		logging.LogDebug(fmt.Sprintf("nmap: %s", conf))
		logging.LogDebug("nmap: initiating request scanner")
	}
	scanner, err := RequestScanner(input, conf)
	if err != nil {
		logging.LogError(err.Error())
		response = err
		return
	} else {
		logging.LogInfo("nmap: Initiating nmap scan")
	}
	scanner.SetDebug(debug)
	scanner.RunScan()
	if scanner.Failed {
		logging.LogWarn("nmap: User requested scan failed")
		response = scanner.Error
	} else {
		logging.LogInfo("nmap: Returning scan results")
		response = scanner.Results
		dbEntry = ParseScanToDbNetwork(scanner.Target, scanner.Results)
	}
	return
}

func ParseScanToDbNetwork(network string, results *NmapRun) (net formats.DbNetwork) {
	net.Subnet = network
	for _, host := range results.Hosts {
		ip, mac := getScanAddrs(host.Addresses)
		dbHost := formats.DbHost{}
		dbHost.IP = ip
		dbHost.MAC = mac
		for _, port := range host.Ports.Ports {
			service := formats.DbService{}
			service.Port = fmt.Sprintf("%s/%s", port.PortId, port.Protocol)
			service.Name = port.Service.Name
			dbHost.Services = append(dbHost.Services, service)
		}
		net.Hosts = append(net.Hosts, dbHost)
	}
	return
}

func checkEnvMode() (mode string) {
	mode = os.Getenv(NMAP_DISCOVERY_ENV_VAR)
	if mode == "" {
		return NMAP_DEFAULT_DISCOVERY_MODE
	}
	return
}

func getScanAddrs(addrs []Address) (ip string, mac string) {
	for _, addr := range addrs {
		if addr.AddrType == "ipv4" {
			ip = addr.Addr
		} else if addr.AddrType == "mac" {
			mac = addr.Addr
		}
	}
	return
}
