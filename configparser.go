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
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Server    ServerConfig
	Discovery DiscoveryConfig
}

type ServerConfig struct {
	ListenAddress string `toml:"listen"`
	ListenPort    string `toml:"port"`
}

type DiscoveryConfig struct {
	Enabled  bool
	Mode     string
	Networks []string
}

func getConfig(configFile *string) (config Configuration, err error) {
	configFileValue := *configFile
	if configFileValue != "" {
		config, err = decodeConfigurationFile(configFileValue)
		return
	}
	config = getDefaultConfig()
	return
}

func decodeConfigurationFile(configFile string) (config Configuration, err error) {
	configText, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}
	_, err = toml.Decode(string(configText), &config)
	return
}

func parseEnvironmentConfiguration() (config Configuration) {
	return
}

func getDefaultConfig() (config Configuration) {
	config.Server.ListenAddress = "127.0.0.1"
	config.Server.ListenPort = "8080"
	config.Discovery.Enabled = true
	config.Discovery.Mode = "ping"
	return
}
