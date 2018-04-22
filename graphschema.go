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

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
)

type Graph struct {
	Store *graph.Handle
}

func getMemoryGraph() (graph Graph, err error) {
	// Create a brand new graph
	store, err := cayley.NewMemoryGraph()
	if err != nil {
		return
	}
	graph.Store = store
	return
}

func (g Graph) GetSubnets() (subnets []string) {
	p := cayley.StartPath(g.Store, quad.String("Subnets")).Out(quad.String("Subnet"))
	p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		valueString := fmt.Sprint(nativeValue)
		if !contains(subnets, valueString) {
			subnets = append(subnets, valueString)
		}
	})
	return
}

func (g Graph) GetHostsBySubnet(subnet string) (hosts []string) {
	p := cayley.StartPath(g.Store, quad.String("Hosts")).Out(quad.String(subnet))
	p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		valueString := fmt.Sprint(nativeValue)
		hosts = append(hosts, valueString)
	})
	return
}

func (g Graph) AddScanResultsByNetwork(network string, results *NmapRun) {
	subnets := g.GetSubnets()
	for _, host := range results.Hosts {
		for _, address := range host.Addresses {
			if address.AddrType == "ipv4" {
				if network != "" {
					g.AddHost(network, address.Addr)
					logInfo(fmt.Sprintf("Added %s to memory graph under %s", address.Addr, network))
				} else {
					formatDefault := formatDefaultNetmask(address.Addr)
					if !contains(subnets, formatDefault) {
						logInfo(fmt.Sprintf("Adding new subnet to memory graph: %s", formatDefault))
						g.AddNetwork(formatDefault)
						subnets = g.GetSubnets()
					}
					g.AddHost(formatDefault, address.Addr)
					logInfo(fmt.Sprintf("Added %s to memory graph under %s", address.Addr, formatDefault))
				}
			}
		}
	}
}

func (g Graph) AddNetwork(network string) {
	g.Store.AddQuad(quad.Make("Subnets", "Subnet", network, nil))
}

func (g Graph) AddHost(network string, host string) {
	g.Store.AddQuad(quad.Make("Hosts", network, host, nil))
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}
