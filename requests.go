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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
)

func logRequest(req *http.Request) {
	msg := fmt.Sprintf("%s : %s %s : %s", req.RemoteAddr, req.Method, req.Proto, req.RequestURI)
	logInfo(fmt.Sprintf("Received Request: %s", msg))
}

func formatResponse(data interface{}) (string, error) {
	dumped, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(dumped) + "\n", nil
}

func receivedScan(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	decoder := json.NewDecoder(req.Body)
	input := &ReqInput{}
	err := decoder.Decode(&input)
	if err != nil {
		logError("Invalid JSON in request payload")
		logError(fmt.Sprintf("\t:%s", input))
		io.WriteString(w, "{\"error\": \"invalid request json\"}\n")
		return
	} else {
		logInfo(fmt.Sprintf("Parsed Request: %s", input))
	}
	defer req.Body.Close()
	scanner, err := RequestScanner(input)
	if err != nil {
		errString := err.Error()
		logError(errString)
		io.WriteString(w, fmt.Sprintf("{\"error\": \"%s\"}\n", err))
		return
	} else {
		logInfo("Initiating nmap scan")
	}
	scanner.RunScan()
	if scanner.Failed {
		logWarn("User requested scan failed")
		response, _ := formatResponse(scanner.Error)
		io.WriteString(w, response)
	} else {
		logInfo("Returning scan results")
		response, _ := formatResponse(scanner.Results)
		io.WriteString(w, response)
	}
}

func (g Graph) IterateNetworks(w http.ResponseWriter, req *http.Request) {
	hosts := make(map[string]*GraphedNetwork)
	p := cayley.StartPath(g.Store, quad.String("Subnets")).Out(quad.String("Subnet"))
	p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		valueString := fmt.Sprint(nativeValue)
		_, ok := hosts[valueString]
		if !ok {
			network := &GraphedNetwork{}
			hosts[valueString] = network
		}
	})
	for subnet, graphedNetwork := range hosts {
		p := cayley.StartPath(g.Store, quad.String("Hosts")).Out(quad.String(subnet))
		p.Iterate(nil).EachValue(nil, func(value quad.Value) {
			nativeValue := quad.NativeOf(value)
			valueString := fmt.Sprint(nativeValue)
			graphedNetwork.Hosts = append(graphedNetwork.Hosts, valueString)
		})
	}
	response, _ := formatResponse(hosts)
	io.WriteString(w, response)
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

type GraphedNetwork struct {
	Hosts []string
}
