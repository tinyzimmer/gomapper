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
	"io"
	"log"
	"net/http"
)

func receivedScan(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	input := &ReqInput{}
	err := decoder.Decode(&input)
	if err != nil {
		log.Println("Invalid JSON in request payload")
		log.Println(input)
		io.WriteString(w, "{\"error\": \"invalid request json\"}\n")
		return
	}
	defer req.Body.Close()
	scanner, err := RequestScanner(input)
	if err != nil {
		log.Println("Failed to initiate scanner")
		log.Println(err)
		io.WriteString(w, "{\"error\": \"failed to initiate a scanner\"}\n")
		return
	}
	scanner.RunScan()
	dumped, err := json.MarshalIndent(scanner.Results, "", "    ")
	if err != nil {
		errString := err.Error()
		io.WriteString(w, errString)
		return
	}
	io.WriteString(w, string(dumped)+"\n")
}
