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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Scanner struct {
	Executable   string
	Input        *ScannerInput
	Target       string
	RawArgs      []string
	ComputedArgs []string
	RawEnforce   bool
	Xml          string
	Results      interface{}
	Failed       bool
}

type ErrorResponse struct {
	Error  string
	Stderr string
}

func InitScanner(input *ScannerInput) Scanner {
	scanner := Scanner{}
	scanner.Failed = false
	scanner.Input = input
	file, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	scanner.Xml = file.Name()
	scanner.ParseInput()
	return scanner
}

func (s *Scanner) ParseInput() {
	if s.Input.Target == "" {
		err := errors.New("No target provided")
		log.Println(err)
		s.Failed = true
		s.ReturnFail(err, "")
	} else {
		s.SetTarget(s.Input.Target)
	}
	if s.Input.CustomExec != "" {
		s.SetExec(s.Input.CustomExec)
	} else {
		s.SetExec("nmap")
	}
	if len(s.Input.RawArgs) > 0 {
		s.SetRawInput(s.Input.RawArgs)
	} else {
		s.SetHelperInput()
	}
}

func (s *Scanner) GetHelperArgs() []string {
	var computedArgs []string
	if s.Input.Method == "tcp-connect" {
		computedArgs = append(computedArgs, "-sT")
	} else if s.Input.Method == "tcp-syn" {
		computedArgs = append(computedArgs, "-sS")
	} else if s.Input.Method == "tcp-ack" {
		computedArgs = append(computedArgs, "-sA")
	} else if s.Input.Method == "udp" {
		computedArgs = append(computedArgs, "-sU")
	}
	if s.Input.Ports != "" {
		computedArgs = append(computedArgs, "-p")
		computedArgs = append(computedArgs, s.Input.Ports)
	}
	computedArgs = append(computedArgs, "-oX")
	computedArgs = append(computedArgs, s.Xml)
	computedArgs = append(computedArgs, s.Target)
	return computedArgs
}

func (s *Scanner) SetHelperInput() {
	s.RawEnforce = false
}

func (s *Scanner) SetRawInput(args []string) {
	s.RawEnforce = true
	s.RawArgs = args
}

func (s *Scanner) SetExec(executable string) {
	s.Executable = executable
}

func (s *Scanner) SetTarget(target string) {
	s.Target = target
}

func (s *Scanner) RunScan() {
	if s.RawEnforce == true {
		s.RunRawArgScan()
	} else {
		s.RunHelperScan()
	}
}

func (s *Scanner) RunHelperScan() {
	stdout, stderr := createPipes()
	args := s.GetHelperArgs()
	cmd := exec.Command(s.Executable, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	s.HandleReturn(err, stdout.String(), stderr.String())
}

func (s *Scanner) RunRawArgScan() {
	stdout, stderr := createPipes()
	rawArgs := append(s.RawArgs, "-oX")
	rawArgs = append(rawArgs, s.Xml)
	rawArgs = append(rawArgs, s.Target)
	cmd := exec.Command(s.Executable, rawArgs...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	s.HandleReturn(err, stdout.String(), stderr.String())
}

func (s *Scanner) HandleReturn(err error, stdout string, stderr string) {
	if err != nil {
		msg := stderr
		s.ReturnFail(err, msg)
	} else {
		s.Results = ParseRun(s.Xml)
	}
}

func (s *Scanner) ReturnFail(err error, msg string) {
	log.Println(fmt.Sprint(err) + ": " + msg)
	response := ErrorResponse{}
	response.Error = fmt.Sprint(err)
	response.Stderr = msg
	s.Results = response
}

func ParseRun(filePath string) *NmapRun {
	file, _ := ioutil.ReadFile(filePath)
	res := &NmapRun{}
	xml.Unmarshal([]byte(string(file)), &res)
	os.Remove(filePath)
	return res
}

func createPipes() (*bytes.Buffer, *bytes.Buffer) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	return &stdout, &stderr
}
