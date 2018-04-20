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
	"errors"
	"fmt"
)

type ReqInput struct {
	CustomExec string   `json:"customExec"`
	Target     string   `json:"target"`
	Method     string   `json:"method"`
	Ports      string   `json:"ports"`
	Detection  string   `json:"detection"`
	Script     string   `json:"script"`
	ScriptArgs string   `json:"scriptArgs"`
	RawArgs    []string `json:"rawArgs"`
}

func RequestScanner(input *ReqInput) (Scanner, error) {
	scanner := Scanner{}
	scanner.Failed = false
	scanner.ReqInput = input
	xml, err := getOutXml()
	if err != nil {
		logError(err.Error())
		return scanner, err
	}
	scanner.Xml = xml
	target, err := checkTarget(input.Target)
	if err != nil {
		logError(err.Error())
		return scanner, err
	} else {
		scanner.SetTarget(target)
	}
	if len(input.RawArgs) > 0 {
		scanner.SetRawInput(input.RawArgs)
	} else {
		scanner.SetHelperInput()
	}
	if input.CustomExec != "" {
		scanner.SetExec(input.CustomExec)
	} else {
		scanner.SetExec("nmap")
	}
	return scanner, nil
}

func checkTarget(target string) (string, error) {
	if target == "" {
		err := errors.New("No target provided")
		return "", err
	}
	return target, nil
}

func GetHelperArgs(input *ReqInput, xml string) ([]string, error) {
	var computedArgs []string
	if input.Method == "tcp-connect" {
		computedArgs = append(computedArgs, "-sT")
	} else if input.Method == "tcp-syn" {
		computedArgs = append(computedArgs, "-sS")
	} else if input.Method == "tcp-ack" {
		computedArgs = append(computedArgs, "-sA")
	} else if input.Method == "udp" {
		computedArgs = append(computedArgs, "-sU")
	} else if input.Method != "" {
		err := errors.New("Invalid scan method")
		return nil, err
	}
	if input.Detection == "full" {
		computedArgs = append(computedArgs, "-A")
	} else if input.Detection == "os" {
		computedArgs = append(computedArgs, "-O")
	} else if input.Detection != "" {
		err := errors.New("Invalid Detection Method")
		return nil, err
	}
	if input.Script != "" {
		computedArgs = append(computedArgs, fmt.Sprintf("--script=%s", input.Script))
	}
	if input.ScriptArgs != "" {
		computedArgs = append(computedArgs, fmt.Sprintf("--script-args=%s", input.ScriptArgs))
	}
	if input.Ports != "" {
		computedArgs = append(computedArgs, "-p")
		computedArgs = append(computedArgs, input.Ports)
	}
	computedArgs = append(computedArgs, "-oX")
	computedArgs = append(computedArgs, xml)
	computedArgs = append(computedArgs, input.Target)
	return computedArgs, nil
}
