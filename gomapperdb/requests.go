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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tinyzimmer/gomapper/config"
	"github.com/tinyzimmer/gomapper/formats"
	"github.com/tinyzimmer/gomapper/logging"
)

func logRequest(req *http.Request) {
	msg := fmt.Sprintf("%s : %s %s : %s", req.RemoteAddr, req.Method, req.Proto, req.RequestURI)
	logging.LogInfo(fmt.Sprintf("Received Request: %s", msg))
}

func formatResponse(data interface{}) (string, error) {
	dumped, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(dumped) + "\n", nil
}

func (db MemoryDatabase) RunPlugins(config config.Configuration) {
	for _, plugin := range db.Plugins.Plugins {
		for _, netw := range config.Discovery.Networks {
			dbNetwork, err := plugin.ScanNetwork(netw)
			if err != nil {
				logging.LogError(err.Error())
			} else {
				db.AddNetwork(dbNetwork)
			}
		}
		discoveredNetworks, err := plugin.DiscoverNetworks()
		if err != nil {
			logging.LogError(err.Error())
		} else {
			for _, netw := range discoveredNetworks {
				netString := netw.String()
				dbNetwork, err := plugin.ScanNetwork(netString)
				if err != nil {
					logging.LogError(err.Error())
				} else {
					db.AddNetwork(dbNetwork)
				}
			}
		}
	}
}

func (db MemoryDatabase) ReceivedScan(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	decoder := json.NewDecoder(req.Body)
	input := &formats.ReqInput{}
	err := decoder.Decode(&input)
	if err != nil {
		logging.LogError("Invalid Json in request payload")
		response, _ := formatResponse(err)
		io.WriteString(w, response)
		return
	} else {
		logging.LogInfo(fmt.Sprintf("Parsed Request: %s", input))
	}
	defer req.Body.Close()
	for _, plugin := range db.Plugins.Plugins {
		for _, method := range plugin.Methods {
			if method == input.Method {
				res, dbNetwork := plugin.HandleScanRequest(input)
				db.AddNetwork(dbNetwork)
				response, _ := formatResponse(res)
				io.WriteString(w, response)
			}
		}
	}
}

func (db MemoryDatabase) IterateNetworks(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	data := db.GetAllNetworks()
	response, _ := formatResponse(data)
	io.WriteString(w, response)
}
