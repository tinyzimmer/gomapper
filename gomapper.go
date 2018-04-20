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
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr, err := getAddr()
	if err != nil {
		log.Fatal(err)
	} else {
		if isPrivateAddr(addr) {
			log.Println(fmt.Sprintf("Listening on private address: %s:8080", addr))
		} else {
			log.Println(fmt.Sprintf("Listening on public address %s:8080", addr))
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/scan", receivedScan)
		http.ListenAndServe(fmt.Sprintf("%s:8080", addr), mux)
	}
}
