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

package baseplugin

import (
	"github.com/tinyzimmer/gomapper/formats"
	"github.com/tinyzimmer/gomapper/gomapperplugins"
)

type BasePlugin struct{}

func LoadPlugin() ([]string, gomapperplugins.PluginInterface, error) {
	methods := []string{}
	p := BasePlugin{}
	return methods, p, nil
}

func (p BasePlugin) ScanNetwork(conf map[string]interface{}, network string) (dbEntry formats.DbNetwork, err error) {
	return
}

func (p BasePlugin) DiscoverNetworks(conf map[string]interface{}) (networks []string, err error) {
	return
}

func (p BasePlugin) HandleScanRequest(conf map[string]interface{}, input *formats.ReqInput) (response interface{}, dbEntry formats.DbNetwork) {
	return
}
