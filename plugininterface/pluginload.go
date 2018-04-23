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
	"plugin"

	"github.com/tinyzimmer/gomapper/logging"
	"github.com/tinyzimmer/gomapper/nmapresult"
)

type LoadedPlugins struct {
	Plugins []GomapperPlugin
}

type GomapperPlugin struct {
	Name            string
	SymbolInterface *plugin.Plugin
}

func LoadPlugins(plugins []string) (loadedPlugins LoadedPlugins) {
	for _, mod := range plugins {
		if mod != "" {
			pluginFile := fmt.Sprintf("plugins/%s/%s.so", mod, mod)
			p, err := plugin.Open(pluginFile)
			if err != nil {
				logging.LogError(err.Error())
			} else {
				loadedPlugin := GomapperPlugin{}
				loadedPlugin.Name = mod
				loadedPlugin.SymbolInterface = p
				loadedPlugins.Plugins = append(loadedPlugins.Plugins, loadedPlugin)
				logging.LogInfo(fmt.Sprintf("Loaded plugin: %s", mod))
			}
		}
	}
	return
}

func (g GomapperPlugin) OnScanComplete(nmapRun *nmapresult.NmapRun, db interface{}) {
	run, err := g.SymbolInterface.Lookup("OnScanComplete")
	if err != nil {
		logging.LogWarn(g.formatError(err))
		return
	}
	runFunc, ok := run.(func(*nmapresult.NmapRun, interface{}) error)
	if !ok {
		err := errors.New("Bad OnScanComplete implementation")
		logging.LogError(g.formatError(err))
		return
	}
	if err := runFunc(nmapRun, db); err != nil {
		logging.LogError(g.formatError(err))
	}
	return
}

func (g GomapperPlugin) formatError(err error) string {
	return fmt.Sprintf("%s: %s", g.Name, err.Error())
}
