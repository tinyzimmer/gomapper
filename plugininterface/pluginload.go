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
	"fmt"
	"plugin"

	"github.com/tinyzimmer/gomapper/logging"
	"github.com/tinyzimmer/gomapper/nmapresult"
)

type LoadedPlugins struct {
	Plugins []*GomapperPlugin
}

type GomapperPlugin struct {
	Name         string
	Interface    *plugin.Plugin
	RunInterface interface{}
}

func LoadPlugins(plugins []string) (loadedPlugins LoadedPlugins) {
	for _, mod := range plugins {
		if mod != "" {
			pluginFile := fmt.Sprintf("plugins/%s/%s.so", mod, mod)
			p, err := plugin.Open(pluginFile)
			if err != nil {
				logging.LogError(err.Error())
			} else {
				loadedPlugin := &GomapperPlugin{}
				loadedPlugin.Name = mod
				loadedPlugin.Interface = p
				loadedPlugin.checkInterfaces()
				loadedPlugins.Plugins = append(loadedPlugins.Plugins, loadedPlugin)
				logging.LogInfo(fmt.Sprintf("Loaded plugin: %s", mod))
			}
		}
	}
	return
}

func (g GomapperPlugin) checkInterfaces() {
	r, err := g.Interface.Lookup("NmapRunInterface")
	if err != nil {
		logging.LogWarn(g.formatError(err))
	} else {
		g.RunInterface = *r.(*nmapresult.NmapRun)
	}
	return
}

func (g GomapperPlugin) OnScanComplete(nmapRun *nmapresult.NmapRun) {
	f, err := g.Interface.Lookup("OnScanComplete")
	if err != nil {
		logging.LogError(g.formatError(err))
		return
	}
	if g.RunInterface != nil {
		g.RunInterface = *nmapRun
	}
	f.(func())()
}

func (g GomapperPlugin) formatError(err error) string {
	return fmt.Sprintf("%s: %s", g.Name, err.Error())
}
