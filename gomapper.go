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
	"net/http"
)

type InputPayload struct {
	Target string   `json:"target"`
	Args   []string `json:"args"`
}

func receivedScan(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	input := &InputPayload{}
	err := decoder.Decode(&input)
	if err != nil {
		errString := err.Error()
		io.WriteString(w, errString)
		return
	}
	defer req.Body.Close()
	results := RunScan(input)
	dumped, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		errString := err.Error()
		io.WriteString(w, errString)
		return
	}
	io.WriteString(w, string(dumped))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/scan", receivedScan)
	http.ListenAndServe(":8080", mux)
}
