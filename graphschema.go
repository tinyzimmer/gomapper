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

func (g Graph) AddNetwork(network string) {
	g.Store.AddQuad(quad.Make("Subnets", "Subnet", network, nil))
}

func (g Graph) AddHost(network string, host string) {
	g.Store.AddQuad(quad.Make("Hosts", network, host, nil))
}
