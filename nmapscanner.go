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
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Scanner struct {
	Executable string
	Target     string
	Args       []string
	Xml        string
	Results    interface{}
}

type ErrorResponse struct {
	Error  string
	Stderr string
}

func InitScanner(target string, args []string) Scanner {
	scanner := Scanner{}
	scanner.Target = target
	scanner.Args = args
	scanner.Executable = "nmap"
	file, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	scanner.Xml = file.Name()
	defer os.Remove(file.Name())
	return scanner
}

func (s *Scanner) SetExec(executable string) {
	s.Executable = executable
}

func (s *Scanner) RunScan() {
	args := strings.Join(s.Args[:], " ")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(s.Executable, args, s.Target, "-oX", s.Xml)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + stderr.String())
		response := ErrorResponse{}
		response.Error = fmt.Sprint(err)
		response.Stderr = stderr.String()
		s.Results = response
	}
	s.Results = s.ParseRun()
}

func (s *Scanner) ParseRun() *NmapRun {
	file, _ := ioutil.ReadFile(s.Xml)
	res := &NmapRun{}
	xml.Unmarshal([]byte(string(file)), &res)
	replaceString := fmt.Sprintf("-oX %s ", s.Xml)
	res.Args = strings.Replace(res.Args, replaceString, "", int(1))
	return res
}
