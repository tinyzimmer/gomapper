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

	"github.com/tinyzimmer/gomapper/config"
	"github.com/tinyzimmer/gomapper/gomapperdb"
	"github.com/tinyzimmer/gomapper/logging"
	"github.com/tinyzimmer/gomapper/netutils"
	"github.com/tinyzimmer/gomapper/plugininterface"
)

func startHttpListener(addr net.IP, port string, db gomapperdb.MemoryDatabase) {
	mux := http.NewServeMux()
	mux.HandleFunc("/scan", db.ReceivedScan)
	mux.HandleFunc("/query", db.IterateNetworks)
	if netutils.IsPrivateAddr(addr) {
		logging.LogInfo(fmt.Sprintf("Listening on private address: %s:%s", addr, port))
	} else {
		logging.LogInfo(fmt.Sprintf("Listening on public address %s:%s", addr, port))
	}
	http.ListenAndServe(fmt.Sprintf("%s:%s", addr, port), mux)
}

func main() {
	var configFile = flag.String("config", "", "toml configuration file")
	flag.Parse()
	config, err := config.GetConfig(configFile)
	if err != nil {
		logging.LogError(err.Error())
		return
	}
	addr := netutils.GetIpObj(config.Server.ListenAddress)
	port := config.Server.ListenPort
	plugins := plugininterface.LoadPlugins(config.Plugins.EnabledPlugins)
	localAddr, db, err := gomapperdb.SetupNetworkDiscovery(plugins)
	if err == nil && config.Discovery.Enabled {
		go gomapperdb.LocalNetworkDiscovery(localAddr, db, plugins, config)
	}
	startHttpListener(addr, port, db)
}
