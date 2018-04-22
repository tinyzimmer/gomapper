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
	"flag"
	"fmt"
	"net"
	"net/http"
)

func startHttpListener(addr net.IP, port string, graph Graph) {
	mux := http.NewServeMux()
	mux.HandleFunc("/scan", graph.receivedScan)
	mux.HandleFunc("/query", graph.IterateNetworks)
	if isPrivateAddr(addr) {
		logInfo(fmt.Sprintf("Listening on private address: %s:%s", addr, port))
	} else {
		logInfo(fmt.Sprintf("Listening on public address %s:%s", addr, port))
	}
	http.ListenAndServe(fmt.Sprintf("%s:%s", addr, port), mux)
}

func main() {
	var configFile = flag.String("config", "", "toml configuration file")
	flag.Parse()
	config, err := getConfig(configFile)
	if err != nil {
		logError(err.Error())
		return
	}
	addr := getIpObj(config.Server.ListenAddress)
	port := config.Server.ListenPort
	localAddr, graph, err := setupNetworkDiscovery()
	if err == nil && config.Discovery.Enabled {
		go localNetworkDiscovery(localAddr, graph, config)
	}
	startHttpListener(addr, port, graph)
}
