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
	"os"
	"strings"

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
	Debug    bool
}

func getConfig(configFile *string) (config Configuration, err error) {
	configFileValue := *configFile
	if configFileValue != "" {
		config, err = decodeConfigurationFile(configFileValue)
		return
	}
	config = parseEnvironmentConfiguration()
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
	config.Server.ListenAddress = os.Getenv("GOMAPPER_LISTEN_ADDRESS")
	config.Server.ListenPort = os.Getenv("GOMAPPER_LISTEN_PORT")
	config.Discovery.Enabled = checkEnvBool(os.Getenv("GOMAPPER_DISCOVERY_ENABLED"), true)
	config.Discovery.Mode = os.Getenv("GOMAPPER_DISCOVERY_MODE")
	config.Discovery.Debug = checkEnvBool(os.Getenv("GOMAPPER_DISCOVERY_DEBUG"), false)
	config.Discovery.Networks = checkEnvNetworks(os.Getenv("GOMAPPER_DISCOVERY_NETWORKS"))
	if undefined(config.Server.ListenAddress) {
		config.Server.ListenAddress = getDefault("listenAddress")
	}
	if undefined(config.Server.ListenPort) {
		config.Server.ListenPort = getDefault("listenPort")
	}
	if undefined(config.Discovery.Mode) {
		config.Discovery.Mode = getDefault("discoveryMode")
	}
	return
}

func checkEnvNetworks(value string) (networks []string) {
	if value == "" {
		return
	}
	networks = strings.Split(value, ",")
	return
}

func checkEnvBool(value string, def bool) bool {
	if value == "0" {
		return false
	}
	if value == "1" {
		return true
	}
	if value == "false" {
		return false
	}
	if value == "true" {
		return true
	}
	return def
}

func undefined(value string) bool {
	if value == "" {
		return true
	}
	return false
}

func getDefault(config string) (value string) {
	if config == "listenAddress" {
		value = "127.0.0.1"
	}
	if config == "listenPort" {
		value = "8080"
	}
	if config == "discoveryMode" {
		value = "ping"
	}
	return
}
