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
	"os"
	"os/exec"
)

type Scanner struct {
	Executable   string
	ReqInput     *ReqInput
	Target       string
	RawArgs      []string
	ComputedArgs []string
	RawEnforce   bool
	Xml          string
	Results      *NmapRun
	Error        ErrorResponse
	Failed       bool
	Debug        bool
}

type ErrorResponse struct {
	Error  string
	Stderr string
}

func InitScanner(target string, debug bool) (Scanner, error) {
	scanner := Scanner{}
	scanner.SetTarget(target)
	scanner.Failed = false
	scanner.SetExec("nmap")
	scanner.SetDebug(debug)
	xml, err := getOutXml()
	if err != nil {
		err := errors.New("Could not initiate scanner")
		return scanner, err
	}
	scanner.Xml = xml
	return scanner, nil
}

func (s *Scanner) SetDebug(debug bool) {
	if debug {
		logDebug("Initializing scanner with debug")
	}
	s.Debug = debug
}

func (s *Scanner) SetHelperInput() {
	s.RawEnforce = false
}

func (s *Scanner) SetRawInput(args []string) {
	s.RawEnforce = true
	s.RawArgs = args
}

func (s *Scanner) SetStealthDiscovery() {
	s.RawEnforce = true
	s.RawArgs = append(s.RawArgs, "-P0")
	s.RawArgs = append(s.RawArgs, "-sS")
}

func (s *Scanner) SetPingDiscovery() {
	s.RawEnforce = true
	s.RawArgs = append(s.RawArgs, "-sn")
}

func (s *Scanner) SetConnectDiscovery() {
	s.RawEnforce = true
	s.RawArgs = append(s.RawArgs, "-sT")
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
	args, err := GetHelperArgs(s.ReqInput, s.Xml)
	if err != nil {
		s.Failed = true
		s.HandleReturn(err, "", "")
	} else {
		cmd := exec.Command(s.Executable, args...)
		if s.Debug {
			logDebug(fmt.Sprint(cmd))
		}
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err := cmd.Run()
		s.HandleReturn(err, stdout.String(), stderr.String())
	}
}

func (s *Scanner) RunRawArgScan() {
	stdout, stderr := createPipes()
	rawArgs := append(s.RawArgs, "-oX")
	rawArgs = append(rawArgs, s.Xml)
	rawArgs = append(rawArgs, s.Target)
	cmd := exec.Command(s.Executable, rawArgs...)
	if s.Debug {
		logDebug(fmt.Sprint(cmd))
	}
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
		results := ParseRun(s.Xml)
		s.Results = results
		logInfo(fmt.Sprintf("Scan completed in %s seconds", floatToString(results.RunStats.Finished.Elapsed)))
	}
}

func (s *Scanner) ReturnFail(err error, msg string) {
	if msg != "" {
		logError(fmt.Sprint(err) + ": " + msg)
	} else {
		logError(fmt.Sprint(err))
	}
	response := ErrorResponse{}
	response.Error = fmt.Sprint(err)
	response.Stderr = msg
	s.Failed = true
	s.Error = response
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

func getOutXml() (string, error) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		logError(err.Error())
		return "", err
	}
	return file.Name(), nil
}
