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
	"encoding/xml"
	"io/ioutil"
)

type NmapRun struct {
	Scanner   string    `xml:"scanner,attr"`
	Args      string    `xml:"args,attr"`
	Start     int       `xml:"start,attr"`
	StartStr  string    `xml:"startstr,attr"`
	Version   string    `xml:"version,attr"`
	ScanInfo  ScanInfo  `xml:"scaninfo"`
	Verbose   Verbose   `xml:"verbose"`
	Debugging Debugging `xml:"debugging"`
	Hosts     []Host    `xml:"host"`
	RunStats  RunStats  `xml:"runstats"`
}

type ScanInfo struct {
	Type        string `xml:"type,attr"`
	Protocol    string `xml:"protocol,attr"`
	NumServices int    `xml:"numservices,attr"`
	Services    string `xml:"services,attr"`
}

type Host struct {
	StartTime int       `xml:"starttime,attr"`
	EndTime   int       `xml:"endtime,attr"`
	Status    Status    `xml:"status"`
	Address   Address   `xml:"address"`
	Hostnames Hostnames `xml:"hostnames"`
	Ports     Ports     `xml:"ports"`
	Times     Times     `xml:"times"`
}

type Status struct {
	State     string `xml:"state,attr"`
	Reason    string `xml:"reason,attr"`
	ReasonTTL int    `xml:"reason_ttl,attr"`
}

type Address struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype"`
}

type Hostnames struct {
	Hostnames []Hostname `xml:"hostname"`
}

type Hostname struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type Ports struct {
	ExtraPorts ExtraPorts `xml:"extraports"`
	Ports      []Port     `xml:"port"`
}

type Port struct {
	Protocol string  `xml:"protocol,attr"`
	PortId   int     `xml:"portid,attr"`
	State    Status  `xml:"state"`
	Service  Service `xml:"service"`
}

type Service struct {
	Name   string `xml:"name,attr"`
	Method string `xml:"method,attr"`
	Conf   string `xml:"conf,attr"`
}

type ExtraPorts struct {
	State        string        `xml:"state,attr"`
	Count        int           `xml:"count,attr"`
	ExtraReasons []ExtraReason `xml:"extrareasons"`
}

type ExtraReason struct {
	Reason string `xml:"reason,attr"`
	Count  int    `xml:"count,attr"`
}

type Times struct {
	Srtt   int `xml:"srtt,attr"`
	Rttvar int `xml:"rttvar,attr"`
	To     int `xml:"to,attr"`
}

type RunStats struct {
	Finished      Finished      `xml:"finished"`
	FinishedHosts FinishedHosts `xml:"hosts"`
}

type Finished struct {
	Time    int     `xml:"time,attr"`
	TimeStr string  `xml:"timestr,attr"`
	Elapsed float32 `xml:"elapsed,attr"`
	Summary string  `xml:"summary,attr"`
	Exit    string  `xml:"exit,attr"`
}

type FinishedHosts struct {
	Up    int `xml:"up,attr"`
	Down  int `xml:"down,attr"`
	Total int `xml:"total,attr"`
}

type Verbose struct {
	Level int `xml:"level,attr"`
}

type Debugging struct {
	Level int `xml:level,attr`
}

func ParseRun(result string) *NmapRun {
	file, _ := ioutil.ReadFile(result)
	res := &NmapRun{}
	xml.Unmarshal([]byte(string(file)), &res)
	return res
}
