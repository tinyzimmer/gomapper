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

package gomapperplugins

import (
	"net"

	"github.com/tinyzimmer/gomapper/formats"
)

type PluginInterface interface {
	DiscoverNetworks(map[string]interface{}) ([]net.IPNet, error)
	ScanNetwork(map[string]interface{}, string) (formats.DbNetwork, error)
	HandleScanRequest(map[string]interface{}, *formats.ReqInput) (interface{}, formats.DbNetwork)
}
