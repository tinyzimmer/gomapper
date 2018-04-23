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

	"github.com/tinyzimmer/gomapper/formats"
	"github.com/tinyzimmer/gomapper/logging"
	"github.com/tinyzimmer/gomapper/scanner"
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

func (db MemoryDatabase) ReceivedScan(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	decoder := json.NewDecoder(req.Body)
	input := &formats.ReqInput{}
	err := decoder.Decode(&input)
	if err != nil {
		logging.LogError("Invalid JSON in request payload")
		logging.LogError(fmt.Sprintf("\t:%s", input))
		io.WriteString(w, "{\"error\": \"invalid request json\"}\n")
		return
	} else {
		logging.LogInfo(fmt.Sprintf("Parsed Request: %s", input))
	}
	defer req.Body.Close()
	scanner, err := scanner.RequestScanner(input)
	if err != nil {
		errString := err.Error()
		logging.LogError(errString)
		io.WriteString(w, fmt.Sprintf("{\"error\": \"%s\"}\n", err))
		return
	} else {
		logging.LogInfo("Initiating nmap scan")
	}
	scanner.RunScan()
	if scanner.Failed {
		logging.LogWarn("User requested scan failed")
		response, _ := formatResponse(scanner.Error)
		io.WriteString(w, response)
	} else {
		logging.LogInfo("Adding scan results to graph")
		db.AddScanResultsByNetwork(scanner.Target, scanner.Results)
		logging.LogInfo("Returning scan results")
		response, _ := formatResponse(scanner.Results)
		io.WriteString(w, response)
		for _, plugin := range db.Plugins.Plugins {
			plugin.OnScanComplete(scanner.Results, db)
		}
	}
}

func (db MemoryDatabase) IterateNetworks(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	data := db.GetAllNetworks()
	response, _ := formatResponse(data)
	io.WriteString(w, response)
}
