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
)

func receivedScan(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	input := &ReqInput{}
	err := decoder.Decode(&input)
	if err != nil {
		logError("Invalid JSON in request payload")
		logError(fmt.Sprintf("\t:s", input))
		io.WriteString(w, "{\"error\": \"invalid request json\"}\n")
		return
	} else {
		logInfo(fmt.Sprintf("Received Request: %s", input))
	}
	defer req.Body.Close()
	scanner, err := RequestScanner(input)
	if err != nil {
		logError("Failed to initiate scanner")
		io.WriteString(w, fmt.Sprintf("{\"error\": \"%s\"}\n", err))
		return
	} else {
		logInfo("Initiating nmap scan")
	}
	scanner.RunScan()
	logInfo("Returning scan results")
	dumped, err := json.MarshalIndent(scanner.Results, "", "    ")
	if err != nil {
		errString := err.Error()
		io.WriteString(w, errString)
		return
	}
	io.WriteString(w, string(dumped)+"\n")
}