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

package awsvpc

import (
	"github.com/tinyzimmer/gomapper/formats"
	"github.com/tinyzimmer/gomapper/gomapperplugins"
)

type AwsPlugin struct{}

func LoadPlugin() ([]string, gomapperplugins.PluginInterface, error) {
	methods := []string{}
	p := AwsPlugin{}
	return methods, p, nil
}

func (p AwsPlugin) ScanNetwork(conf map[string]interface{}, network string) (dbEntry formats.DbNetwork, err error) {
	return
}

func (p AwsPlugin) DiscoverNetworks(conf map[string]interface{}) (networks []string, err error) {
	return
}

func (p AwsPlugin) HandleScanRequest(conf map[string]interface{}, input *formats.ReqInput) (response interface{}, dbEntry formats.DbNetwork) {
	return
}
