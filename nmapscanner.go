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
	"log"
	"os/exec"
	"strings"
)

func RunScan(input *InputPayload) *NmapRun {
	log.Println(input)
	xmlFile := "out.xml"
	args := strings.Join(input.Args[:], " ")
	cmd := exec.Command("nmap", args, input.Target, "-oX", xmlFile)
	err := cmd.Run()
	results := ParseRun(xmlFile)
	if err != nil {
		log.Fatalln(err)
	}
	return results
}
