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

package plugininterface

import (
	"errors"
	"fmt"

	"github.com/tinyzimmer/gomapper/config"
	"github.com/tinyzimmer/gomapper/formats"
	"github.com/tinyzimmer/gomapper/gomapperplugins"
	"github.com/tinyzimmer/gomapper/logging"

	"github.com/tinyzimmer/gomapper/gomapperplugins/awsvpc"
	"github.com/tinyzimmer/gomapper/gomapperplugins/nmap"
)

type LoadedPlugins struct {
	Plugins []GomapperPlugin
}

type GomapperPlugin struct {
	Name      string
	Methods   []string
	Interface *gomapperplugins.PluginInterface
	Config    map[string]interface{}
}

func LoadPlugins(conf config.Configuration) (loadedPlugins LoadedPlugins) {
	for _, mod := range conf.EnabledPlugins {
		if mod != "" {
			loadedPlugin := GomapperPlugin{}
			loadedPlugin.Name = mod
			methods, _, err := loadPluginInterface(mod)
			if err != nil {
				logging.LogError(fmt.Sprintf("engines: Failed to load plugin: %s", mod))
			} else {
				c, _ := conf.Plugins[mod]
				if c != nil {
					c["debug"] = conf.Debug
				} else {
					c = make(map[string]interface{})
					c["debug"] = conf.Debug
				}
				loadedPlugin.Config = c
				loadedPlugin.Methods = methods
				loadedPlugins.Plugins = append(loadedPlugins.Plugins, loadedPlugin)
				logging.LogInfo(fmt.Sprintf("engines: Loaded plugin: %s", mod))
				logging.LogInfo(fmt.Sprintf("engines: %s methods: %s", mod, loadedPlugin.Methods))
			}
		}
	}
	return
}

func loadPluginInterface(mod string) (methods []string, interf gomapperplugins.PluginInterface, err error) {
	switch mod {
	case "nmap":
		methods, interf, err = nmap.LoadPlugin()
	case "awsvpc":
		methods, interf, err = awsvpc.LoadPlugin()
	default:
		err = errors.New("Tried to load invalid mod interface")
	}
	return
}

func (g GomapperPlugin) CheckInterface() ([]string, bool) {
	methods, _, err := g.GetInterface()
	if err != nil {
		return nil, false
	} else {
		return methods, true
	}
}

func (g GomapperPlugin) GetInterface() (methods []string, inter gomapperplugins.PluginInterface, err error) {
	methods, inter, err = loadPluginInterface(g.Name)
	if err != nil {
		g.logError(err)
		return
	}
	return
}

func (g GomapperPlugin) DiscoverNetworks() ([]string, error) {
	_, inter, _ := g.GetInterface()
	networks, err := inter.DiscoverNetworks(g.Config)
	if err != nil {
		return nil, err
	}
	return networks, nil
}

func (g GomapperPlugin) ScanNetwork(network string) (formats.DbNetwork, error) {
	_, inter, _ := g.GetInterface()
	dbNetwork, err := inter.ScanNetwork(g.Config, network)
	return dbNetwork, err
}

func (g GomapperPlugin) HandleScanRequest(req *formats.ReqInput) (interface{}, formats.DbNetwork) {
	_, inter, _ := g.GetInterface()
	response, dbNetwork := inter.HandleScanRequest(g.Config, req)
	return response, dbNetwork
}

func (g GomapperPlugin) logInfo(msg string) {
	logging.LogInfo(fmt.Sprintf("%s: %s", g.Name, msg))
}

func (g GomapperPlugin) logWarn(msg string) {
	logging.LogWarn(fmt.Sprintf("%s: %s", g.Name, msg))
}

func (g GomapperPlugin) logError(err error) {
	logging.LogError(fmt.Sprintf("%s: %s", g.Name, err.Error()))
}
