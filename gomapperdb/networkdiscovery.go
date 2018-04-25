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
	"github.com/tinyzimmer/gomapper/logging"
	"github.com/tinyzimmer/gomapper/netutils"
	"github.com/tinyzimmer/gomapper/plugininterface"
)

func notifyDiscoveryDisabled() {
	logging.LogWarn("Network discovery is disabled")
}

func SetupNetworkDiscovery(plugins plugininterface.LoadedPlugins) (db MemoryDatabase, err error) {
	addr, err := netutils.GetAddr()
	if err != nil {
		logging.LogError(err.Error())
		notifyDiscoveryDisabled()
		return
	}
	db, err = GetMemoryDatabase(plugins, addr)
	if err != nil {
		logging.LogError(err.Error())
		notifyDiscoveryDisabled()
		return
	}
	return
}
