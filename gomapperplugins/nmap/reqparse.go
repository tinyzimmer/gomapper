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

package nmap

import (
	"errors"
	"fmt"

	"github.com/tinyzimmer/gomapper/formats"
	"github.com/tinyzimmer/gomapper/logging"
)

func RequestScanner(input *formats.ReqInput, conf map[string]interface{}) (Scanner, error) {
	scanner := Scanner{}
	scanner.Failed = false
	scanner.ReqInput = input
	xml, err := getOutXml()
	if err != nil {
		logging.LogError(err.Error())
		return scanner, err
	}
	scanner.Xml = xml
	target, err := checkTarget(input.Target)
	if err != nil {
		logging.LogError(err.Error())
		return scanner, err
	} else {
		scanner.SetTarget(target)
	}
	scanner.SetHelperInput(input.Method, input.Options)
	scanner.SetExec("nmap")
	return scanner, nil
}

func checkTarget(target string) (string, error) {
	if target == "" {
		err := errors.New("No target provided")
		return "", err
	}
	return target, nil
}

func checkScanMethod(method string) (string, error) {
	if method == "nmap_connect" {
		return "-sT", nil
	} else if method == "nmap_syn" {
		return "-sS", nil
	} else if method == "nmap_ack" {
		return "-sA", nil
	} else if method == "nmap_udp" {
		return "-sU", nil
	} else if method == "nmap_ping" {
		return "-sn", nil
	} else if method != "" {
		err := errors.New("Invalid scan method")
		return "", err
	}
	return "", nil

}

func checkDetectionMethod(method string) (string, error) {
	if method == "full" {
		return "-A", nil
	} else if method == "os" {
		return "-O", nil
	} else if method != "" {
		err := errors.New("Invalid Detection Method")
		return "", err
	}
	return "", nil
}

func GetHelperArgs(method string, opts map[string]string) ([]string, error) {
	logging.LogInfo("Determining scan arguments")
	var computedArgs []string
	method, err := checkScanMethod(method)
	if err == nil {
		computedArgs = append(computedArgs, method)
	} else {
		return nil, err
	}
	detect, detectOk := opts["detection"]
	ports, portsOk := opts["ports"]
	script, scriptOk := opts["script"]
	scriptArgs, scriptArgsOk := opts["scriptArgs"]
	if portsOk {
		computedArgs = append(computedArgs, "-p")
		computedArgs = append(computedArgs, ports)
	}
	if detectOk {
		detection, err := checkDetectionMethod(detect)
		if err != nil {
			return nil, err
		}
		computedArgs = append(computedArgs, detection)
	}
	if scriptOk {
		computedArgs = append(computedArgs, fmt.Sprintf("--script=%s", script))
		if scriptArgsOk {
			computedArgs = append(computedArgs, fmt.Sprintf("--script-args=%s", scriptArgs))
		}
	}
	return computedArgs, nil
}
